#! /usr/bin/env bash

SCRIPT_DIR=$(dirname "$(readlink -f "$BASH_SOURCE")")
SUNBEAM_SH="$SCRIPT_DIR/../sunbeam.sh"
[[ -f "$SUNBEAM_SH" ]] || { printf 'Cannot find sunbeam.sh\n' >&2; exit 1; }

# Source sunbeam.sh but save the original batch function so we can restore it
source "$SUNBEAM_SH" || { printf 'Failed to source sunbeam.sh\n' >&2; exit 1; }

#-------------------------------------------------------------------------------
#
# run-tests()
#
run-tests()
{
    local tests=(
        test-flag-parsed
        test-default-path-uses-home
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
# make-mock-daylight-sh()
#
# @internal
# Create a minimal daylight.sh that responds to gen-completion-script
# and list-bash-funcs.
#
make-mock-daylight-sh()
{
    local dst=$1/daylight.sh
    cat > "$dst" << 'MOCKEOF'
#! /usr/bin/env bash
gen-completion-script-batch() { while read -r l; do printf '%s\n' "$l"; done; }
list-bash-funcs() { echo cmd1; echo cmd2; }
case ${1:-} in
    gen-completion-script-batch) gen-completion-script-batch;;
    list-bash-funcs) list-bash-funcs;;
esac
MOCKEOF
    chmod +x "$dst"
    printf '%s' "$dst"
}


#-------------------------------------------------------------------------------
#
# setup-mock()
#
# @internal
# Override download-daylight-batch with a mock, set up temp dir.
# Returns via echo the temp dir path.
#
setup-mock()
{
    local tmpDir; tmpDir=$(mktemp -d) || return 1
    make-mock-daylight-sh "$tmpDir" >/dev/null

    eval "original_batch_$(date +%s)()" { download-daylight-batch '"$@"' ';' } 2>/dev/null || true

    download-daylight-batch()
    {
        local args=("$@")
        local dst=""
        local i=0
        while (( i < ${#args[@]} )); do
            if [[ "${args[i]}" != --* ]] && [[ "${args[i]}" != -* ]]; then
                dst="${args[i]}"
            fi
            (( i++ ))
        done
        [[ -n "$dst" ]] && mkdir -p "$dst"
        make-mock-daylight-sh "$dst" >/dev/null
    }

    printf '%s' "$tmpDir"
}


#-------------------------------------------------------------------------------
#
# test-flag-parsed()
#
test-flag-parsed()
{
    local tmpDir; tmpDir=$(setup-mock) || return
    local compFile="/tmp/test-142-comp.bash"

    download-daylight --gen-bash-completions "$compFile" --branch main "$tmpDir"

    if [[ ! -f "$compFile" ]]; then
        printf '  FAIL: completion file was not created at %s\n' "$compFile"
        rm -f "$compFile"
        rm -rf "$tmpDir"
        return 1
    fi

    if [[ ! -s "$compFile" ]]; then
        printf '  FAIL: completion file is empty\n'
        rm -f "$compFile"
        rm -rf "$tmpDir"
        return 1
    fi

    rm -f "$compFile"
    rm -rf "$tmpDir"
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-default-path-uses-home()
#
test-default-path-uses-home()
{
    local tmpDir; tmpDir=$(mktemp -d) || return
    local home_bak=$HOME
    export HOME="$tmpDir/home"
    mkdir -p "$HOME"

    make-mock-daylight-sh "$tmpDir" >/dev/null

    download-daylight-batch()
    {
        local args=("$@")
        local dst=""
        local i=0
        while (( i < ${#args[@]} )); do
            if [[ "${args[i]}" != --* ]] && [[ "${args[i]}" != -* ]]; then
                dst="${args[i]}"
            fi
            (( i++ ))
        done
        [[ -n "$dst" ]] && mkdir -p "$dst"
        make-mock-daylight-sh "$dst" >/dev/null
    }

    download-daylight --gen-bash-completions --branch main "$tmpDir"

    local expected="$HOME/bash-completion.d/daylight.sh"
    if [[ ! -f "$expected" ]]; then
        printf '  FAIL: default path %s not found\n' "$expected"
        export HOME="$home_bak"
        rm -rf "$tmpDir"
        return 1
    fi

    export HOME="$home_bak"
    rm -rf "$tmpDir"
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

