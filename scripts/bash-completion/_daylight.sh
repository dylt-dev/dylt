#! /usr/bin/env bash

### Helper functions

# Various bash completion values, global variables, etc. Useful for development/debugging.
status () {
	printf '%-16s: DONE=%d N=%d COMP_CWORD=%d cur=[%s] <%s>\n' "$1()" "$DONE" "$N" "$COMP_CWORD" "${COMP_WORDS[COMP_CWORD]}" "${COMP_WORDS[*]}" >> /tmp/dylt.log
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


on-main () {

}

# Register completion handler with function
complete -F _dylt dylt


