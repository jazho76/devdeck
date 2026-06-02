#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: tests/install/matrix.sh <fast|full>

fast: Minimal Fedora install contract scenarios intended for PR CI.
full: Currently aliases fast.
USAGE
}

if [[ $# -ne 1 ]]; then
  usage >&2
  exit 2
fi

suite="$1"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

case "$suite" in
  fast|full)
    scenarios=(
      "fedora fresh-minimal"
      "fedora unmanaged-nvim-config"
      "fedora unmanaged-tmux-config"
    )
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac

for entry in "${scenarios[@]}"; do
  image="${entry%% *}"
  scenario="${entry#* }"
  printf '\n==> %s / %s\n' "$image" "$scenario"
  "$root/tests/install/run.sh" "$image" "$scenario"
done
