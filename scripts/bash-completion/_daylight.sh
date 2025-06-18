#! /usr/bin/env bash


### Helper functions

colorComment=$'\033[96m'
reset=$'\033[0m' 

# Style text + print ANSI reset when done
comment () {
	printf '%s%s%s\n' "$colorComment" "$1" "$reset" >> /tmp/dylt.log
}

# Set COMPREPLY to (), and set DONE
complete-with-empty () {
	COMPREPLY=()
	DONE=1
}

# Set COMPREPLY to files matched by the current token ($1), and set DONE
complete-with-files () {
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: complete-with-files $token\n' >&2; return 1; }
	local token=$1

	compopt -o filenames
	COMPREPLY=($(compgen -f -- "$token"))
	DONE=1
}

# Set COMPREPLY for the elements of a wordlist (#1) matched by the current 
# token ($2), and set DONE
complete-with-words () {
	# shellcheck disable=SC2016
	(( $# == 2 )) || { printf 'Usage: complete-with-words $wordlist $token\n' >&2; return 1; }
	local wordlist=$1
	local token=$2

	COMPREPLY=($(compgen -W "$wordlist" -- "$token"))
	DONE=1
}

# Get the next token from the command line, and increment N. The current
# token is just $COMP_WORDS indexed by N.
get-token () {
	local -n _ref=$1	
	_ref=${COMP_WORDS[N]}
	((N++))	
	printf '%-16s: token=%s N=%d\n' "get-token()" "$_ref" "$N" >> /tmp/dylt.log
}

# Various bash completion values, global variables, etc. Useful for development/debugging.
status () {
	printf '%-16s: DONE=%d N=%d COMP_CWORD=%d COMP_TYPE=%d cur=[%s] <%s>\n' "$1()" "$DONE" "$N" "$COMP_CWORD" "$COMP_TYPE" "${COMP_WORDS[COMP_CWORD]}" "${COMP_WORDS[*]}" >> /tmp/dylt.log
}

### Main script

# Entry point. Set global variables and call the main function
_dylt () {
	echo Initializing N >> /tmp/dylt.log
	N=1
	DONE=0
	status '_dylt'
	on-main
	printf '\n\n' >> /tmp/dylt.log
}


# dylt
on-main () {
    # anything goes
    argvals=(call config get host init list misc vm watch)
    flags=()
    get-token token
    status on-main 
    if on-last-token; then
        comment "$(printf "latest token; complete with flags or subcommands")"
        case $token in
            -*)	complete-with-words	"${flags[*]}"	"$token";;
            *)	complete-with-words	"${argvals[*]}"	"$token";;
        esac
        return
    fi
	
	# Handle subcommand
	comment "$(printf "Ready for what's next: token=%s\n" "$token")"
	case $token in
		call)	on-call;	status X-on-call ;;
		config)	on-config;	status X-on-config ;;
		get)	on-get;		status X-on-get ;;
		host)	on-host;	status X-on-host ;;
		init)	on-init;	status X-on-init ;;
		list)   on-list;	status X-on-list ;;
		misc)	on-misc;	status X-on-misc ;;
		vm)		on-vm;		status X-on-vm ;;
		watch)	on-watch;	status X-on-watch ;;
		*) complete-with-empty;;
	esac
			
	# Done
	comment "$(printf 'on-main() - done; last token=%s' "$token")"
}

# dylt call [--flag flagval]
#
# @note valid functions + args for daylight.sh could be used for autocompletion
on-call () {
	local flags
	
	#  - complete with flags only
	flags=(--script-path)
	get-token token
	status on-call
	if on-last-token; then
		comment "$(printf "we have arrived at the latest token; time to generate completions")"
		case $token in
			-*) complete-with-words "${flags[*]}" "$token";;
			*)  complete-with-empty;;
		esac
		return
	fi
	
	# next token should be a flag. go to handler function, or complete with empty.
	case "$token" in
		--script-path)	on-call-scriptPath;	status X-on-c-sp;;
		*)	complete-with-empty;;
	esac
}


# dylt call --script-path /path/to/dylight.sh
#
# complete with files
on-call-scriptPath () {
	get-token token
	status on-call-scriptPath
	if on-last-token; then
		comment "$(printf "on latest token; complete with files")"
		complete-with-files "$token"
	fi

	# Done
	comment "$(printf 'on-call-scriptPath() - done; last token=%s' "$token")"
}

# dylt config subcommand
#
# complete with subcommand
on-config () {
    argvals=(get set show)
    get-token token
    status on-config
    if on-last-token; then
        comment 'generate completions for subcommand'
        case "$token" in
            *) complete-with-words "${argvals[*]}" "$token";;
        esac
        return
    fi

    # Handle subcommand
    comment "$(printf 'Calling subcommand handler (token=%s)\n' "$token")"
    case "$token" in
        get)  on-config-get;  status X-on-config-get;;
		set)  on-config-set;  status X-on-config-set;;
		show) on-config-show; status X-on-config-show;;
        *) complete-with-empty;;
    esac

    # Done
    comment "$(printf 'on-config() - done; last token=%s' "$token")"
}


# dylt config get key
on-config-get () {
	local argvals

	# key - complete with argvals
	argvals=(name age luckyNumber)
	get-token token
	status on-config-get-1
	if on-last-token; then
		comment 'generate completions - argument values only'
		case $token in
			*) complete-with-words "${argvals[*]}" "$token";;
		esac
		return
	fi
	
	# Done
	comment "$(printf 'config-get() - done; last token=%s' "$token")"
}


# dylt config set key val
on-config-set () {
	local argvals

	# key - complete with argvals
	argvals=(name age luckyNumber)
	get-token token
	status on-config-set-1
	if on-last-token; then
		comment 'generate completions - argument values only'
		case "$token" in
			*) complete-with-words "${argvals[*]}" "$token";;
		esac
		return
	fi
	
	# val (complete with empty)
	get-token token
	status on-config-set-2
	if on-last-token; then
		comment 'handle latest token; complete with empty'
		complete-with-empty
		return
	fi

	# Done
	comment "$(printf 'on-config-set() - done; last token=%s' "$token")"
}


# on config show
#
# No args or flags - complete with empty
on-config-show () {
	status on-config-show
	complete-with-empty
}


# on get key
#
# freeform arg - complete with empty
# @note contraining `key` on valid cluster keys might be good
on-get () {
	status on-get
	complete-with-empty
}


# on host subcmd
on-host () {
	local argvals flags

	# complete subcommand
	argvals=(init)
	get-token token
	status 
	if on-last-token; then
		comment 'generate completions for subcommand'
		complete-with-words "${argvals[*]}" "$token"
		return
	fi
	
	# next token should be subcommand. go to handler function, or complete with empty.
	case "$token" in
		init) on-host-init; status X-on-host-init;;
		*)    complete-with-empty;;
	esac
}


# dylt host init uid gid
on-host-init () {
	# first arg: uid (complete with empty)
	get-token token
	status on-host-init-1
	if on-last-token; then
		comment "$(printf "we have arrived at the latest token; time to generate completions")"
		complete-with-empty
		return
	fi

	# second arg: gid (complete with empty)
	get-token token
	status on-host-init-2
	if on-last-token; then
		comment "$(printf "we have arrived at the latest token; time to generate completions")"
		complete-with-empty
		return
	fi

	 # Done
	comment "$(printf 'on-host-init() - done; last token=%s' "$token")"
}


# dylt init --flag
on-init () {
    #  - complete with flags only
    flags=(--etcd-domain)
    get-token token
    status on-init-1
    if on-last-token; then
        comment "$(printf "we have arrived at the latest token; time to generate completions")"
        case $token in
            -*) complete-with-words "${flags[*]}" "$token";;
            *)  complete-with-empty;;
        esac
        return
    fi
    
    # next token should be a flag. go to handler function, or complete with empty.
    case "$token" in
        --etcd-domain)	on-init-etcdDomain;	status X-on-init-etcdDomain;;
        *)	complete-with-empty;;
    esac

     # Done
    comment "$(printf 'on-init() - done; last token=%s' "$token")"
}


# dylt init --etcd-domain etcd.cluster.domain
#
# Freeform - complete with empty
on-init-etcdDomain () {
    get-token token
    status on-init-etcdDomain
    if on-last-token; then
        comment "$(printf "on latest token; complete with empty")"
        complete-with-empty
    fi

    # Done
    comment "$(printf 'on-init-etcdDomain() - done; last token=%s' "$token")"
}


# dylt list
#
# No args, no flags - complete with empty
on-list () {
    status on-list
    complete-with-empty

    # Done
    comment "on-list() - done"
}


# dylt misc subcommand
#
# complete with subcommand
on-misc () {
    argvals=(create-two-node-cluster gen-etcd-run-script)
    get-token token
    status on-misc
    if on-last-token; then
        comment 'generate completions for subcommand'
        case "$token" in
            *) complete-with-words "${argvals[*]}" "$token";;
        esac
        return
    fi

    # Handle subcommand
    comment "$(printf 'Calling subcommand handler (token=%s)\n' "$token")"
    case "$token" in
        create-two-node-cluster)    on-misc-createTwoNodeCluster;   status X-on-misc-createTwoNodeCluster;;
        gen-etcd-run-script)        on-misc-genEtcdRunScript;       status X-on-misc-genEtcdRunScript;;
        *) complete-with-empty;;
    esac

    # Done
    comment "$(printf 'on-misc() - done; last token=%s' "$token")"
}


# dylt misc create-two-node-cluster
#
# No args, no flags - complete with empty
on-misc-createTwoNodeCluster () {
    status on-createTwoNodeCluster
    complete-with-empty

    # Done
    comment "on-createTwoNodeCluster() - done"
}


# dylt misc gen-etcd-run-script
#
# No args, no flags - complete with empty
on-misc-createTwoNodeCluster () {
    status on-createTwoNodeCluster
    complete-with-empty

    # Done
    comment "on-createTwoNodeCluster() - done"
}


# dylt vm subcommand
#
# complete with subcommand
on-vm () {
    argvals=(add all del get list set)
    get-token token
    status on-vm
    if on-last-token; then
        comment 'generate completions for subcommand'
        case "$token" in
            *) complete-with-words "${argvals[*]}" "$token";;
        esac
        return
    fi

    # Handle subcommand
    comment "$(printf 'Calling subcommand handler (token=%s)\n' "$token")"
    case "$token" in
		add)	on-vm-add;	status X-on-vm-add;;
		all)	on-vm-all;	status X-on-vm-all;;
		del)	on-vm-del;	status X-on-vm-del;;
		get)	on-vm-get;	status X-on-vm-get;;
		list)	on-vm-list;	status X-on-vm-list;;
		set)	on-vm-set;	status X-on-vm-set;;
        *) comment "$(printf 'unrecognized subcommand: %s\n' "$token")"; complete-with-empty ;;
    esac

    # Done
    comment "$(printf 'on-vm() - done; last token=%s' "$token")"
}   


# dylt vm add name fqdn
#
# name = complete with empty
# fqdn = complete with empty
on-vm-add () {
    # name - freeform; complete with empty
    get-token token
    status on-vm-add-1
    if on-last-token; then
        comment 'name: freeform; complete with empty'
        complete-with-empty
        return
    fi
    
    # fqdn - freeform; complete with empty
    get-token token
    status on-vm-add-2
    if on-last-token; then
        comment 'fqdn: freeform; complete with empty'
        complete-with-empty
        return
    fi

    # Done
    comment "$(printf 'on-vm-add() - done; last token=%s' "$token")"
}


# dylt vm all
#
# No args, no flags - complete with empty
on-vm-all () {
    status on-vm-all
    complete-with-empty

    # Done
    comment "on-vm-all() - done"
}


# dylt vm del name
#
# name: freeform; complete with empty
# @note name could be populated w list of vms
on-vm-del () {
    # name - freeform; complete with empty
    get-token token
    status on-vm-del-1
    if on-last-token; then
        comment 'name: freeform; complete with empty'
        complete-with-empty
        return
    fi
    
    # Done
    comment "$(printf 'on-vm-del() - done; last token=%s' "$token")"
}


# dylt vm get name
#
# name: freeform, complete with empty
# @note name could be populated w a list of vms
on-vm-get () {
    # name - freeform; complete with empty
    get-token token
    status on-vm-get-1
    if on-last-token; then
        comment 'name - freeform; complete with empty'
        complete-with-empty
        return
    fi

    # Done
    comment "$(printf 'on-vm-get() - done; last token=%s' "$token")"
}


# dylt vm list
#
# No args, no flags - complete with empty
on-vm-list () {
    status on-vm-list
    complete-with-empty

    # Done
    comment "on-vm-list() - done"
}


# dylt vm set name key val
#
# name: freeform; complete with empty
# key:  freeform; complete with empty
# val:  freeform; complete with empty
# @note name could be populated w a dynamic list of vms
# @note val could be populated w a static list of valid keys
on-vm-set () {
    # name - freeform; complete with empty
    get-token token
    status on-vm-set-1
    if on-last-token; then
        comment 'name - freeform; complete with empty'
        complete-with-empty
        return
    fi
    
    # key - freeform; complete with empty
    get-token token
    status on-vm-set-2
    if on-last-token; then
        comment 'key - freeform; complete with empty'
        complete-with-empty
        return
    fi
    
    # val - freeform; complete with empty
    get-token token
    status on-vm-set-3
    if on-last-token; then
        comment 'val - freeform; complete with empty'
        complete-with-empty
        return
    fi
    
    # Done
    comment "$(printf 'on-vm-set() - done; last token=%s' "$token")"
}


# dylt watch subcommand
#
# complete with subcommand
on-watch () {
    argvals=(script svc)
    get-token token
    status on-watch
    if on-last-token; then
        comment 'generate completions for subcommand'
        case "$token" in
            *) complete-with-words "${argvals[*]}" "$token";;
        esac
        return
    fi

    # Handle subcommand
    comment "$(printf 'Calling subcommand handler (token=%s)\n' "$token")"
    case "$token" in
		script)	on-watch-script;	status X-on-watch-script;;
		svc)	on-watch-svc;		status X-on-watch-svc;;
        *) comment "$(printf 'unrecognized subcommand: %s\n' "$token")"; complete-with-empty ;;
    esac

    # Done
    comment "$(printf 'on-watch() - done; last token=%s' "$token")"
}


# dylt watch script name targetPath
#
# name: freeform; complete with empty
# targetPath: freeform complete with empty
# @note name could be dynamically populated with a list of available scripts
on-watch-script () {
    # name - freeform; complete with empty
    get-token token
    status on-watch-script-1
    if on-last-token; then
        comment 'name - freeform; complete with empty'
        complete-with-empty
        return
    fi

    # targetPath - freeform; complete with empty
    get-token token
    status on-watch-script-2
    if on-last-token; then
        comment 'targetPath - freeform; complete with empty'
        complete-with-empty
        return
    fi

    # Done
    comment "$(printf 'on-watch-script() - done; last token=%s' "$token")"
}


# dylt watch svc name
#
# name: freeform; complete with empty
# @note name could be dynamically populated with a list of available services
on-watch-svc () {
    # name - freeform; complete with empty
    get-token token
    status on-watch-svc-1
    if on-last-token; then
        comment 'name - freeform; complete with empty'
        complete-with-empty
        return
    fi

    # Done
    comment "$(printf 'on-watch-svc() - done; last token=%s' "$token")"
}