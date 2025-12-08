# piper-speak

Simple text-to-speech wrapper for [Piper TTS](https://github.com/rhasspy/piper) on Linux.

- **piper-speak** - Pipe text to speech with double-buffered playback
- **speak-selection** - Read highlighted text aloud (Wayland)

## Features

- Double-buffered audio generation for gapless playback
- Automatic text chunking (by paragraph, then sentence)
- Handles long documents without memory issues
- Graceful Ctrl+C handling

## Installation

### [Arch Linux (AUR)](https://aur.archlinux.org/packages/piper-speak)

```bash
yay -S piper-speak
```

Includes the default voice model (`en_US-lessac-medium`).

### Build from source

Requires Go 1.21+

```bash
git clone https://github.com/kgn/piper-speak.git
cd piper-speak
go build -o piper-speak ./cmd/piper-speak/
sudo install -Dm755 piper-speak /usr/bin/piper-speak
sudo install -Dm755 scripts/speak-selection /usr/bin/speak-selection

# Download a voice model
mkdir -p ~/.local/share/piper/voices
curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx" \
    -o ~/.local/share/piper/voices/en_US-lessac-medium.onnx
curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx.json" \
    -o ~/.local/share/piper/voices/en_US-lessac-medium.onnx.json
```

### Dependencies

- [piper-tts](https://github.com/rhasspy/piper) - `yay -S piper-tts-bin`
- pipewire - `pacman -S pipewire pipewire-pulse`
- wl-clipboard - `pacman -S wl-clipboard` (for speak-selection)

## Usage

### piper-speak

```bash
# Pipe text
echo "Hello world" | piper-speak

# Direct argument
piper-speak "Hello world"

# Adjust speed (lower = faster, default: 0.7)
echo "Slower speech" | piper-speak --speed 1.0

# Run in background
echo "Background speech" | piper-speak --bg

# Read a long document (automatically chunked)
cat large-document.txt | piper-speak

# Adjust chunk size (default: 500 characters)
cat book.txt | piper-speak --chunk-size 800
```

### speak-selection

Reads highlighted text aloud. Press the hotkey again to stop.

Add a keybinding in your window manager:

**Hyprland:**
```conf
bind = SUPER SHIFT, period, exec, speak-selection
```

**Sway:**
```conf
bindsym $mod+Shift+period exec speak-selection
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PIPER_VOICE_DIR` | `~/.local/share/piper/voices` | User voice model directory |
| `PIPER_VOICE` | `en_US-lessac-medium` | Default voice model |

Voice models are searched in user directory first, then `/usr/share/piper-speak/voices`.

## Voices

Download additional voices from [rhasspy/piper-voices](https://huggingface.co/rhasspy/piper-voices):

```bash
# Example: download a different voice
curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ryan/high/en_US-ryan-high.onnx" \
    -o ~/.local/share/piper/voices/en_US-ryan-high.onnx
curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ryan/high/en_US-ryan-high.onnx.json" \
    -o ~/.local/share/piper/voices/en_US-ryan-high.onnx.json

# Use it
echo "Hello" | piper-speak --voice en_US-ryan-high
```

Popular English voices:

| Voice | Quality | Description |
|-------|---------|-------------|
| `en_US-lessac-medium` | Medium | Clear American English (default) |
| `en_US-lessac-high` | High | Higher quality, slower |
| `en_US-amy-medium` | Medium | Female American English |
| `en_US-ryan-medium` | Medium | Male American English |
| `en_GB-alan-medium` | Medium | Male British English |

## License

MIT
