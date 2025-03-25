#! /usr/bin/env bash

main_bash () 
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

main_go ()
{
    [[ -d /opt/bin/ ]] || { printf 'Non-existent folder: /opt/bin/\n' >&2; return 1; }
    [[ -x /opt/bin/dylt ]] || { printf "/opt/bin/dylt not found or not executable\n" >&2; return 1; }
    
    local scriptKey=/daylight.sh
    local targetPath=/opt/bin/daylight.sh
    /opt/bin/dylt watch script "$scriptKey" "$targetPath"
}

main_go "$@"


