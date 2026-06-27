#! /usr/bin/env bash

SCRIPT_DIR=$(dirname "$(readlink -f "$BASH_SOURCE")")
SUNBEAM_SH="$SCRIPT_DIR/../sunbeam.sh"
[[ -f "$SUNBEAM_SH" ]] || { printf 'Cannot find sunbeam.sh\n' >&2; exit 1; }

source "$SCRIPT_DIR/test-utils.sh" || exit 1
source "$SUNBEAM_SH" || { printf 'Failed to source sunbeam.sh\n' >&2; exit 1; }

#-------------------------------------------------------------------------------
#
# Mock gh() to capture calls without hitting GitHub
#
gh() {
    GH_CMD="$1"
    GH_ARGS=("$@")
    case "$1" in
        label)
            case "$2" in
                list)
                    printf '[{"name":"bug","color":"d73a4a","description":"Bug"},{"name":"enhancement","color":"a2eeef","description":"Feature request"},{"name":"documentation","color":"0075ca","description":"Docs"},{"name":"duplicate","color":"cfd3d7","description":"Duplicate issue or PR"},{"name":"invalid","color":"e4e669","description":"Not valid"}]'
                    ;;
                create|edit|delete)
                    return 0
                    ;;
                *)
                    printf 'unexpected gh label action: %s\n' "$2" >&2
                    return 1
                    ;;
            esac
            ;;
        *)
            printf 'unexpected gh command: %s\n' "$1" >&2
            return 1
            ;;
    esac
}

export -f gh
export GH_CMD GH_ARGS

#-------------------------------------------------------------------------------
#
# run-tests()
#
run-tests()
{
    local tests=(
        test-no-args
        test-unknown-flag
        test-all-with-positional
        test-dry-run-rename
        test-dry-run-create
        test-dry-run-delete
        test-snapshot-written
        test-skip-compliant-labels
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
test-no-args()
{
    fail-check "no args" 'Usage' reconcile-labels
}


#-------------------------------------------------------------------------------
#
# test-unknown-flag()
#
test-unknown-flag()
{
    fail-check "unknown flag" 'Unknown flag' reconcile-labels --bogus owner/repo
}


#-------------------------------------------------------------------------------
#
# test-all-with-positional()
#
test-all-with-positional()
{
    fail-check "--all with positional" 'error: --all' reconcile-labels --all owner/repo
}


#-------------------------------------------------------------------------------
#
# test-dry-run-rename()
#
# Verify dry-run mode prints rename without calling gh label edit
#
test-dry-run-rename()
{
    local output
    output=$(reconcile-labels --dry-run dylt-dev/test-repo 2>&1) || {
        printf '  FAIL (dry-run rename): reconcile-labels failed\n'
        return 1
    }
    if [[ $output != *'would rename: enhancement -> feature'* ]]; then
        printf '  FAIL (dry-run rename): expected "would rename" message, got:\n  %s\n' "$output"
        return 1
    fi
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-dry-run-create()
#
# Verify dry-run prints create messages for missing task label
#
test-dry-run-create()
{
    # Mock gh label list to return only bug and enhancement (no task)
    gh() {
        case "$1/$2" in
            label/list)
                printf '[{"name":"bug","color":"d73a4a","description":"Bug"},{"name":"enhancement","color":"a2eeef","description":"Feature request"}]'
                ;;
            *)
                return 0
                ;;
        esac
    }
    local output
    output=$(reconcile-labels --dry-run dylt-dev/test-repo 2>&1) || {
        printf '  FAIL (dry-run create): reconcile-labels failed\n'
        return 1
    }
    if [[ $output != *'would create: task'* ]]; then
        printf '  FAIL (dry-run create): expected "would create: task" message, got:\n  %s\n' "$output"
        return 1
    fi
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-dry-run-delete()
#
# Verify dry-run prints delete messages for extra labels
#
test-dry-run-delete()
{
    local output
    output=$(reconcile-labels --dry-run dylt-dev/test-repo 2>&1) || {
        printf '  FAIL (dry-run delete): reconcile-labels failed\n'
        return 1
    }
    if [[ $output != *'would delete: documentation'* ]]; then
        printf '  FAIL (dry-run delete): expected "would delete: documentation", got:\n  %s\n' "$output"
        return 1
    fi
    if [[ $output != *'would delete: duplicate'* ]]; then
        printf '  FAIL (dry-run delete): expected "would delete: duplicate", got:\n  %s\n' "$output"
        return 1
    fi
    if [[ $output != *'would delete: invalid'* ]]; then
        printf '  FAIL (dry-run delete): expected "would delete: invalid", got:\n  %s\n' "$output"
        return 1
    fi
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-snapshot-written()
#
# Verify snapshot file is written when a local clone is found
#
test-snapshot-written()
{
    local tmpdir
    tmpdir=$(mktemp -d) || { printf '  FAIL (snapshot): mktemp failed\n'; return 1; }
    mkdir -p "$tmpdir/.git" "$tmpdir/misc"

    # Override HOME to point at tmpdir and create project dir
    local orig_home=$HOME
    export HOME=$(mktemp -d)
    mkdir -p "$HOME/src/test-repo/.git"

    # Set up remote
    git -C "$HOME/src/test-repo" init --quiet 2>/dev/null
    git -C "$HOME/src/test-repo" remote add origin https://github.com/dylt-dev/test-repo.git 2>/dev/null

    # Mock label list
    gh() {
        printf '[{"name":"bug","color":"d73a4a","description":"Bug"}]'
    }

    reconcile-labels dylt-dev/test-repo 2>/dev/null || true

    if [[ ! -f "$HOME/src/test-repo/misc/gh-labels.json" ]]; then
        printf '  FAIL (snapshot): misc/gh-labels.json not created\n'
        export HOME=$orig_home
        rm -rf "$HOME"
        return 1
    fi
    export HOME=$orig_home
    rm -rf "$HOME"
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-skip-compliant-labels()
#
# Verify nothing is renamed/created/deleted when only bug/feature/task exist
#
test-skip-compliant-labels()
{
    gh() {
        if [[ "$1/$2" == label/list ]]; then
            printf '[{"name":"bug","color":"d73a4a","description":"Bug"},{"name":"feature","color":"a2eeef","description":"Feature"},{"name":"task","color":"c5def5","description":"Task"}]'
        else
            printf 'FAIL: unexpected gh call: %s\n' "$*" >&2
            return 1
        fi
    }
    local output
    output=$(reconcile-labels --dry-run dylt-dev/test-repo 2>&1) || {
        printf '  FAIL (skip): reconcile-labels failed on compliant repo\n'
        return 1
    }
    if [[ $output == *'would rename'* || $output == *'would create'* || $output == *'would delete'* ]]; then
        printf '  FAIL (skip): expected no changes for compliant repo, got:\n  %s\n' "$output"
        return 1
    fi
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# main()
#
case ${1:-} in
    *)  run-tests;;
esac


if ! (return 0 2>/dev/null); then
    main "$@"
fi
