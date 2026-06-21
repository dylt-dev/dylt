#! /usr/bin/env bash

SCRIPT_DIR=$(dirname "$(readlink -f "$BASH_SOURCE")")
source "$SCRIPT_DIR/test-utils.sh" || exit 1
REPO_DIR=$(dirname "$SCRIPT_DIR")


#-------------------------------------------------------------------------------
#
# run-tests()
#
run-tests()
{
    local tests=(
        test-goreleaser-yaml-core
        test-goreleaser-full-yaml-exists
        test-goreleaser-full-yaml-content
        test-workflow-goreleaser-trigger
        test-workflow-release-config
        test-workflow-nightly-format
        test-workflow-nightly-tag-format
    )
    local total=0 passed=0 failed=0
    for t in "${tests[@]}"; do
        printf '=== %s ===\n' "$t"
        (( total++ ))
        if "$t"; then
            (( passed++ ))
        else
            (( failed++ ))
        fi
    done
    printf '\n%d total, %d passed, %d failed\n' "$total" "$passed" "$failed"
    return "$failed"
}


#-------------------------------------------------------------------------------
#
# test-goreleaser-yaml-core()
#
# Verify .goreleaser.yaml builds only amd64 across darwin/linux/windows
#
test-goreleaser-yaml-core()
{
    local f="$REPO_DIR/.goreleaser.yaml"
    [[ -f "$f" ]] || { printf '  FAIL: %s not found\n' "$f"; return 1; }

    grep -qF '  - amd64' "$f" || { printf '  FAIL: missing amd64 goarch\n'; return 1; }
    grep -qF '  - darwin' "$f" || { printf '  FAIL: missing darwin goos\n'; return 1; }
    grep -qF '  - linux' "$f"  || { printf '  FAIL: missing linux goos\n';  return 1; }
    grep -qF '  - windows' "$f" || { printf '  FAIL: missing windows goos\n'; return 1; }

    grep -qF 'arm64'  "$f" && { printf '  FAIL: arm64 should not be in core config\n'; return 1; }
    grep -qF '  - arm' "$f" && { printf '  FAIL: arm should not be in core config\n'; return 1; }
    grep -qF '  - "386"' "$f" && { printf '  FAIL: 386 should not be in core config\n'; return 1; }
    grep -qF 'goarm:' "$f" && { printf '  FAIL: goarm should not be in core config\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-goreleaser-full-yaml-exists()
#
# Verify .goreleaser.full.yaml is present
#
test-goreleaser-full-yaml-exists()
{
    [[ -f "$REPO_DIR/.goreleaser.full.yaml" ]] || {
        printf '  FAIL: .goreleaser.full.yaml not found\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-goreleaser-full-yaml-content()
#
# Verify .goreleaser.full.yaml has full matrix + windows/arm ignore
#
test-goreleaser-full-yaml-content()
{
    local f="$REPO_DIR/.goreleaser.full.yaml"
    grep -qF 'arm64'      "$f" || { printf '  FAIL: missing arm64\n'; return 1; }
    grep -qF '  - arm'    "$f" || { printf '  FAIL: missing arm\n'; return 1; }
    grep -qF '  - "386"'  "$f" || { printf '  FAIL: missing 386\n'; return 1; }
    grep -qF 'goarm:'     "$f" || { printf '  FAIL: missing goarm\n'; return 1; }
    grep -qF '  - 6'      "$f" || { printf '  FAIL: missing goarm 6\n'; return 1; }
    grep -qF '  - 7'      "$f" || { printf '  FAIL: missing goarm 7\n'; return 1; }
    grep -qF 'windows'    "$f" || { printf '  FAIL: missing windows/arm ignore\n'; return 1; }
    grep -qF '  goarch: arm' "$f" || { printf '  FAIL: missing arm in ignore\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-workflow-goreleaser-trigger()
#
# Verify goreleaser.yaml workflow only triggers on v*-nightly-* tags
# and uses the default core config (no --config flag)
#
test-workflow-goreleaser-trigger()
{
    local f="$REPO_DIR/.github/workflows/goreleaser.yaml"
    [[ -f "$f" ]] || { printf '  FAIL: %s not found\n' "$f"; return 1; }
    grep -qF "v*-nightly-*" "$f" || { printf '  FAIL: missing v*-nightly-* tag pattern\n'; return 1; }
    grep -qF "args: release --clean" "$f" || { printf '  FAIL: missing goreleaser release\n'; return 1; }
    grep -qF -- '--config=' "$f" && { printf '  FAIL: should not use --config (defaults to core)\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-workflow-release-config()
#
# Verify release.yml has -nightly- skip condition and --config flag
#
test-workflow-release-config()
{
    local f="$REPO_DIR/.github/workflows/release.yml"
    [[ -f "$f" ]] || { printf '  FAIL: %s not found\n' "$f"; return 1; }
    grep -qF "contains(github.ref_name, '-nightly-')" "$f" || {
        printf '  FAIL: missing nightly skip condition\n'; return 1; }
    grep -qF -- '--config=.goreleaser.full.yaml' "$f" || {
        printf '  FAIL: missing --config flag\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-workflow-nightly-format()
#
# Verify nightly-release.yml uses YYYY-MM-DD date format and curl-based dedup
#
test-workflow-nightly-format()
{
    local f="$REPO_DIR/.github/workflows/nightly-release.yml"
    [[ -f "$f" ]] || { printf '  FAIL: %s not found\n' "$f"; return 1; }
    grep -qF '%Y-%m-%d' "$f" || { printf '  FAIL: missing YYYY-MM-DD date format\n'; return 1; }
    grep -qF 'git/ref/tags' "$f" || { printf '  FAIL: missing curl-based tag existence check\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-workflow-nightly-tag-format()
#
# Verify nightly-release.yml builds tags matching v<ver>-nightly-<date>[-label]
#
test-workflow-nightly-tag-format()
{
    local f="$REPO_DIR/.github/workflows/nightly-release.yml"
    grep -qF 'TAG="v${{' "$f" || { printf '  FAIL: missing version prefix in tag\n'; return 1; }
    grep -qF 'nightly-${DATE}' "$f" || { printf '  FAIL: missing nightly-DATE in tag template\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# main()
#
main()
{
    if (( $# >= 1 )); then
        local cmd=$1; shift
        case "$cmd" in
            run-tests)                              run-tests;;
            test-goreleaser-yaml-core)               test-goreleaser-yaml-core;;
            test-goreleaser-full-yaml-exists)         test-goreleaser-full-yaml-exists;;
            test-goreleaser-full-yaml-content)        test-goreleaser-full-yaml-content;;
            test-workflow-goreleaser-trigger)         test-workflow-goreleaser-trigger;;
            test-workflow-release-config)             test-workflow-release-config;;
            test-workflow-nightly-format)             test-workflow-nightly-format;;
            test-workflow-nightly-tag-format)         test-workflow-nightly-tag-format;;
            *)                                        printf 'Unknown test: %s\n' "$cmd" >&2; exit 1;;
        esac
    else
        run-tests
    fi
}

if ! (return 0 2>/dev/null); then
    main "$@"
fi
