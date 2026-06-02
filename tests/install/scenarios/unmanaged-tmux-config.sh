#!/usr/bin/env bash
set -euo pipefail

source /opt/devdeck-src/tests/install/lib/assert.sh
source /opt/devdeck-src/tests/install/lib/env.sh
source /opt/devdeck-src/tests/install/lib/state.sh

setup_test_home
install_devdeck_cli_from_source

note "Creating unmanaged tmux config"
mkdir -p "$HOME/.config/tmux"
printf 'unmanaged\n' > "$HOME/.config/tmux/tmux.conf"

note "Install should fail closed instead of overwriting unmanaged tmux config"
assert_command_fails devdeck install --toolsets c_cpp,config,markdown --skip-lazy-install
assert_file "$HOME/.config/tmux/tmux.conf"
assert_contains "$HOME/.config/tmux/tmux.conf" "unmanaged"
assert_not_exists "$HOME/.config/tmux/plugins/tpm"
