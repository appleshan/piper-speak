#!/bin/bash
# piper-speak installer
set -e

echo "Installing piper-speak..."

# Check dependencies
check_dep() {
    if ! command -v "$1" &>/dev/null; then
        echo "Error: $1 is required but not installed" >&2
        echo "Install with: $2" >&2
        exit 1
    fi
}

check_dep piper-tts "yay -S piper-tts-bin"
check_dep pw-play "pacman -S pipewire-pulse"
check_dep wl-paste "pacman -S wl-clipboard"
check_dep curl "pacman -S curl"

# Install scripts
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$INSTALL_DIR"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

cp "$SCRIPT_DIR/bin/piper-speak" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/bin/speak-selection" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/bin/piper-speak-install" "$INSTALL_DIR/"

chmod +x "$INSTALL_DIR/piper-speak"
chmod +x "$INSTALL_DIR/speak-selection"
chmod +x "$INSTALL_DIR/piper-speak-install"

echo "Installed scripts to $INSTALL_DIR"

# Download default voice
"$INSTALL_DIR/piper-speak-install"

echo ""
echo "Installation complete!"
echo ""
echo "Usage:"
echo "  echo 'Hello world' | piper-speak"
echo "  piper-speak 'Hello world'"
echo ""
echo "For speak-selection, add a keybinding in your window manager:"
echo "  Hyprland: bind = SUPER SHIFT, period, exec, speak-selection"
