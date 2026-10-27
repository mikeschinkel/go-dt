#!/usr/bin/env bash
# ── DO NOT EDIT ─────────────────────────────────────────────────────
# This suite belongs to E-2452 and records what was true when E-2452
# landed. Edit it only if you ARE E-2452. If your change breaks an
# assertion here, leave it alone — see .endless/tasks/CLAUDE.md.
source "$(dirname "${BASH_SOURCE[0]}")/../_harness.sh"

set -u

ROOT="$(git -C "${ENDLESS_VERIFY_DIR}" rev-parse --show-toplevel)" \
    || setup_error "cannot find the repository root"
cd "${ROOT}" || setup_error "cannot cd to ${ROOT}"

RELEASE_YML=.github/workflows/release.yml
BASE="$(git merge-base main HEAD)" || setup_error "cannot find the merge base with main"

# The 'Generate release notes' step's script, exactly as release.yml has it
step=$(yq '.jobs[].steps[] | select(.name == "Generate release notes") | .run' "${RELEASE_YML}") \
    || setup_error "yq cannot read ${RELEASE_YML}"
[[ -n "${step}" ]] || setup_error "release.yml has no 'Generate release notes' step"

section "Fail fast: the release-notes step is valid bash"
if out=$(bash -n <<<"${step}" 2>&1); then
    report_pass "the 'Generate release notes' script parses"
else
    report_fail "the 'Generate release notes' script parses" "no syntax error" "${out}"
    summary
fi

# run_step runs release.yml's step for one module over prev..tag, in a scratch
# clone checked out at tag (the release workflow runs at the tag's commit), and
# prints what it writes to GITHUB_OUTPUT. GitHub's ${{ }} expressions are
# replaced with this repository's values.
scratch=$(mktemp -d) || setup_error "cannot make a scratch directory"
trap 'rm -rf "${scratch}"' EXIT
git clone -q "${ROOT}" "${scratch}/repo" || setup_error "cannot clone ${ROOT}"
script=${step//'${{ github.repository }}'/mikeschinkel/go-dt}
script=${script//'${{ github.event.repository.name }}'/go-dt}
run_step() {
    local module="$1" dir="$2" prev="$3" tag="$4" out="${scratch}/out"
    : >"${out}"
    git -C "${scratch}/repo" checkout -q --detach "${tag}" || setup_error "cannot check out ${tag}"
    (cd "${scratch}/repo" \
        && MODULE="${module}" DIR="${dir}" VERSION="${tag##*/}" TAG="${tag}" \
           TITLE="go-dt/${dir} ${tag##*/}" PREV_TAG="${prev}" DESCRIPTION="" \
           GITHUB_OUTPUT="${out}" bash -e <<<"${script}") \
        || setup_error "the release-notes step failed for ${module} ${tag}"
    cat "${out}"
}

# listed prints the short hashes of the commits a body lists, one per line
listed() {
    sed -n 's/^- \[\*\*\([0-9a-f]*\)\*\*\].*/\1/p' <<<"$1"
}

# only_endless succeeds when every file a commit changes is under .endless/
only_endless() {
    ! git show --name-only --format= "$1" | grep -qv '^\.endless/'
}

section "Root's notes leave out commits that touch only .endless/ (P1)"
assert_contains "root's PATHS excludes .endless after the nested modules" \
    "PATHS=(-- . ':(exclude)dtx' ':(exclude)appinfo' ':(exclude)dtglob' ':(exclude).endless')" "${step}"

section "Root's v0.6.0..v0.6.1 notes, made by release.yml's step (P2, P3)"
body=$(run_step root . v0.6.0 v0.6.1)
hashes=$(listed "${body}")
endless_only=$(for h in $(git log --format=%h v0.6.0..v0.6.1); do
    only_endless "${h}" && echo "${h}"
done)
[[ -n "${endless_only}" ]] || setup_error "v0.6.0..v0.6.1 has no .endless-only commit to leave out"
leaked=""
for h in ${hashes}; do
    only_endless "${h}" && leaked+="${h} "
done
assert_eq "no listed commit touches only .endless/ ($(wc -l <<<"${endless_only}" | tr -d ' ') left out)" "" "${leaked}"
assert_contains "the PathSegments.IsRoot commit is listed" "be26016" "${hashes}"
# 2cab020 changes .endless/ and files outside it, so it stays
files=$(git show --name-only --format= 2cab020)
assert_eq "2cab020 changes files both inside and outside .endless/" "yes yes" \
    "$(grep -q '^\.endless/' <<<"${files}" && echo yes) $(grep -qv '^\.endless/' <<<"${files}" && echo yes)"
assert_contains "the mixed commit 2cab020 is still listed" "2cab020" "${hashes}"
assert_eq "root's notes list be26016, 2cab020 and b9a036c" \
    "be26016 2cab020 b9a036c" "$(tr '\n' ' ' <<<"${hashes}" | sed 's/ $//')"
assert_contains "the heading names the project and version" "## Domain Types for Go (`dt`) v0.6.1" "${body}"
assert_contains "the Full Changelog line is kept" \
    "**Full Changelog**: https://github.com/mikeschinkel/go-dt/compare/v0.6.0...v0.6.1" "${body}"

section "Root's v0.5.0..v0.6.0 notes are unchanged by the exclusion"
assert_eq "root's v0.5.0..v0.6.0 notes still list 14 commits" "14" \
    "$(listed "$(run_step root . v0.5.0 v0.6.0)" | wc -l | tr -d ' ')"

section "Nested modules' release-notes PATHS are unchanged (P4)"
assert_contains "a nested module's PATHS is still its own directory" 'PATHS=(-- "$DIR")' "${step}"
changed=$(git diff -U0 "${BASE}" HEAD -- "${RELEASE_YML}" | grep '^[-+] ' | grep -v '^[-+] *#')
assert_eq "the only code line release.yml changes is root's PATHS" \
    "-            PATHS=(-- . ':(exclude)dtx' ':(exclude)appinfo' ':(exclude)dtglob')
+            PATHS=(-- . ':(exclude)dtx' ':(exclude)appinfo' ':(exclude)dtglob' ':(exclude).endless')" \
    "${changed}"

section "Scope: no other workflow, Go code or module file changes (P6, P7)"
assert_eq "outside .endless/, only release.yml differs from main" "${RELEASE_YML}" \
    "$(git diff --name-only "${BASE}" HEAD -- . ':(exclude).endless')"

summary
