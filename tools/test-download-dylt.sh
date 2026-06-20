#! /usr/bin/env bash


SCRIPT_DIR=$(dirname "$(readlink -f "$BASH_SOURCE")")
SUNBEAM_SH="$SCRIPT_DIR/../sunbeam.sh"
[[ -f "$SUNBEAM_SH" ]] || { printf 'Cannot find sunbeam.sh\n' >&2; exit 1; }

TOKEN="${TOKEN:-}"
_positional=()
while (( $# > 0 )); do
    case $1 in
        --token) shift; TOKEN=$1 ;;
        -*)      printf 'Unknown flag: %s\n' "$1" >&2; exit 1 ;;
        *)       _positional+=("$1") ;;
    esac
    shift
done
set -- "${_positional[@]}"

command -v jq &>/dev/null || { printf 'jq is required but was not found\n' >&2; exit 1; }


#-------------------------------------------------------------------------------
#
# sb()
#
# @internal
# Shorthand for calling sunbeam.sh via case dispatch
#
sb()
{
    "$SUNBEAM_SH" "$@"
}


#-------------------------------------------------------------------------------
#
# run-tests()
#
# Run all test scenarios and report results
#
run-tests()
{
    local tests=(
        test-no-args
        test-nonexistent-dir
        test-token-no-value
        test-unknown-flag
        test-branch-mode
        test-release-mode
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
# test-no-args()
#
# Calling download-dylt-batch with no arguments should fail with usage
#
test-no-args()
{
    local output
    output=$(sb download-dylt-batch 2>&1) && {
        printf '  FAIL: expected error but succeeded\n'
        return 1
    }
    printf '%s' "$output" | grep -q 'Usage:' || {
        printf '  FAIL: unexpected output: %s\n' "$output"
        return 1
    }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-nonexistent-dir()
#
# Calling download-dylt-batch with a non-existent destination folder should fail
#
test-nonexistent-dir()
{
    local output
    output=$(sb download-dylt-batch /does/not/exist 2>&1) && {
        printf '  FAIL: expected error but succeeded\n'
        return 1
    }
    printf '%s' "$output" | grep -q 'Non-existent folder' || {
        printf '  FAIL: unexpected output: %s\n' "$output"
        return 1
    }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-token-no-value()
#
# --token without a value should return a clear error
#
test-token-no-value()
{
    local output
    output=$(sb download-dylt-batch --token 2>&1) && {
        printf '  FAIL: expected error but succeeded\n'
        return 1
    }
    printf '%s' "$output" | grep -q 'requires a value' || {
        printf '  FAIL: unexpected output: %s\n' "$output"
        return 1
    }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-unknown-flag()
#
# Unknown flag should return a clear error
#
test-unknown-flag()
{
    local output
    output=$(sb download-dylt-batch --bogus 2>&1) && {
        printf '  FAIL: expected error but succeeded\n'
        return 1
    }
    printf '%s' "$output" | grep -q 'Unknown flag' || {
        printf '  FAIL: unexpected output: %s\n' "$output"
        return 1
    }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-branch-mode()
#
# Download sunbeam.sh from the main branch via raw GitHub
#
test-branch-mode()
{
    local tmpDir; tmpDir=$(mktemp -d) || return 1
    sb download-dylt-batch --branch main "$tmpDir" || {
        printf '  FAIL: branch download failed\n'
        rm -rf "$tmpDir"; return 1
    }
    [[ -f "$tmpDir/sunbeam.sh" ]] || {
        printf '  FAIL: sunbeam.sh not downloaded\n'
        rm -rf "$tmpDir"; return 1
    }
    printf '  PASS\n'
    rm -rf "$tmpDir"
}


#-------------------------------------------------------------------------------
#
# test-release-mode()
#
# Download the latest dylt release artifact to a temp directory.
# Requires --token for reliable API access (no rate limiting).
# Without a token, the test is skipped.
#
test-release-mode()
{
    [[ -n "$TOKEN" ]] || { printf '  SKIP: no token provided\n'; return 0; }
    local tmpDir; tmpDir=$(mktemp -d) || return 1
    sb download-dylt-batch --release --latest --token "$TOKEN" "$tmpDir" || {
        printf '  FAIL: release download failed\n'
        rm -rf "$tmpDir"; return 1
    }
    local assetFile
    assetFile=$(ls "$tmpDir"/dylt_*.tar.gz 2>/dev/null | head -1)
    [[ -n "$assetFile" ]] && [[ -f "$assetFile" ]] || {
        printf '  FAIL: no dylt tar.gz found in %s\n' "$tmpDir"
        rm -rf "$tmpDir"; return 1
    }
    printf '  PASS (%s)\n' "$(basename "$assetFile")"
    rm -rf "$tmpDir"
}


#-------------------------------------------------------------------------------
#
# main()
#
# Dispatch to a specific test or run all
#
main ()
{
    if (( $# >= 1 )); then
        local cmd=$1
        shift
        case "$cmd" in
            run-tests)               run-tests;;
            test-no-args)            test-no-args;;
            test-nonexistent-dir)    test-nonexistent-dir;;
            test-token-no-value)     test-token-no-value;;
            test-unknown-flag)       test-unknown-flag;;
            test-branch-mode)        test-branch-mode;;
            test-release-mode)       test-release-mode;;
            *)                       printf 'Unknown test: %s\n' "$cmd" >&2; exit 1;;
        esac
    else
        run-tests
    fi
}


if ! (return 0 2>/dev/null); then
    main "$@"
fi
