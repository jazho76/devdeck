#!/usr/bin/env bash
set -euo pipefail

NVIM_VERSION="0.12.2"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
NVIM_CONFIG_DIR="$REPO_DIR/nvim"

NVIM_ARCHIVE="nvim-linux-x86_64.tar.gz"
NVIM_URL="https://github.com/neovim/neovim/releases/download/v${NVIM_VERSION}/${NVIM_ARCHIVE}"

TARGETS=(
  "$HOME/.local/nvim"
  "$HOME/.local/bin/nvim"
  "$HOME/.config/nvim"
  "$HOME/.local/share/nvim"
  "$HOME/.local/state/nvim"
  "$HOME/.cache/nvim"
)

if [ ! -f "$NVIM_CONFIG_DIR/init.lua" ]; then
  echo "Neovim config not found: $NVIM_CONFIG_DIR"
  exit 1
fi

existing_targets=()
for target in "${TARGETS[@]}"; do
  if [ -e "$target" ] || [ -L "$target" ]; then
    existing_targets+=("$target")
  fi
done

echo "Installing Neovim ${NVIM_VERSION}"
echo "Config source: ${NVIM_CONFIG_DIR}"

if [ "${#existing_targets[@]}" -gt 0 ]; then
  echo
  echo "The following paths already exist and will be deleted:"
  for target in "${existing_targets[@]}"; do
    echo "  ${target}"
  done
  echo
  read -r -p "Continue? [y/N] " answer

  case "$answer" in
    y|Y|yes|YES)
      ;;
    *)
      echo "Aborted."
      exit 1
      ;;
  esac
fi

rm -rf "$HOME/.local/nvim"
rm -f "$HOME/.local/bin/nvim"
rm -rf "$HOME/.config/nvim"
rm -rf "$HOME/.local/share/nvim"
rm -rf "$HOME/.local/state/nvim"
rm -rf "$HOME/.cache/nvim"

mkdir -p "$HOME/.local" "$HOME/.local/bin" "$HOME/.config"

curl -L "$NVIM_URL" -o "$NVIM_ARCHIVE"
tar xzf "$NVIM_ARCHIVE"
rm "$NVIM_ARCHIVE"

mv nvim-linux-x86_64 "$HOME/.local/nvim"
ln -sfn "$HOME/.local/nvim/bin/nvim" "$HOME/.local/bin/nvim"
ln -sfn "$NVIM_CONFIG_DIR" "$HOME/.config/nvim"

echo
echo 'Make sure ~/.local/bin is on your PATH:'
echo 'export PATH="$HOME/.local/bin:$PATH"'
echo
echo "Done. Start Neovim with: nvim"
