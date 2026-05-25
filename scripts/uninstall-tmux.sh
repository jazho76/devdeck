#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
TMUX_CONFIG_SOURCE="$REPO_DIR/tmux"
TMUX_CONFIG_TARGET="$HOME/.config/tmux"
TMUX_PLUGIN_DIR="$HOME/.local/share/tmux/plugins"

remove_if_symlink_to() {
  local link_path="$1"
  local expected_target="$2"

  if [ ! -L "$link_path" ]; then
    if [ -e "$link_path" ]; then
      echo "Keeping unmanaged path: $link_path"
    fi
    return
  fi

  actual_target="$(readlink -f "$link_path")"
  expected_target="$(readlink -f "$expected_target")"

  if [ "$actual_target" = "$expected_target" ]; then
    rm "$link_path"
    echo "Removed symlink: $link_path"
  else
    echo "Keeping unmanaged symlink: $link_path -> $actual_target"
  fi
}

remove_if_symlink_to "$TMUX_CONFIG_TARGET" "$TMUX_CONFIG_SOURCE"

if [ -d "$TMUX_PLUGIN_DIR" ]; then
  rm -rf "$TMUX_PLUGIN_DIR"
  echo "Removed tmux plugins: $TMUX_PLUGIN_DIR"
else
  echo "No tmux plugin directory found: $TMUX_PLUGIN_DIR"
fi

echo "Kept tmux binary untouched."
