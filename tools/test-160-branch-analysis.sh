#! /usr/bin/env bash

SCRIPT_DIR=$(dirname "$(readlink -f "$BASH_SOURCE")")
SUNBEAM_SH="$SCRIPT_DIR/../sunbeam.sh"
[[ -f "$SUNBEAM_SH" ]] || { printf 'Cannot find sunbeam.sh\n' >&2; exit 1; }

source "$SCRIPT_DIR/test-utils.sh" || exit 1
source "$SUNBEAM_SH" || { printf 'Failed to source sunbeam.sh\n' >&2; exit 1; }


test-gh-branch-list-usage()
{
    gh-branch-list a 2>/dev/null && {
        printf '  FAIL (gh-branch-list-usage): expected failure with args\n'
        return 1
    }
    printf '  PASS\n'
}


test-gh-branch-compare-usage()
{
    gh-branch-compare 2>/dev/null && {
        printf '  FAIL (gh-branch-compare-usage): expected failure with no args\n'
        return 1
    }
    gh-branch-compare a b c 2>/dev/null && {
        printf '  FAIL (gh-branch-compare-usage): expected failure with 3 args\n'
        return 1
    }
    printf '  PASS\n'
}


test-gh-branch-merged-usage()
{
    gh-branch-merged a b 2>/dev/null && {
        printf '  FAIL (gh-branch-merged-usage): expected failure with 2 args\n'
        return 1
    }
    printf '  PASS\n'
}


test-gh-branch-clean-usage()
{
    gh-branch-clean a b 2>/dev/null && {
        printf '  FAIL (gh-branch-clean-usage): expected failure with 2 args\n'
        return 1
    }
    printf '  PASS\n'
}


test-gh-branch-list-runs()
{
    command -v gh >/dev/null 2>&1 || {
        printf '  SKIP (gh-branch-list-runs): gh not installed\n'
        return 0
    }

    local out
    out=$(gh-branch-list 2>/dev/null) || {
        printf '  FAIL (gh-branch-list-runs): returned non-zero\n'
        return 1
    }
    printf '%s' "$out" | grep -q 'main' || {
        printf '  FAIL (gh-branch-list-runs): "main" not in branch list\n'
        return 1
    }
    printf '  PASS\n'
}


run-tests()
{
    local tests=(
        test-gh-branch-list-usage
        test-gh-branch-compare-usage
        test-gh-branch-merged-usage
        test-gh-branch-clean-usage
        test-gh-branch-list-runs
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


main()
{
    case ${1:-all} in
        test-gh-branch-list-usage)     test-gh-branch-list-usage ;;
        test-gh-branch-compare-usage)  test-gh-branch-compare-usage ;;
        test-gh-branch-merged-usage)   test-gh-branch-merged-usage ;;
        test-gh-branch-clean-usage)    test-gh-branch-clean-usage ;;
        test-gh-branch-list-runs)      test-gh-branch-list-runs ;;
        all)                           run-tests ;;
        *)                             printf 'Unknown test: %s\n' "$1" >&2; exit 1 ;;
    esac
}


if ! (return 0 2>/dev/null); then
    main "$@"
fi
