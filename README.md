# piper-speak

Simple text-to-speech wrapper for [Piper TTS](https://github.com/rhasspy/piper) on Linux.

- **piper-speak** - Pipe text to speech with speed control
- **speak-selection** - Read highlighted text aloud (Wayland)
- **piper-speak-install** - Download voice models

## Installation

### [Arch Linux (AUR)](https://aur.archlinux.org/packages/piper-speak)

```bash
yay -S piper-speak
```

Includes the default voice model (`en_US-lessac-medium`).

### Manual

```bash
git clone https://github.com/kgn/piper-speak.git
cd piper-speak
./install.sh
```

### Dependencies

- [piper-tts](https://github.com/rhasspy/piper) - `yay -S piper-tts-bin`
- pipewire - `pacman -S pipewire pipewire-pulse`
- wl-clipboard - `pacman -S wl-clipboard`

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

### piper-speak-install

```bash
# Install default voice (en_US-lessac-medium)
piper-speak-install

# List available voices
piper-speak-install --list

# Install specific voice
piper-speak-install en_US-ryan-high
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PIPER_VOICE_DIR` | `~/.local/share/piper/voices` | Voice model directory |
| `PIPER_VOICE` | `en_US-lessac-medium` | Default voice model |

## Voices

Popular English voices:

| Voice | Quality | Description |
|-------|---------|-------------|
| `en_US-lessac-medium` | Medium | Clear American English (default) |
| `en_US-lessac-high` | High | Higher quality, slower |
| `en_US-amy-medium` | Medium | Female American English |
| `en_US-ryan-medium` | Medium | Male American English |
| `en_GB-alan-medium` | Medium | Male British English |

See all voices: [rhasspy/piper-voices](https://huggingface.co/rhasspy/piper-voices)

## License

MIT
