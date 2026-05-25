#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
TMUX_CONFIG_SOURCE="$REPO_DIR/tmux"
TMUX_CONFIG_TARGET="$HOME/.config/tmux"
TMUX_DATA_DIR="$HOME/.local/share/tmux"
TPM_DIR="$TMUX_DATA_DIR/plugins/tpm"
TPM_URL="https://github.com/tmux-plugins/tpm"

confirm_delete() {
  local target="$1"

  if [ ! -e "$target" ] && [ ! -L "$target" ]; then
    return
  fi

  echo
  read -r -p "Delete existing ${target}? [y/N] " answer

  case "$answer" in
    y|Y|yes|YES)
      rm -rf "$target"
      ;;
    *)
      echo "Keeping ${target}."
      ;;
  esac
}

if [ ! -f "$TMUX_CONFIG_SOURCE/tmux.conf" ]; then
  echo "Tmux config not found: $TMUX_CONFIG_SOURCE"
  exit 1
fi

if [ -e "$TMUX_CONFIG_TARGET" ] && [ ! -L "$TMUX_CONFIG_TARGET" ]; then
  echo "Refusing to overwrite existing path: $TMUX_CONFIG_TARGET"
  echo "Move it aside and run this script again."
  exit 1
fi

echo "Installing tmux config"
echo "Config source: ${TMUX_CONFIG_SOURCE}"

confirm_delete "$TMUX_DATA_DIR"

mkdir -p "$HOME/.config" "$TMUX_DATA_DIR/plugins"
ln -sfn "$TMUX_CONFIG_SOURCE" "$TMUX_CONFIG_TARGET"
echo "Linked tmux config: $TMUX_CONFIG_TARGET -> $TMUX_CONFIG_SOURCE"

if [ -d "$TPM_DIR/.git" ]; then
  git -C "$TPM_DIR" pull --ff-only
elif [ -e "$TPM_DIR" ]; then
  echo "Refusing to overwrite existing path: $TPM_DIR"
  echo "Move it aside and run this script again."
  exit 1
else
  git clone "$TPM_URL" "$TPM_DIR"
fi

"$TPM_DIR/bin/install_plugins"

echo "TPM ready: $TPM_DIR"
echo "Done. Start tmux with: tmux"
