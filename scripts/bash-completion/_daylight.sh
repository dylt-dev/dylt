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
	
	# Call appropriate function for latest token
	comment "$(printf "Ready for what's next: token=%s\n" "$token")"
	case $token in
		--script-path) on-call-scriptPath; status X-on-c-sp;;
		*) complete-with-empty
	esac
}


# dylt call --script-path /path/to/daylight.sh
on-call-scriptPath () {
	# --script-path value - complete with files
	get-token token
	status on-call-scriptPath-1
	if on-last-token; then
		comment "$(printf "latest token; generate completions from files")"
		complete-with-files "$token"
		return
	fi
	
	# Done
	comment "$(printf 'on-call-scriptPath() - done; last token=%s' "$token")"
}


# dylt config subcommand
on-config () {
	# complete subcommand
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
	comment "$(printf 'Ready for next step (token=%s)\n' "$token")"
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


# on host init uid gid
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













# Register completion handler with function
complete -F _dylt dylt


