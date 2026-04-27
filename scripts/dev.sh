#!/usr/bin/env sh
set -eu

repo_root="$(git rev-parse --show-toplevel)"

exec go -C "$repo_root" run ./cmds/byob "$@"

