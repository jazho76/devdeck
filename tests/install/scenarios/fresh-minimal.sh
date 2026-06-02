#!/usr/bin/env bash
set -euo pipefail

source /opt/devdeck-src/tests/install/lib/assert.sh
source /opt/devdeck-src/tests/install/lib/env.sh
source /opt/devdeck-src/tests/install/lib/state.sh

setup_test_home
install_devdeck_cli_from_source

note "Installing with default minimal toolsets and no Lazy bootstrap"
devdeck install --toolsets c_cpp,config,markdown --skip-lazy-install

assert_devdeck_cli
assert_devdeck_source
assert_tmux_config_link
assert_tpm_installed
assert_nvim_config_link
assert_nvim_available
assert_default_toolsets
assert_nvim_binary_starts
assert_command_succeeds devdeck doctor
assert_command_succeeds devdeck workspace --help
assert_command_succeeds devdeck ws --help
assert_command_succeeds devdeck upgrade --help
