#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

all_flag=""
if [ "${1:-}" = "--all" ]; then
  all_flag="--all"
fi

run_step() {
  local label="$1"
  shift

  echo
  echo "==> $label"
  "$@"
}

run_step "Installing tmux" "$SCRIPT_DIR/_install-tmux.sh"
run_step "Configuring toolsets" "$SCRIPT_DIR/_configure-toolsets.py" $all_flag
run_step "Installing Neovim" "$SCRIPT_DIR/_install-nvim.sh"

echo
echo "Done. Re-pick toolsets anytime with: scripts/_configure-toolsets.py"
