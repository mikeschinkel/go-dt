#!/bin/sh
# Part of E-2567's verify suite. DO NOT EDIT unless you ARE E-2567.
# Prints every exported declaration of the root package, as go doc shows it,
# that existed where this branch left main and no longer exists unchanged.
# Empty output means no existing exported symbol was renamed, removed or had
# its signature changed.
set -eu
root=$(git rev-parse --show-toplevel)
base=$(git merge-base HEAD main)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
git -C "$root" archive "$base" | tar -x -C "$tmp"
export GOWORK=off
decls() {
	(cd "$1" && go doc -all .) |
		grep -E '^(func|type) |^[[:space:]]+[A-Z][A-Za-z0-9_]* +(=|[A-Z][A-Za-z0-9_]* +=)' |
		tr -s ' \t' ' ' | sed 's/^ //' | sort -u
}
decls "$tmp" >"$tmp/.before"
decls "$root" >"$tmp/.after"
comm -23 "$tmp/.before" "$tmp/.after"
