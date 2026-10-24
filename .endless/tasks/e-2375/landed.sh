#!/bin/sh
# ── DO NOT EDIT ─────────────────────────────────────────────────────
# This suite belongs to E-2375 and records what was true when E-2375
# landed. Edit it only if you ARE E-2375. If your change breaks an
# assertion here, leave it alone — see .endless/tasks/CLAUDE.md.
#
# Proves the kept changes from main's working tree landed, the dropped
# one did not, and every module builds and passes its tests.
# Runs from the repo root of the worktree under test.

n=0
fails=0
ok() { n=$((n + 1)); echo "ok $n - $1"; }
not_ok() { n=$((n + 1)); fails=$((fails + 1)); echo "not ok $n - $1"; }
check() { desc=$1; shift; if "$@" >/dev/null 2>&1; then ok "$desc"; else not_ok "$desc"; fi; }

cd "$(git rev-parse --show-toplevel)" || { echo "Bail out! not in a git repo"; exit 2; }
export GOEXPERIMENT=jsonv2

echo "TAP version 13"

# Fail fast: the task's own unit test.
if ! go test -count=1 -run '^TestPathSegments_IsRoot$' . >/dev/null 2>&1; then
    echo "not ok 1 - TestPathSegments_IsRoot passes"
    echo "1..1"
    exit 1
fi
ok "TestPathSegments_IsRoot passes"

check "dtx/test requires go-dt v0.5.0" grep -qE '^[[:space:]]+github.com/mikeschinkel/go-dt v0\.5\.0$' dtx/test/go.mod
check "dtx/test requires dtx v0.2.1" grep -qE '^[[:space:]]+github.com/mikeschinkel/go-dt/dtx v0\.2\.1$' dtx/test/go.mod
check "dtx/test has no zero pseudo-version" sh -c '! grep -q v0.0.0-00010101000000 dtx/test/go.mod'
check "dropped dtx/go.sum is not tracked" sh -c '! git ls-files --error-unmatch dtx/go.sum'

for d in . appinfo appinfo/test dtglob dtglob/test dtx dtx/test examples/basic_usage test; do
    check "$d builds and vets" sh -c "cd '$d' && go vet ./..."
done
check "make test passes (root and all test modules)" make test

echo "1..$n"
[ "$fails" -eq 0 ]
