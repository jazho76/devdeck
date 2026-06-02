#!/usr/bin/env bash

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

note() {
  printf '\n## %s\n' "$*"
}

assert_exists() {
  local path="$1"
  [[ -e "$path" ]] || fail "expected path to exist: $path"
}

assert_not_exists() {
  local path="$1"
  [[ ! -e "$path" ]] || fail "expected path to be absent: $path"
}

assert_dir() {
  local path="$1"
  [[ -d "$path" ]] || fail "expected directory: $path"
}

assert_file() {
  local path="$1"
  [[ -f "$path" ]] || fail "expected file: $path"
}

assert_executable() {
  local path="$1"
  [[ -x "$path" ]] || fail "expected executable: $path"
}

assert_symlink_to() {
  local link="$1"
  local expected="$2"
  [[ -L "$link" ]] || fail "expected symlink: $link"

  local actual
  actual="$(readlink "$link")"
  [[ "$actual" == "$expected" ]] || fail "expected $link -> $expected, got $actual"
}

assert_contains() {
  local path="$1"
  local needle="$2"
  grep -Fq -- "$needle" "$path" || fail "expected $path to contain: $needle"
}

assert_not_contains() {
  local path="$1"
  local needle="$2"
  if grep -Fq -- "$needle" "$path"; then
    fail "expected $path not to contain: $needle"
  fi
}

assert_command_succeeds() {
  "$@" || fail "command failed: $*"
}

assert_command_fails() {
  if "$@"; then
    fail "expected command to fail: $*"
  fi
}
