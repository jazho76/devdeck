#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FONT_CHOICE="ask"
TOOLSET_ARGS=()

usage() {
  cat <<'USAGE'
Usage: scripts/install.sh [--toolsets name[,name...]] [--font|--no-font] [--help]

Options:
  --toolsets LIST  Configure toolsets without opening the fzf picker.
  --font           Install the bundled Nerd Font without prompting.
  --no-font        Skip Nerd Font installation without prompting.
  --help           Show this help.
USAGE
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --toolsets)
      if [ "$#" -lt 2 ]; then
        echo "Missing value for --toolsets" >&2
        exit 1
      fi
      TOOLSET_ARGS=(--toolsets "$2")
      shift 2
      ;;
    --toolsets=*)
      TOOLSET_ARGS=(--toolsets "${1#--toolsets=}")
      shift
      ;;
    --font)
      FONT_CHOICE="install"
      shift
      ;;
    --no-font)
      FONT_CHOICE="skip"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

run_step() {
  local label="$1"
  shift

  echo
  echo "==> $label"
  "$@"
}

run_step "Installing tmux" "$SCRIPT_DIR/install-tmux.sh"
run_step "Configuring toolsets" "$SCRIPT_DIR/configure-toolsets.sh" "${TOOLSET_ARGS[@]}"
run_step "Installing Neovim" "$SCRIPT_DIR/install-nvim.sh"

echo
case "$FONT_CHOICE" in
  install)
    run_step "Installing font" "$SCRIPT_DIR/install-font.sh"
    ;;
  skip)
    echo "Skipping font. Run scripts/install-font.sh later if you want it."
    ;;
  ask)
    read -r -p "Install Nerd Font? [y/N] " reply || reply=""
    case "$reply" in
      [yY] | [yY][eE][sS])
        run_step "Installing font" "$SCRIPT_DIR/install-font.sh"
        ;;
      *)
        echo "Skipping font. Run scripts/install-font.sh later if you want it."
        ;;
    esac
    ;;
esac

echo
echo "Done. Re-pick toolsets anytime with: scripts/configure-toolsets.sh"
