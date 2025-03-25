#! /usr/bin/env bash

main ()
{
    [[ -d /opt/bin/ ]] || { printf 'Non-existent folder: /opt/bin/\n' >&2; return 1; }
    [[ -x /opt/bin/dylt ]] || { printf "/opt/bin/dylt not found or not executable\n" >&2; return 1; }
    
    /opt/bin/dylt watch svc
}

main "$@"


