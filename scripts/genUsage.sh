#! /usr/bin/env bash

handle-file () {
    # shellcheck disable=SC2016
    (( $# == 1 )) || { printf 'Usage: handle-file $path\n' >&2; return 1; }
    local path=$1

    printf '$path=%s\n' "$path"
    # shellcheck disable=SC2016
    [[ -f "$path" ]] || { printf 'Non-existent path: $path\n' >&2; return 1; }
    grep --after-context=20 'func .*create.*Command' "$path" # | grep case
    echo
}

list-create-command-files () {
    # shellcheck disable=SC2016
    (( $# >= 0 && $# <= 1 )) || { printf 'Usage: list-create-command-files [$folder]\n' >&2; return 1; }
    local folder=${1:-.}
    grep --recursive --include "*.go" --files-with-matches 'create.*Command' "$folder"
}

main () {
    while IFS= read -r line; do
        printf '%s\n' "$line" || return
        handle-file "$line"
    done < <(list-create-command-files)
}