# Determines the first non-option word of the command line. This
# is usually the command
_sy_get_firstword() {
	local firstword i

	firstword=
	for ((i = 1; i < ${#COMP_WORDS[@]}; ++i)); do
		if [[ ${COMP_WORDS[i]} != -* ]]; then
			firstword=${COMP_WORDS[i]}
			break
		fi
	done

	echo $firstword
}

# Determines the last non-option word of the command line. This
# is usally a sub-command
_sy_get_lastword() {
	local lastword i

	lastword=
	for ((i = 1; i < ${#COMP_WORDS[@]}; ++i)); do
		if [[ ${COMP_WORDS[i]} != -* ]] && [[ -n ${COMP_WORDS[i]} ]] && [[ ${COMP_WORDS[i]} != $cur ]]; then
			lastword=${COMP_WORDS[i]}
		fi
	done

	echo $lastword
}

cmdsDylt=(
    call
    config
    get
    host
    init
    list
    misc
    vm
    watch
)            

cmdsDyltConfig=(
    get
    set
    show
)

_dylt () {
	local cur prev firstword lastword complete_words complete_options

	# Don't break words at : and =, see [1] and [2]
	COMP_WORDBREAKS=${COMP_WORDBREAKS//[:=]}

	cur=${COMP_WORDS[COMP_CWORD]}
	prev=${COMP_WORDS[COMP_CWORD-1]}
	firstword=$(_sy_get_firstword)
	lastword=$(_sy_get_lastword)
    
    COMPREPLY=()

    # As a convenience, assign COMP_WORDS to $@. This lets us consume args via `shift`
    set -- "${COMP_WORDS[@]}"
    shift
    do-dylt $@
}

# dylt
do-dylt () {
#    printf '$COMP_CWORD=%s\n' "$COMP_CWORD"
    local cword=${COMP_CWORD}
    # Confirm that we are here because `dylt` was the first word
    if (( cword==1 )); then
        COMPREPLY=( $(compgen -W "${cmdsDylt[*]}" -- "$cur") )
    else
        # $1 is the next token after dylt
        case $1 in 
            call) do-dylt-call; return;;
			config) do-dylt-config; return;;
			get) do-dylt-get; return;;
			host) do-dylt-host; return;;
			init) do-dylt-init; return;;
			list) do-dylt-list; return;;
			misc) do-dylt-misc; return;;
			vm) do-dylt-vm; return;;
			watch) do-dylt-watch; return;;
            *) echo MEAT
        esac
    fi
}

do-dylt-call () {
	echo call
}


# dylt config
do-dylt-config () {
	echo config
    # at this point we have parsed `dylt config`
    # that explains how we got here. but it doesn't say much about where we are
    # possibilities
    #   - just started
    #   - in the middle of a flag
    #   - in the middle of a flag's parameter
    #   - completed a flag's parameter
    # sy --loglevel achieves this. let's see how.
}


do-dylt-get () {
    :
}


do-dylt-host () {
	echo host
}


do-dylt-init () {
	echo init
}


do-dylt-list () {
    COMPREPLY=()
}


do-dylt-misc () {
	echo misc
}


do-dylt-vm () {
	echo vm
}


do-dylt-watch () {
	echo watch
}

complete -F _dylt dylt