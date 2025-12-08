package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	systemVoiceDir   = "/usr/share/piper-speak/voices"
	defaultVoice     = "en_US-lessac-medium"
	defaultSpeed     = 0.7
	defaultChunkSize = 500
)

func main() {
	speed := flag.Float64("speed", defaultSpeed, "Length scale (lower = faster)")
	background := flag.Bool("bg", false, "Run in background")
	voice := flag.String("voice", "", "Voice model name")
	chunkSize := flag.Int("chunk-size", defaultChunkSize, "Max characters per chunk")
	help := flag.Bool("help", false, "Show help")
	flag.BoolVar(help, "h", false, "Show help")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	// Get voice from flag, env, or default
	voiceName := *voice
	if voiceName == "" {
		voiceName = os.Getenv("PIPER_VOICE")
	}
	if voiceName == "" {
		voiceName = defaultVoice
	}

	// Find voice model
	voicePath, err := findVoice(voiceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'piper-speak-install %s' to download it\n", voiceName)
		os.Exit(1)
	}

	// Get text from args or stdin
	text := strings.Join(flag.Args(), " ")
	if text == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		text = string(data)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		fmt.Fprintln(os.Stderr, "Error: No text provided")
		os.Exit(1)
	}

	if *background {
		// Fork to background
		cmd := exec.Command(os.Args[0], append([]string{
			"--speed", fmt.Sprintf("%f", *speed),
			"--voice", voiceName,
			"--chunk-size", fmt.Sprintf("%d", *chunkSize),
		}, text)...)
		cmd.Start()
		os.Exit(0)
	}

	// Split into chunks and speak with double buffering
	chunks := splitIntoChunks(text, *chunkSize)
	speakChunks(chunks, voicePath, *speed)
}

func printHelp() {
	fmt.Println(`Usage: echo "text" | piper-speak [options]
       piper-speak [options] "text to speak"

Options:
  --speed NUM       Length scale (default: 0.7, lower = faster)
  --bg              Run in background
  --voice NAME      Voice model name (default: en_US-lessac-medium)
  --chunk-size NUM  Max characters per chunk (default: 500)

Environment:
  PIPER_VOICE_DIR  User voice directory (default: ~/.local/share/piper/voices)
  PIPER_VOICE      Default voice model name

Voice models are loaded from user dir first, then /usr/share/piper-speak/voices`)
}

func findVoice(name string) (string, error) {
	// Check user directory first
	userDir := os.Getenv("PIPER_VOICE_DIR")
	if userDir == "" {
		home, _ := os.UserHomeDir()
		userDir = filepath.Join(home, ".local", "share", "piper", "voices")
	}

	userPath := filepath.Join(userDir, name+".onnx")
	if _, err := os.Stat(userPath); err == nil {
		return userPath, nil
	}

	// Check system directory
	systemPath := filepath.Join(systemVoiceDir, name+".onnx")
	if _, err := os.Stat(systemPath); err == nil {
		return systemPath, nil
	}

	return "", fmt.Errorf("voice model not found: %s", name)
}

func splitIntoChunks(text string, maxChunkSize int) []string {
	var chunks []string

	// Split by double newlines (paragraphs) first
	paragraphs := strings.Split(text, "\n\n")

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		// If paragraph is too long, split by sentences
		if len(para) > maxChunkSize {
			sentences := splitSentences(para)
			chunks = append(chunks, sentences...)
		} else {
			chunks = append(chunks, para)
		}
	}

	return chunks
}

func splitSentences(text string) []string {
	var sentences []string
	var current strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		r := scanner.Text()
		current.WriteString(r)

		// End of sentence markers
		if r == "." || r == "!" || r == "?" {
			s := strings.TrimSpace(current.String())
			if s != "" {
				sentences = append(sentences, s)
			}
			current.Reset()
		}
	}

	// Don't forget remaining text
	if s := strings.TrimSpace(current.String()); s != "" {
		sentences = append(sentences, s)
	}

	return sentences
}

func speakChunks(chunks []string, voicePath string, speed float64) {
	if len(chunks) == 0 {
		return
	}

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel for generated wav files ready to play
	type wavFile struct {
		path string
		err  error
	}
	readyChan := make(chan wavFile, 2) // Buffer of 2 for double buffering

	// Track temp files for cleanup
	var tempFiles []string
	cleanup := func() {
		for _, f := range tempFiles {
			os.Remove(f)
		}
	}
	defer cleanup()

	// Start generating first chunk
	go func() {
		wav, err := generateWav(chunks[0], voicePath, speed, 0)
		readyChan <- wavFile{wav, err}
	}()

	for i := 0; i < len(chunks); i++ {
		// Start generating next chunk while we wait/play current
		if i+1 < len(chunks) {
			nextIdx := i + 1
			go func(idx int) {
				wav, err := generateWav(chunks[idx], voicePath, speed, idx)
				readyChan <- wavFile{wav, err}
			}(nextIdx)
		}

		// Wait for current chunk to be ready
		select {
		case <-sigChan:
			return
		case wf := <-readyChan:
			if wf.err != nil {
				fmt.Fprintf(os.Stderr, "Error generating audio: %v\n", wf.err)
				continue
			}
			tempFiles = append(tempFiles, wf.path)

			// Play current chunk (blocking)
			playDone := make(chan error, 1)
			go func() {
				playDone <- playWav(wf.path)
			}()

			select {
			case <-sigChan:
				// Kill playback on interrupt
				exec.Command("pkill", "-f", "pw-play.*piper-speak").Run()
				return
			case err := <-playDone:
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error playing audio: %v\n", err)
				}
			}

			// Clean up played file immediately
			os.Remove(wf.path)
		}
	}
}

func generateWav(text, voicePath string, speed float64, index int) (string, error) {
	tmpFile := fmt.Sprintf("/tmp/piper-speak-%d-%d.wav", os.Getpid(), index)

	cmd := exec.Command("piper-tts",
		"--model", voicePath,
		"--length_scale", fmt.Sprintf("%.2f", speed),
		"--output_file", tmpFile,
	)
	cmd.Stdin = strings.NewReader(text)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("piper-tts failed: %w", err)
	}

	return tmpFile, nil
}

func playWav(path string) error {
	cmd := exec.Command("pw-play", path)
	return cmd.Run()
}
