#!/usr/bin/env bash

setup_test_home() {
  export HOME="${DEVDECK_TEST_HOME:-/home/devdeck-test}"
  export XDG_CONFIG_HOME="$HOME/.config"
  export XDG_DATA_HOME="$HOME/.local/share"
  export XDG_STATE_HOME="$HOME/.local/state"
  export XDG_CACHE_HOME="$HOME/.cache"
  export PATH="$HOME/.local/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
  export DEVDECK_SOURCE="${DEVDECK_SOURCE:-/opt/devdeck-src}"
  export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"

  mkdir -p \
    "$XDG_CONFIG_HOME" \
    "$XDG_DATA_HOME" \
    "$XDG_STATE_HOME" \
    "$XDG_CACHE_HOME" \
    "$HOME/.local/bin"
}

install_devdeck_cli_from_source() {
  note "Building devdeck CLI from mounted source"
  (cd "$DEVDECK_SOURCE/cli" && go build -buildvcs=false -o "$HOME/.local/bin/devdeck" .)
  assert_executable "$HOME/.local/bin/devdeck"
}

reset_test_home() {
  rm -rf "$HOME"
  setup_test_home
}
