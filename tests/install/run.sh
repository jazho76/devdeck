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
mount_src=/opt/devdeck-src
staged_src=/tmp/devdeck-src

"$engine" build -t "$image_tag" -f "$containerfile" "$repo_root/tests/install/images"
"$engine" run --rm \
  --name "devdeck-install-${image}-${scenario}" \
  -v "$repo_root:$mount_src:ro" \
  -e DEVDECK_SOURCE="$staged_src" \
  -e DEVDECK_TEST_HOME=/home/devdeck-test \
  "$image_tag" \
  bash -euo pipefail -c '
    mount="$1"; staged="$2"; scenario="$3"
    git config --global --add safe.directory "$mount"
    mkdir -p "$staged"
    git -C "$mount" ls-files -z --cached --others --exclude-standard \
      | while IFS= read -r -d "" path; do
          [[ -e "$mount/$path" ]] && printf "%s\0" "$path"
        done \
      | tar -C "$mount" --null -T - -cf - \
      | tar -C "$staged" -xf -
    exec bash "$staged/tests/install/scenarios/$scenario.sh"
  ' _ "$mount_src" "$staged_src" "$scenario"
