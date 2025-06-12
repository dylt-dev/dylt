# COMP_WORDS=(dylt foo)
# COMP_NWORD=1


get-token () {
	local -n _ref=$1	
	_ref=${COMP_WORDS[N]}
	printf 'get-token(): token=%s\n' "$_ref" >> /tmp/dylt.log
	((N++))	
	printf 'after: N=%d\n' "$N" >> /tmp/dylt.log
}

peek-token () {
	token=${COMP_WORDS[N]}
	printf 'peek-token(): token=%s\n' "$token" >> /tmp/dylt.log
}

dumpargs () {
	printf '(%d) <%s> N=%d cur=[%s]\n' "$COMP_CWORD" "${COMP_WORDS[*]}" "$N" "${COMP_WORDS[COMP_CWORD]}" >> /tmp/dylt.log

}

_f () {
	echo Initializing N >> /tmp/dylt.log
	N=1
	dumpargs
	on-f
}

on-f () {
	# If we're here, the user has entered in the command name plus whitespace so that the
	# command has been tokenized.
	#
	# This means the next token is either in progress, or it isn't.
	# If it's in progress, then COMP_CWORD=N.
	# Else, the next token has been completed. We inspect it, and move on.
	local cmds=(call config get list)
	local flags=()
	get-token token
	if (( N > COMP_CWORD )); then
		printf "current token is in progress: no more looking" >> /tmp/dylt.log
		case $token in
			-*) COMPREPLY=($(compgen -W "${flags[*]}" "$token"))
			    return;;
			*)  COMPREPLY=($(compgen -W "${cmds[*]}" "$token"))
			    return;;
		esac
	else
		printf "Ready for what's next: token=%s\n" "$token" >> /tmp/dylt.log
		case $token in
			list) on-list; return;;
			*) COMPREPLY=(); return;;
		esac
	fi

}

on-list () {
	# this is a terminal state. we COMPREPLY=() and return
	COMPREPLY=()
}






on-foo () {
	# We got here as a result of 3 possiblities
	# - Ready for subcommand
	# - Within flag
	# - Within flagval
	# A tricky edgecase is within a flagval, particularly when a flagval might be a valid subcommand. context/state is necessary
	dumpargs 'on-foo'
	cmds=(bar bum)
	flags=(--value)
	get-token token
	local cur=${COMP_WORDS[COMP_CWORD]}
	if [[ -z "$cur" ]]; then
		cur=${COMP_WORDS[COMP_CWORD-1]}
	fi
	case $cur in
		--value) on-foo-flag
			 return;;
		-*) COMPREPLY=($(compgen -W "${flags[*]}" -- "$cur"))
		    return ;;
		*) COMPREPLY=($(compgen -W "${cmds[*]}" -- "$cur"))
		   return ;;
	esac
}

on-foo-flag () {
	dumpargs 'on-foo-flag'
	get-token token
	local cur=${COMP_WORDS[COMP_CWORD]}
		
}

complete -F _f f
