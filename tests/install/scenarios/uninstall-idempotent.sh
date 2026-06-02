#!/usr/bin/env bash
set -euo pipefail

source /opt/devdeck-src/tests/install/lib/assert.sh
source /opt/devdeck-src/tests/install/lib/env.sh
source /opt/devdeck-src/tests/install/lib/state.sh

setup_test_home
install_devdeck_cli_from_source

devdeck install --toolsets c_cpp,config,markdown --skip-lazy-install
mkdir -p "$HOME/.local/share/nvim" "$HOME/.local/state/nvim" "$HOME/.cache/nvim"
printf 'runtime\n' > "$HOME/.local/share/nvim/sentinel"
printf 'state\n' > "$HOME/.local/state/nvim/sentinel"
printf 'cache\n' > "$HOME/.cache/nvim/sentinel"

note "First uninstall"
devdeck uninstall

assert_not_exists "$HOME/.config/nvim"
assert_not_exists "$HOME/.config/tmux"
assert_not_exists "$HOME/.local/share/devdeck/source"
assert_not_exists "$HOME/.local/share/tmux/plugins"
assert_file "$HOME/.local/share/nvim/sentinel"
assert_file "$HOME/.local/state/nvim/sentinel"
assert_file "$HOME/.cache/nvim/sentinel"

note "Second uninstall should be idempotent"
devdeck uninstall
assert_file "$HOME/.local/share/nvim/sentinel"
assert_file "$HOME/.local/state/nvim/sentinel"
assert_file "$HOME/.cache/nvim/sentinel"
