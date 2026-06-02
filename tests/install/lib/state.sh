#!/usr/bin/env bash

source "$DEVDECK_SOURCE/tests/install/lib/assert.sh"

assert_devdeck_cli() {
  assert_executable "$HOME/.local/bin/devdeck"
  assert_command_succeeds devdeck version
}

assert_devdeck_source() {
  assert_dir "$HOME/.local/share/devdeck/source"
  assert_file "$HOME/.local/share/devdeck/source/README.md"
  assert_dir "$HOME/.local/share/devdeck/source/nvim"
  assert_dir "$HOME/.local/share/devdeck/source/tmux"
}

assert_tmux_config_link() {
  assert_symlink_to \
    "$HOME/.config/tmux" \
    "$HOME/.local/share/devdeck/source/tmux"
  assert_file "$HOME/.config/tmux/tmux.conf"
}

assert_tpm_installed() {
  assert_dir "$HOME/.local/share/tmux/plugins/tpm"
  assert_executable "$HOME/.local/share/tmux/plugins/tpm/bin/install_plugins"
  assert_not_exists "$HOME/.config/tmux/plugins/tpm"
}

assert_nvim_config_link() {
  assert_symlink_to \
    "$HOME/.config/nvim" \
    "$HOME/.local/share/devdeck/source/nvim"
  assert_file "$HOME/.config/nvim/init.lua"
}

assert_nvim_available() {
  if [[ -x "$HOME/.local/share/devdeck/nvim/bin/nvim" ]]; then
    assert_symlink_to \
      "$HOME/.local/bin/nvim" \
      "$HOME/.local/share/devdeck/nvim/bin/nvim"
    assert_executable "$HOME/.local/bin/nvim"
    return
  fi

  command -v nvim >/dev/null 2>&1 || fail "expected Neovim either managed by Devdeck or available on PATH"
}

assert_nvim_binary_starts() {
  local nvim_bin
  if [[ -x "$HOME/.local/bin/nvim" ]]; then
    nvim_bin="$HOME/.local/bin/nvim"
  else
    nvim_bin="$(command -v nvim)"
  fi

  "$nvim_bin" --headless -u NONE '+qa'
}

assert_default_toolsets() {
  local file="$HOME/.local/share/devdeck/source/nvim/lua/config/toolsets-local.lua"
  assert_file "$file"
  assert_contains "$file" "'c_cpp'"
  assert_contains "$file" "'config'"
  assert_contains "$file" "'markdown'"
}

assert_no_toolsets() {
  local file="$HOME/.local/share/devdeck/source/nvim/lua/config/toolsets-local.lua"
  assert_file "$file"
  assert_not_contains "$file" "'c_cpp'"
  assert_not_contains "$file" "'config'"
  assert_not_contains "$file" "'markdown'"
}
