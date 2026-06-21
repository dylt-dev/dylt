#! /usr/bin/env bash
# Manual functional test — download daylight.sh via sunbeam.sh, generate
# bash completions, and verify they auto-wire into your shell:
#
#   dst=/tmp/dl-test && mkdir -p "$dst"
#   source ./sunbeam.sh
#   download-daylight --gen-bash-completions ~/.bash-completion.d/daylight.sh "$dst"
#   grep 'daylight.sh' ~/.bashrc
#   source ~/.bashrc
#   daylight.sh <tab>
#
# Or interactively:
#
#   download-daylight "$dst"
#   # answer y at the prompt
#   grep 'daylight.sh' ~/.bashrc
#
# This test script automates the underlying function behaviour; the manual
# test above exercises the full end-to-end pipeline.

SCRIPT_DIR=$(dirname "$(readlink -f "$BASH_SOURCE")")
source "$SCRIPT_DIR/test-utils.sh" || exit 1
SUNBEAM_SH="$SCRIPT_DIR/../sunbeam.sh"
[[ -f "$SUNBEAM_SH" ]] || { printf 'Cannot find sunbeam.sh\n' >&2; exit 1; }
source "$SUNBEAM_SH" || { printf 'Failed to source sunbeam.sh\n' >&2; exit 1; }

INSTALL_SH="$SCRIPT_DIR/../install-sunbeam.sh"


#-------------------------------------------------------------------------------
#
# run-tests()
#
run-tests()
{
    local tests=(
        test-install-sunbeam-calls-add-to-bashrc
        test-download-daylight-wires-bashrc
        test-download-daylight-no-duplicate
        test-download-daylight-interactive-wires-bashrc
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
# make-stub-daylight()
#
# @internal
# Create a minimal daylight.sh stub in the given directory so that the
# inner bash subprocess invoked by download-daylight can run
# gen-completion-script and list-bash-funcs.
#
make-stub-daylight()
{
    local dst=$1
    mkdir -p "$dst"
    cat > "$dst/daylight.sh" << 'EOF'
gen-completion-script()
{
    while read -r l; do printf '%s\n' "$l"; done
}
list-bash-funcs()
{
    printf 'cmd1\ncmd2\n'
}
if [[ $# -gt 0 ]]; then
    case "$1" in
        gen-completion-script) gen-completion-script;;
        list-bash-funcs) list-bash-funcs;;
    esac
fi
EOF
    chmod +x "$dst/daylight.sh"
}


#-------------------------------------------------------------------------------
#
# test-install-sunbeam-calls-add-to-bashrc()
#
# Verify install-sunbeam.sh invokes add-to-bashrc after install
#
test-install-sunbeam-calls-add-to-bashrc()
{
    [[ -f "$INSTALL_SH" ]] || { printf '  FAIL: install-sunbeam.sh not found\n'; return 1; }
    grep -qF 'add-to-bashrc' "$INSTALL_SH" || {
        printf '  FAIL: install-sunbeam.sh does not call add-to-bashrc\n'; return 1; }
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-download-daylight-wires-bashrc()
#
# Verify download-daylight appends source line to .bashrc after generating
# completions (--gen-bash-completions path)
#
test-download-daylight-wires-bashrc()
{
    local tmpdir; tmpdir=$(mktemp -d) || return 1
    local bashrc="$tmpdir/.bashrc"
    touch "$bashrc"

    # Mock download-daylight-batch to create a stub daylight.sh
    eval 'download-daylight-batch() { make-stub-daylight "${@: -1}"; return 0; }' || true

    HOME="$tmpdir" download-daylight --gen-bash-completions "$tmpdir/comps" "$tmpdir/dst" 2>/dev/null || {
        printf '  FAIL: download-daylight failed\n'
        rm -rf "$tmpdir"
        return 1
    }

    # Verify the source line was added
    grep -qF "source $tmpdir/comps" "$bashrc" || {
        printf '  FAIL: source line not found in .bashrc\n'
        cat "$bashrc" >&2
        rm -rf "$tmpdir"
        return 1
    }

    rm -rf "$tmpdir"
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-download-daylight-no-duplicate()
#
# Verify running download-daylight twice does not duplicate the source line
#
test-download-daylight-no-duplicate()
{
    local tmpdir; tmpdir=$(mktemp -d) || return 1
    local bashrc="$tmpdir/.bashrc"
    touch "$bashrc"

    eval 'download-daylight-batch() { make-stub-daylight "${@: -1}"; return 0; }' || true

    HOME="$tmpdir" download-daylight --gen-bash-completions "$tmpdir/comps" "$tmpdir/dst" 2>/dev/null || {
        printf '  FAIL: first call failed\n'; rm -rf "$tmpdir"; return 1; }

    # Second call with a different dst dir to avoid conflicts
    HOME="$tmpdir" download-daylight --gen-bash-completions "$tmpdir/comps" "$tmpdir/dst2" 2>/dev/null || {
        printf '  FAIL: second call failed\n'; rm -rf "$tmpdir"; return 1; }

    local count
    count=$(grep -c "source $tmpdir/comps" "$bashrc" 2>/dev/null || true)
    (( count == 1 )) || {
        printf '  FAIL: expected 1 source line, got %d\n' "$count"
        cat "$bashrc" >&2
        rm -rf "$tmpdir"
        return 1
    }

    rm -rf "$tmpdir"
    printf '  PASS\n'
}


#-------------------------------------------------------------------------------
#
# test-download-daylight-interactive-wires-bashrc()
#
# Verify the interactive path (no --gen-bash-completions, user answers y)
# also wires .bashrc
#
test-download-daylight-interactive-wires-bashrc()
{
    local tmpdir; tmpdir=$(mktemp -d) || return 1
    local bashrc="$tmpdir/.bashrc"
    touch "$bashrc"

    eval 'download-daylight-batch() { make-stub-daylight "${@: -1}"; return 0; }' || true

    printf 'y\n' | HOME="$tmpdir" download-daylight "$tmpdir/dst" 2>/dev/null || {
        printf '  FAIL: download-daylight failed\n'
        rm -rf "$tmpdir"
        return 1
    }

    grep -qF "source $tmpdir/bash-completion.d/daylight.sh" "$bashrc" || {
        printf '  FAIL: source line not found in .bashrc\n'
        cat "$bashrc" >&2
        rm -rf "$tmpdir"
        return 1
    }

    rm -rf "$tmpdir"
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
            run-tests)                                   run-tests;;
            test-install-sunbeam-calls-add-to-bashrc)     test-install-sunbeam-calls-add-to-bashrc;;
            test-download-daylight-wires-bashrc)          test-download-daylight-wires-bashrc;;
            test-download-daylight-no-duplicate)          test-download-daylight-no-duplicate;;
            test-download-daylight-interactive-wires-bashrc) test-download-daylight-interactive-wires-bashrc;;
            *)                                            printf 'Unknown test: %s\n' "$cmd" >&2; exit 1;;
        esac
    else
        run-tests
    fi
}

if ! (return 0 2>/dev/null); then
    main "$@"
fi
