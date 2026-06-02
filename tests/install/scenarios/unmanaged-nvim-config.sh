#!/usr/bin/env bash
set -euo pipefail

source /opt/devdeck-src/tests/install/lib/assert.sh
source /opt/devdeck-src/tests/install/lib/env.sh
source /opt/devdeck-src/tests/install/lib/state.sh

setup_test_home
install_devdeck_cli_from_source

note "Creating unmanaged Neovim config"
mkdir -p "$HOME/.config/nvim"
printf 'unmanaged\n' > "$HOME/.config/nvim/sentinel"

note "Install should fail closed instead of overwriting unmanaged Neovim config"
assert_command_fails devdeck install --toolsets c_cpp,config,markdown --skip-lazy-install
assert_file "$HOME/.config/nvim/sentinel"
assert_contains "$HOME/.config/nvim/sentinel" "unmanaged"
assert_not_exists "$HOME/.config/nvim/init.lua"
