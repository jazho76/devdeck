#!/usr/bin/env bash
set -euo pipefail

NERD_FONTS_VERSION="3.4.0"
FONT_NAME="DroidSansMono"
FONT_DISPLAY_NAME="DroidSansM Nerd Font"
FONT_ARCHIVE="${FONT_NAME}.zip"
FONT_URL="https://github.com/ryanoasis/nerd-fonts/releases/download/v${NERD_FONTS_VERSION}/${FONT_ARCHIVE}"
FONT_DIR="$HOME/.local/share/fonts/devdeck"

require_command() {
  local command_name="$1"

  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Missing required command: $command_name"
    exit 1
  fi
}

require_command curl
require_command unzip
require_command fc-cache

work_dir="$(mktemp -d)"
trap 'rm -rf "$work_dir"' EXIT

mkdir -p "$FONT_DIR"

echo "Installing ${FONT_DISPLAY_NAME} ${NERD_FONTS_VERSION}"
echo "Font target: $FONT_DIR"

curl -L "$FONT_URL" -o "$work_dir/$FONT_ARCHIVE"
unzip -q "$work_dir/$FONT_ARCHIVE" -d "$work_dir/font"

rm -f "$FONT_DIR"/DroidSansMNerdFont*.otf
find "$work_dir/font" -maxdepth 1 -type f -name 'DroidSansMNerdFont*.otf' -exec cp {} "$FONT_DIR" \;

if ! compgen -G "$FONT_DIR/DroidSansMNerdFont*.otf" >/dev/null; then
  echo "No DroidSansM Nerd Font files were installed."
  exit 1
fi

fc-cache -fv "$FONT_DIR"

echo
if command -v fc-match >/dev/null 2>&1 && fc-match "$FONT_DISPLAY_NAME" | grep -Fq "$FONT_DISPLAY_NAME"; then
  echo "Installed: $FONT_DISPLAY_NAME"
else
  echo "Installed font files. If your terminal does not list $FONT_DISPLAY_NAME, restart it."
fi

echo "Select '${FONT_DISPLAY_NAME}' in your terminal profile."
