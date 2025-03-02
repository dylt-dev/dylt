#! /usr/bin/env bash

main () 
{ 
    printf "Downloading current script ...\n";
    # etcdctl get, redirect to tempfile, and cp to proper location on success
    local tmpfile; tmpfile="$(mktemp --tmpdir daylight.sh.XXXXXX)" || return
    if ! /opt/etcd/etcdctl \
        --discovery-srv hello.dylt.dev \
        get --print-value-only /daylight.sh >"$tmpfile";
    then
        local rc=$?
        printf '%s\n' "Failed to download /daylight.sh from cluster" >&2
        return $?
    else
        printf 'Download succeeded. Copying script to final location ...\n'
        cp "$tmpfile" /opt/bin/daylight.sh
    fi
    printf 'Watching for further updates ....\n'
    /opt/etcd/etcdctl \
        --discovery-srv hello.dylt.dev \
        watch /daylight.sh \
            -- sh -c '{ printf "Downloading update ..."; tmpfile="$(mktemp --tmpdir daylight.sh.XXXXXX)"; /opt/etcd/etcdctl --discovery-srv hello.dylt.dev get --print-value-only /daylight.sh >"$tmpfile"; cp "$tmpfile" /opt/bin/daylight.sh; printf "Complete.\n"; }' || return
}

main "$@"


