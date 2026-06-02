#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: tests/install/run.sh <image> <scenario>

Examples:
  tests/install/run.sh fedora fresh-minimal
  tests/install/run.sh fedora reinstall-idempotent
USAGE
}

if [[ $# -ne 2 ]]; then
  usage >&2
  exit 2
fi

image="$1"
scenario="$2"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
containerfile="$repo_root/tests/install/images/Containerfile.$image"
scenario_file="$repo_root/tests/install/scenarios/$scenario.sh"
engine="${CONTAINER_ENGINE:-}"

if [[ -z "$engine" ]]; then
  if command -v docker >/dev/null 2>&1; then
    engine=docker
  elif command -v podman >/dev/null 2>&1; then
    engine=podman
  else
    echo "Need docker or podman" >&2
    exit 1
  fi
fi

[[ -f "$containerfile" ]] || { echo "No image file: $containerfile" >&2; exit 1; }
[[ -f "$scenario_file" ]] || { echo "No scenario: $scenario_file" >&2; exit 1; }

image_tag="devdeck-install-test:$image"

"$engine" build -t "$image_tag" -f "$containerfile" "$repo_root/tests/install/images"
"$engine" run --rm \
  --name "devdeck-install-${image}-${scenario}" \
  -v "$repo_root:/opt/devdeck-src:ro" \
  -e DEVDECK_SOURCE=/opt/devdeck-src \
  -e DEVDECK_TEST_HOME=/home/devdeck-test \
  "$image_tag" \
  bash "/opt/devdeck-src/tests/install/scenarios/$scenario.sh"
