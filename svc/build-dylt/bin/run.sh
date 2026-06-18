#!/usr/bin/env bash

main ()
{
    local svcDir=/opt/svc/nightly-release
    cd "$svcDir/dylt-dev-dylt/repo" || exit 1
    printf 'Updating dylt repo ...\n'
    git pull --ff-only origin main || {
        printf 'git pull failed — no new commits or network issue\n' >&2
        exit 1
    }

    printf 'Running nightly release ...\n'
    source ./sunbeam.sh || { printf 'Failed to source sunbeam.sh\n' >&2; exit 1; }

    trigger-nightly-release "dylt-dev/dylt" || {
        printf 'Nightly release dispatch failed\n' >&2
        exit 1
    }

    printf 'Done — nightly release workflow dispatched\n'
}

main "$@"
