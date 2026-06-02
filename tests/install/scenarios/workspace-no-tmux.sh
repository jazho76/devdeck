#!/usr/bin/env bash
set -euo pipefail

source /opt/devdeck-src/tests/install/lib/assert.sh
source /opt/devdeck-src/tests/install/lib/env.sh
source /opt/devdeck-src/tests/install/lib/state.sh

setup_test_home
install_devdeck_cli_from_source

note "Workspace commands outside tmux"
assert_command_succeeds devdeck workspace list
assert_command_fails devdeck workspace save outside-tmux
assert_command_fails devdeck workspace restore outside-tmux
assert_not_exists "$HOME/.local/share/devdeck/workspaces/outside-tmux.json"
