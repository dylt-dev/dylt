#! /usr/bin/env bash

#! /usr/bin/env bash


main ()
{
    local COMMANDS_FILE_PATH='./cli-run.sh'
    # shellcheck disable=SC2016
    [[ -f "$COMMANDS_FILE_PATH" ]] || { printf 'Non-existent path: %s\n' "$COMMANDS_FILE_PATH" >&2; return 1; }

    while read -r line; do
        # printf '=%.0s' {1..50}
        # echo
        local output; output=$($line)
        printf '$ %s\n' "$line" || return
        # printf '=%.0s' {1..50}
        # echo
        echo
        printf '%s\n' "$output" || return
        echo
        # printf '=%.0s' {1..50}
        echo
        echo
        echo
    done <"$COMMANDS_FILE_PATH"    
}


(return 0 2>/dev/null) || main "$@"