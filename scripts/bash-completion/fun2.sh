# COMP_WORDS=(dylt foo)
# COMP_NWORD=1


get-token () {
	local -n _ref=$1
	_ref=${COMP_WORDS[N]}
	((N++))	
	printf 'after: N=%d\n' "$N" >> /tmp/dylt.log
	printf 'token=%s\n' "$_ref" >> /tmp/dylt.log
}

dumpargs () {
	fn=$1
	set -- "${COMP_WORDS[@]}"
	sargs=$(printf '[%s]' "$@")
	printf '%s(): (%d) <%s> (%d) %s N=%d\n' "$fn" "$COMP_CWORD" "${COMP_WORDS[*]}" "$#" "$sargs" "$N" >> /tmp/dylt.log

}

_comp () {
	echo Initializing N >> /tmp/dylt.log
	N=1
	get-token token
	on-foo

}

on-foo () {
	printf 'on-foo(): N=%d\n' "$N" >> /tmp/dylt.log
	get-token token
	on-bar
}

on-bar () {
	printf 'on-bar(): N=%d\n' "$N" >> /tmp/dylt.log
	get-token token
}


complete -F _comp f
