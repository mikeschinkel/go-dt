#!/usr/bin/env bash
# ── DO NOT EDIT ─────────────────────────────────────────────────────
# This suite belongs to E-2431 and records what was true when E-2431
# landed. Edit it only if you ARE E-2431. If your change breaks an
# assertion here, leave it alone — see .endless/tasks/CLAUDE.md.
source "$(dirname "${BASH_SOURCE[0]}")/../_harness.sh"

set -u

# Every module is checked as CI sees it: against its go.mod's requires, with
# no go.work. A gotest check cannot reach a nested module (endless E-2438),
# so each one is a `cd <dir> && go ...`.
export GOWORK=off
export GOEXPERIMENT=jsonv2

ROOT="$(git -C "${ENDLESS_VERIFY_DIR}" rev-parse --show-toplevel)" \
    || setup_error "cannot find the repository root"
cd "${ROOT}" || setup_error "cannot cd to ${ROOT}"

PUBLISHABLE=(dtx appinfo dtglob)
ALL_MODULES=(. test dtx dtx/test appinfo appinfo/test dtglob dtglob/test examples/basic_usage)

# replaces prints one "old => new" line per replace directive in a go.mod
replaces() {
    go mod edit -json "$1/go.mod" | jq -r '.Replace // [] | .[] | "\(.Old.Path) => \(.New.Path)"'
}

section "Fail fast: every module vets and tests against its requires"
for dir in "${ALL_MODULES[@]}"; do
    if out=$(cd "${dir}" && go vet ./... 2>&1 && go test ./... 2>&1); then
        report_pass "${dir}: go vet and go test pass with GOWORK=off"
    else
        report_fail "${dir}: go vet and go test pass with GOWORK=off" "exit 0" "${out}"
    fi
done
if (( FAIL_COUNT > 0 )); then
    summary
fi

section "Publishable modules require a published go-dt and carry no replace (P1, P2)"
for dir in "${PUBLISHABLE[@]}"; do
    assert_eq "${dir}/go.mod has no replace directive" "" "$(replaces "${dir}")"
    required=$(go mod edit -json "${dir}/go.mod" \
        | jq -r '.Require[] | select(.Path == "github.com/mikeschinkel/go-dt") | .Version')
    newest=$(printf '%s\n%s\n' "v0.6.1" "${required}" | sort -V | tail -n 1)
    if [[ -n "${required}" && "${newest}" == "${required}" ]]; then
        report_pass "${dir} requires go-dt ${required} (>= v0.6.1)"
    else
        report_fail "${dir} requires go-dt >= v0.6.1" ">= v0.6.1" "${required:-none}"
    fi
    assert_contains "${dir}/go.sum records go-dt ${required}" \
        "github.com/mikeschinkel/go-dt ${required} h1:" "$(cat "${dir}/go.sum" 2>/dev/null)"
done

section "Test and example modules keep one replace, to the module they test (P3)"
assert_eq "test replaces go-dt only" \
    "github.com/mikeschinkel/go-dt => .." "$(replaces test)"
assert_eq "examples/basic_usage replaces go-dt only" \
    "github.com/mikeschinkel/go-dt => ../.." "$(replaces examples/basic_usage)"
for dir in "${PUBLISHABLE[@]}"; do
    assert_eq "${dir}/test replaces ${dir} only" \
        "github.com/mikeschinkel/go-dt/${dir} => .." "$(replaces "${dir}/test")"
done

section "Every go.mod is tidy as CI sees it (P4)"
for dir in "${ALL_MODULES[@]}"; do
    if out=$(cd "${dir}" && go mod tidy -diff 2>&1); then
        report_pass "${dir}: go mod tidy -diff reports nothing"
    else
        report_fail "${dir}: go mod tidy -diff reports nothing" "no diff" "${out}"
    fi
done

section "A justfile replaces the Makefile (P5)"
assert_eq "the Makefile is gone" "no" "$([[ -e Makefile ]] && echo yes || echo no)"
recipes=$(just --summary 2>&1)
for recipe in help test test-unit test-corpus test-all lint fmt vet tidy build examples clean ci ensure-valid work; do
    assert_contains "just offers ${recipe}" " ${recipe} " " ${recipes} "
done

section "just work writes a go.work over every module (P6)"
scratch=$(mktemp -d) || setup_error "cannot make a scratch directory"
git archive HEAD | tar -x -C "${scratch}" || setup_error "cannot export HEAD"
# A module just work has never heard of must still land in go.work
mkdir -p "${scratch}/newmod" && printf 'module example.com/newmod\n\ngo 1.25.3\n' >"${scratch}/newmod/go.mod"
if out=$(cd "${scratch}" && GOWORK= just work 2>&1); then
    listed=$(cd "${scratch}" && GOWORK= go work edit -json | jq -r '.Use[].DiskPath' | sed 's|^\./||' | sort)
    expected=$(cd "${scratch}" && find . -name go.mod | sed 's|^\./||; s|/\{0,1\}go.mod$||; s|^$|.|' | sort)
    assert_eq "go.work uses every directory holding a go.mod" "${expected}" "${listed}"
else
    report_fail "just work succeeds" "exit 0" "${out}"
fi
rm -rf "${scratch}"

section "README documents just work for local development (P7)"
readme=$(cat README.md)
assert_contains "README has a Local Development section" "## Local Development" "${readme}"
assert_contains "README tells you to run just work" "just work" "${readme}"
assert_contains "README explains GOWORK=off" "GOWORK=off" "${readme}"

section "test.yml runs without go.work and covers every module (P8)"
testyml=$(cat .github/workflows/test.yml)
assert_contains "test.yml sets GOWORK off" "GOWORK: 'off'" "${testyml}"
for dir in "${PUBLISHABLE[@]}"; do
    assert_contains "test.yml tests ${dir}" "(cd ${dir} && go test" "${testyml}"
    assert_contains "test.yml tests ${dir}/test" "(cd ${dir}/test && go test" "${testyml}"
done
assert_contains "test.yml checks every go.mod is tidy" "go mod tidy -diff" "${testyml}"

section "release.yml makes directory-prefixed tags safely (P9–P13)"
rel=$(cat .github/workflows/release.yml)
# The option lines between the module input's `options:` and the next input
modules_choice=$(sed -n '/^      module:/,/^      version:/p' .github/workflows/release.yml \
    | sed -n 's/^          - //p' | tr '\n' ' ')
assert_contains "module is a choice input" "type: choice" "${rel}"
assert_eq "module is a choice of root, dtx, appinfo, dtglob" "root dtx appinfo dtglob " "${modules_choice}"
assert_contains "root tags vX.Y.Z" 'TAG="$VERSION"' "${rel}"
assert_contains "a nested module tags <dir>/vX.Y.Z" 'TAG="$MODULE/$VERSION"' "${rel}"
assert_contains "release.yml checks the declared module path" '.Module.Path' "${rel}"
assert_contains "release.yml refuses a replace directive" 'contains replace directives' "${rel}"
assert_contains "release.yml tests with GOWORK off" "GOWORK: 'off'" "${rel}"
assert_contains "release.yml tests the module's test module" 'working-directory: ${{ env.DIR }}/test' "${rel}"
assert_contains "a nested module's previous tag has its prefix" 'git tag -l "$MODULE/v*"' "${rel}"
assert_contains "root's previous tag is unprefixed" "git tag -l 'v*'" "${rel}"
assert_contains "nested notes list only commits touching its directory" 'PATHS=(-- "$DIR")' "${rel}"
assert_contains "root notes omit commits touching only nested modules" \
    "':(exclude)dtx' ':(exclude)appinfo' ':(exclude)dtglob'" "${rel}"
assert_contains "the release title names the module" 'name: ${{ env.TITLE }}' "${rel}"

# Run the tag and previous-tag logic as release.yml would, on this repo's tags
for module in root dtx appinfo dtglob; do
    if [[ "${module}" == root ]]; then
        prev=$(git tag -l 'v*' --sort=-v:refname | head -n 1)
        assert_not_contains "root's previous tag (${prev}) has no prefix" "/" "${prev}"
    else
        prev=$(git tag -l "${module}/v*" --sort=-v:refname | head -n 1)
        assert_contains "${module}'s previous tag (${prev}) has its prefix" "${module}/v" "${prev}"
    fi
done

section "Scope: root's Go code is unchanged (P18)"
assert_eq "no root-level .go file differs from main" "" \
    "$(git diff --name-only "$(git merge-base main HEAD)" HEAD -- ':(glob)*.go')"

section "go-dt carries Endless's harness unchanged (P21)"
for f in _harness.sh _guard.sh; do
    assert_eq ".endless/tasks/${f} is present" "yes" "$([[ -f .endless/tasks/${f} ]] && echo yes || echo no)"
done

summary
