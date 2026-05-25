#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

run_step() {
  local label="$1"
  shift

  echo
  echo "==> $label"
  "$@"
}

run_step "Installing tmux" "$SCRIPT_DIR/install-tmux.sh"
run_step "Configuring toolsets" "$SCRIPT_DIR/configure-toolsets.sh"
run_step "Installing Neovim" "$SCRIPT_DIR/install-nvim.sh"

echo
read -r -p "Install Nerd Font? [y/N] " reply || reply=""
case "$reply" in
  [yY] | [yY][eE][sS])
    run_step "Installing font" "$SCRIPT_DIR/install-font.sh"
    ;;
  *)
    echo "Skipping font. Run scripts/install-font.sh later if you want it."
    ;;
esac

echo
echo "Done. Re-pick toolsets anytime with: scripts/configure-toolsets.sh"
