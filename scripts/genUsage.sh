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

list-cases () {
    # shellcheck disable=SC2016
    (( $# == 1 )) || { printf 'Usage: list-cases $path\n' >&2; return 1; }
    local path=$1
    # shellcheck disable=SC2016
    [[ -f "$path" ]] || { printf 'Non-existent path: $path\n' >&2; return 1; }

    sed -n '/func.*create.*Command/,/^}/p' <"$path" | grep case
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

# How to do a file
#   Run list-create-command-files to get all files
#   Run list-cases to get case lines that belong in usage
#   Create func (*cmd) PrintUsage () {}
#   Paste case lines into the func
#   Run subst to fix up the lines
#   Edit Run() method to check for 0 args and call PrintUsage