#!/usr/bin/env bash
set -euo pipefail

NVIM_VERSION="0.12.2"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
NVIM_CONFIG_SOURCE="$REPO_DIR/nvim"
NVIM_CONFIG_TARGET="$HOME/.config/nvim"
NVIM_ARCHIVE="nvim-linux-x86_64.tar.gz"
NVIM_URL="https://github.com/neovim/neovim/releases/download/v${NVIM_VERSION}/${NVIM_ARCHIVE}"

RESET_TARGETS=(
  "$HOME/.local/share/nvim"
  "$HOME/.local/state/nvim"
  "$HOME/.cache/nvim"
)

reset_existing_state() {
  local existing_targets=()
  local target
  local answer

  for target in "${RESET_TARGETS[@]}"; do
    if [ -e "$target" ] || [ -L "$target" ]; then
      existing_targets+=("$target")
    fi
  done

  if [ "${#existing_targets[@]}" -eq 0 ]; then
    return
  fi

  echo
  echo "The following Neovim data paths already exist:"
  for target in "${existing_targets[@]}"; do
    echo "  ${target}"
  done
  echo
  read -r -p "Delete them and start with clean Neovim state? [y/N] " answer

  case "$answer" in
    y|Y|yes|YES)
      rm -rf "${existing_targets[@]}"
      ;;
    *)
      echo "Keeping existing Neovim state."
      ;;
  esac
}

if [ ! -f "$NVIM_CONFIG_SOURCE/init.lua" ]; then
  echo "Neovim config not found: $NVIM_CONFIG_SOURCE"
  exit 1
fi

if [ -e "$NVIM_CONFIG_TARGET" ] && [ ! -L "$NVIM_CONFIG_TARGET" ]; then
  echo "Refusing to overwrite existing path: $NVIM_CONFIG_TARGET"
  echo "Move it aside and run this script again."
  exit 1
fi

echo "Installing Neovim ${NVIM_VERSION}"
echo "Config source: ${NVIM_CONFIG_SOURCE}"

reset_existing_state

rm -rf "$HOME/.local/nvim"
rm -f "$HOME/.local/bin/nvim"

mkdir -p "$HOME/.local" "$HOME/.local/bin" "$HOME/.config"

work_dir="$(mktemp -d)"
trap 'rm -rf "$work_dir"' EXIT

curl -L "$NVIM_URL" -o "$work_dir/$NVIM_ARCHIVE"
tar xzf "$work_dir/$NVIM_ARCHIVE" -C "$work_dir"

mv "$work_dir/nvim-linux-x86_64" "$HOME/.local/nvim"
ln -sfn "$HOME/.local/nvim/bin/nvim" "$HOME/.local/bin/nvim"
ln -sfn "$NVIM_CONFIG_SOURCE" "$NVIM_CONFIG_TARGET"

"$HOME/.local/bin/nvim" --headless "+Lazy! install" "+qa"

echo
echo 'Make sure ~/.local/bin is on your PATH:'
echo 'export PATH="$HOME/.local/bin:$PATH"'
echo
echo "Done. Start Neovim with: nvim"
