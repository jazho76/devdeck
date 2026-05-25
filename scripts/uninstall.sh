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

run_step "Uninstalling Neovim" "$SCRIPT_DIR/uninstall-nvim.sh"
run_step "Uninstalling tmux" "$SCRIPT_DIR/uninstall-tmux.sh"

echo
echo "Done."
