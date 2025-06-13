# COMP_WORDS=(dylt foo)
# COMP_NWORD=1


complete-with-empty () {
	COMPREPLY=()
	DONE=1
}

complete-with-files () {
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: complete-with-files $token\n' >&2; return 1; }
	local token=$1

	compopt -o filenames
	COMPREPLY=($(compgen -f -- "$token"))
	DONE=1
}

complete-with-words () {
	# shellcheck disable=SC2016
	(( $# == 2 )) || { printf 'Usage: complete-with-words $wordlist $token\n' >&2; return 1; }
	local wordlist=$1
	local token=$2

	COMPREPLY=($(compgen -W "$wordlist" -- "$token"))
	DONE=1
}


get-token () {
	local -n _ref=$1	
	_ref=${COMP_WORDS[N]}
	((N++))	
	printf '%-16s: token=%s N=%d\n' "get-token()" "$_ref" "$N" >> /tmp/dylt.log
}

on-last-token () { 
	(( N > COMP_CWORD ))
}

peek-token () {
	token=${COMP_WORDS[N]}
	printf 'peek-token(): token=%s\n' "$token" >> /tmp/dylt.log
}

white=$'\033[96m'
reset=$'\033[0m' 
comment () {
	printf '%s%s%s\n' "$white" "$1" "$reset" >> /tmp/dylt.log
}

status () {
	printf '%-16s: DONE=%d N=%d COMP_CWORD=%d cur=[%s] <%s>\n' "$1()" "$DONE" "$N" "$COMP_CWORD" "${COMP_WORDS[COMP_CWORD]}" "${COMP_WORDS[*]}" >> /tmp/dylt.log

}

_f () {
	echo Initializing N >> /tmp/dylt.log
	N=1
	DONE=0
	status '_f'
	on-f
	printf '\n\n' >> /tmp/dylt.log
}

on-f () {
	# If we're here, the user has entered in the command name plus whitespace so that the
	# command has been tokenized.
	#
	# This means the next token is either in progress, or it isn't.
	# If it's in progress, then COMP_CWORD=N.
	# Else, the next token has been completed. We inspect it, and move on.
	local cmds=(call config get host init list misc vm watch)
	local flags=()
	get-token token
	status on-f
	if on-last-token; then
		printf "current token is in progress: no more looking" >> /tmp/dylt.log
		case $token in
			-*) COMPREPLY=($(compgen -W "${flags[*]}" -- "$token"))
			    return;;
			*)  COMPREPLY=($(compgen -W "${cmds[*]}" -- "$token"))
			    return;;
		esac
	else
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
			*) COMPREPLY=()
			   DONE=1
			   ;;
		esac
	fi

}

# f call
#
# Some possibilities
#
# f call_
#        ^just finished, ready to accept keystrokes and COMPREPLY
#
# f call_ -
#          ^just began a flag, can COMPREPLY w flags
#
# f call --script-path
#                     ^technically not done with this flag
#
# f call --script-path_
#                      ^just finished the flag; ready to move to flag-handling func
#                      
# f call --script-path /path/to/script.sh_
#                                         ^the flag handler should have processed this flag, and
#                                          returned control to us, and we treat this as in-progress
on-call () {
	local cmds=(foo bar bum)
	local flags=(--script-path)
	get-token token
	status on-call
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		case $token in
			-*) complete-with-words "${flags[*]}" "$token";;
			*)  complete-with-words "${cmds[*]}" "$token";;
		esac
	else
		comment "$(printf "Ready for what's next: token=%s\n" "$token")"
		case $token in
			--script-path)
				on-call-scriptpath
				status X-on-c-sp
				if ((DONE == 0)); then
					get-token token
					status on-call-2
					case $token in
						-*) COMPREPLY=($(compgen -W "${flags[*]}" -- "$token"))
							DONE=1
							;;
						*)  COMPREPLY=($(compgen -W "${cmds[*]}" -- "$token"))
							DONE=1
							;;
					esac
				fi
				;;
			*) COMPREPLY=()
			   DONE=1;;
		esac
	fi
}

# f call --script-path /path/to/daylight.sh
#
# 3 possiblities here
#   - --script-path has been entered, user is working on flagval    N=4 COMP_CWORD=3 cur=[/opt/bin/day] <f call --script-path /opt/bin/day>
#   - flagval is complete, user is ready to return back to command  N=4 COMP_CWORD=4 cur=[] <f call --script-path /opt/bin/daylight.sh >
#   - user has definitely moved on                                  N=4 COMP_CWORD=6 cur=[] <f call --script-path /opt/bin/daylight.sh /opt/bin/daylight.sh update-and-restart >
on-call-scriptpath () {
	get-token token
	status on-call-scriptpath
	if (( N > COMP_CWORD )); then
		comment "$(printf "current token is in progress: COMPREPLY can be generated")"
		complete-with-files "$token"
	else
		comment "$(printf 'flagval is complete; token=' "$token")"
	fi
}

on-config () {
	local cmds=(get set show)
	local flags=()
	
	get-token token
	status on-config
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		case $token in
			-*) complete-with-words "${flags[*]}" "$token";;
			*)  complete-with-words "${cmds[*]}" "$token";;
		esac
	else
		case $token in
			get)  on-config-get;  status X-on-config-get;;
			set)  on-config-set;  status X-on-config-set;;
			show) on-config-show; status X-on-config-show;;
			*)    complete-with-empty;;
		esac
	fi
}


on-config-get () {
	local cmds=()
	local flags=()
	local keys=(name age luckyNumber)

	# First arg: config key
	get-token token
	status on-config-get
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		case $token in
			-*) complete-with-words "${flags[*]}" "$token";;
			*)  complete-with-words "${keys[*]}" "$token";;
		esac
		return
	fi
	
	comment "$(printf 'config-get is complete; token=' "$token")"
}


on-config-set () {
	local cmds=()
	local flags=()
	local keys=(name age luckyNumber)

	# first arg: config key
	get-token token
	status on-config-set-1
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		case $token in
			-*) complete-with-words "${flags[*]}" "$token";;
			*)  complete-with-words "${keys[*]}" "$token";;
		esac
		return
	fi
	
	# second arg: config value
	get-token token
	status on-config-set-2
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		complete-with-empty
		return
	fi

	comment "$(printf 'config-set is complete; token=' "$token")"
}


on-config-show () {
	status on-config-show
	# this is a terminal state. we COMPREPLY=() and return
	complete-with-empty
}

on-get () {
	status on-get
	complete-with-empty
}


on-host () {
	cmds=(init)
	flags=()
	
	# subcmd
	get-token token
	status on-host
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		case $token in
			-*) complete-with-words "${flags[*]}" "$token";;
			*)  complete-with-words "${cmds[*]}" "$token";;
		esac
		return
	fi

	# there are tokens ahead, so switch on subcommand
	case "$token" in
		init) on-host-init; status X-on-host-init;;
		*)    complete-with-empty;;
	esac
}


on-host-init () {
	# first arg: uid (no help)
	get-token token
	status on-host-init-1
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		complete-with-empty
		return
	fi
	
	# second arg: gid (no help)
	get-token token
	status on-host-init-2
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		complete-with-empty
		return
	fi

	comment "$(printf 'on-host-init is complete; token=' "$token")"
}


on-init () {
	cmds=()
	flags=(--etcd-domain)
	
	# 1 - flag
	get-token token
	status on-init-1
	if on-last-token; then
		comment "$(printf "current token is in progress: no more looking")"
		case $token in
			-*) complete-with-words "${flags[*]}" "$token";;
			*)  complete-with-empty
		esac
		return
	fi

	# 2 - flag value (complete with empty)
	get-token token
	status on-init-2
	if on-last-token; then
		complete-with-empty
		return
	fi

	comment "$(printf 'on-init is complete; last token=' "$token")"
}


on-list () {
	status on-list
	# this is a terminal state. we COMPREPLY=() and return
	complete-with-empty
}

# on-foo () {
	# # We got here as a result of 3 possiblities
	# # - Ready for subcommand
	# # - Within flag
	# # - Within flagval
	# # A tricky edgecase is within a flagval, particularly when a flagval might be a valid subcommand. context/state is necessary
	# status 'on-foo'
	# cmds=(bar bum)
	# flags=(--value)
	# get-token token
	# local cur=${COMP_WORDS[COMP_CWORD]}
	# if [[ -z "$cur" ]]; then
		# cur=${COMP_WORDS[COMP_CWORD-1]}
	# fi
	# case $cur in
		# --value) on-foo-flag
			 # return;;
		# -*) COMPREPLY=($(compgen -W "${flags[*]}" -- "$cur"))
		    # return ;;
		# *) COMPREPLY=($(compgen -W "${cmds[*]}" -- "$cur"))
		   # return ;;
	# esac
# }

# on-foo-flag () {
	# status 'on-foo-flag'
	# get-token token
	# local cur=${COMP_WORDS[COMP_CWORD]}
		
# }

complete -F _f f
