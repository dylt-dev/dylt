#! /usr/bin/env bash

SCRIPT_DIR=$(dirname "$(readlink -f "$BASH_SOURCE")")
source "$SCRIPT_DIR/test-utils.sh" || exit 1
SUNBEAM_SH="$SCRIPT_DIR/../sunbeam.sh"
[[ -f "$SUNBEAM_SH" ]] || { printf 'Cannot find sunbeam.sh\n' >&2; exit 1; }
source "$SUNBEAM_SH" || { printf 'Failed to source sunbeam.sh\n' >&2; exit 1; }


#-------------------------------------------------------------------------------
#
# run-tests()
#
run-tests()
{
    local tests=(
        test-download-dylt-batch-unknown-flag
        test-download-dylt-batch-extract-flags
        test-download-dylt-batch-output-dir-remote-name
        test-download-dylt-batch-extract-with-platform
        test-download-daylight-batch-extract-flags
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
# test-download-dylt-batch-unknown-flag()
#
# Verify download-dylt-batch rejects unknown flags
#
test-download-dylt-batch-unknown-flag()
{
    fail-check "unknown flag" "Unknown flag" download-dylt-batch --bogus /nonexistent
}


#-------------------------------------------------------------------------------
#
# test-download-dylt-batch-extract-flags()
#
# Verify --extract/--extract-dir/--extract-name parse correctly
#
test-download-dylt-batch-extract-flags()
{
    fail-check "extract flags" "Non-existent folder" \
        download-dylt-batch --extract --extract-dir /opt --extract-name dylt /nonexistent
}


#-------------------------------------------------------------------------------
#
# test-download-dylt-batch-output-dir-remote-name()
#
# Verify --output-dir and --remote-name pass through
#
test-download-dylt-batch-output-dir-remote-name()
{
    fail-check "output-dir remote-name" "Non-existent folder" \
        download-dylt-batch --output-dir /tmp --remote-name --token sekret /nonexistent
}


#-------------------------------------------------------------------------------
#
# test-download-dylt-batch-extract-with-platform()
#
# Verify extract flags work alongside platform positional
#
test-download-dylt-batch-extract-with-platform()
{
    fail-check "extract with platform" "Non-existent folder" \
        download-dylt-batch --extract --extract-dir /opt --release --token x /nonexistent linux-amd64
}


#-------------------------------------------------------------------------------
#
# test-download-daylight-batch-extract-flags()
#
# Verify download-daylight-batch (sunbeam version) parses extract flags
#
test-download-daylight-batch-extract-flags()
{
    fail-check "daylight extract flags" "Non-existent folder" \
        download-daylight-batch --extract --extract-dir /opt --extract-name x /nonexistent
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
            run-tests)                                  run-tests;;
            test-download-dylt-batch-unknown-flag)      test-download-dylt-batch-unknown-flag;;
            test-download-dylt-batch-extract-flags)     test-download-dylt-batch-extract-flags;;
            test-download-dylt-batch-output-dir-remote-name) test-download-dylt-batch-output-dir-remote-name;;
            test-download-dylt-batch-extract-with-platform) test-download-dylt-batch-extract-with-platform;;
            test-download-daylight-batch-extract-flags) test-download-daylight-batch-extract-flags;;
            *)                                          printf 'Unknown test: %s\n' "$cmd" >&2; exit 1;;
        esac
    else
        run-tests
    fi
}

if ! (return 0 2>/dev/null); then
    main "$@"
fi
