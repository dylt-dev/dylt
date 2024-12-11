# Some completions use a bash function for their source of words.
# This is actually hard to do.
# 
# compgen does have a -F option. It's natural to assume the compgen -F function lets you use
# a function as a word source, just as -W lets you use an array of words, etc. But this isn't
# the case. -F has a very strange implementation that is never useful. It only exists because
# complete has an -F option, which is in fact useful, and the authors wanted complete and 
# compgen have the same set of options. So compgen has at least two options -- -F and -C -- 
# that are not useful, and only exist so that compgen and complete's options match. You could
# say that they exist for - ahem - completeness.
#
# We definitely want to be able to use bash functions or caommands as completion word sources.
# If we can't use compgen -F, or anything that already exists, we'll need to build the support
# ourselves. We'll build around compgen -W. compgen -W is great at taking an array of words, and
# the current partial word, and returning all matches.
#
# That gets us most of the way there. Unfortunately compgen -W returns a list of words one per line,
# and complete & COMPREPLY prefer an array. So we need to convert compgen's response into an array using
# mapfile. This is a bit confusing, since the initial source of words for compgen is typically a one-per-line
# list of words that also needs to be converted to mapfile. So the basic flow looks like this:
#   - Call function or command to produce list of words
#   - Convert wordlist to an array with mapfile
#   - Call compgen to filter word array on COMP_CWORD
#   - Convert filtered wordlist to an array with mapfile
#   - Set COMPREPLY to array
#
# The last two steps are naturally combined, because the last mapfile can take COMPREPLY as an argument. But the other
# steps are a bit challenging to combine, since functions can't return arrays and passing arrays as namerefs is not 
# composable. So what's the best way to make it easy to map an arbitrary command's output to an array, then filter it
# with compgen, then mapfile to COMPREPLY? Maybe something like ...
#
#	mapfile -t COMREPLY <(cmd | func-that-mapfiles-stdin-and-calls compgen -W)
#
# I'm not sure if there's any way to reduce this. But it might be worth trying to play with

# write stdin to tmpfile
# map tmpfile to an array (maybe tee can handle this, maybe not)
# compgen array into a tmpfile
# mapfile tmpfile into COMPREPLY
read-into-compreply ()
{
	local curr=${COMP_WORDS[COMP_CWORD]}
	local last=${COMP_WORDS[COMP_CWORD-1]}

	local tmpStdin; tmpStdin=$(mktemp --tmpdir bc.ric.stdin.XXXXXX) || return
	cat >"tmpStdin" || return 
	local words; mapfile -t words <"tmpStdin" || return
	local tmpCompgen; tmpCompgen=$(mktemp --tmpdir fc.ric.compgen.XXXXXX) || return
	compgen -W "${words[*]}" -- "$curr" >"$tmpCompgen" || return
	mapfile -t COMPREPLY <"$tmpCompgen" || return
}


_sunbeam-sh ()
{
	local curr=$2
	local last=$3

	print-comp-args "$@"
	local mainCmds=(gen-nightly-tagname \
 					 gen-nightly-timestamp \
 					 git-do-nightly-release \
 					 git-download-latest-daylightsh \
 					 git-get-latest-release-spec \
 					 git-get-latest-release-tag \
 					 git-get-latest-release-version \
 					 git-install-latest-daylightsh \
 					 git-install-latest-dylt \
 					 git-tag-nightly \
 				     yesorno \
					)
	local lastCmd=${last##*/}
	case "$lastCmd" in
		sunbeam.sh)
			mapfile -t COMPREPLY < <(compgen -W "${mainCmds[*]}" -- "$curr")
			;;
	esac
	exec 10>&-
}


print-comp-args ()
{
	exec 10>/tmp/sunbeam.sh.fc.txt
	printf 'cmd=%q last=%q curr=%q\n' "$1" "$3" "$2">&10
	printf '%-25s %s\n' COMP_WORDS "$(printf '<%q>' "${COMP_WORDS[@]}")" >&10
	printf '%-25s %d\n' COMP_CWORD "$COMP_CWORD" >&10
	# shellcheck disable=SC2016
	printf '%-25s %s\n' '${COMP_WORDS[COMP_CWORD]}' "${COMP_WORDS[COMP_CWORD]}" >&10
	printf '%-25s %s\n' '${COMP_WORDS[COMP_CWORD-1]}' "${COMP_WORDS[COMP_CWORD-1]}" >&10
	printf '%-25s %d\n' COMP_KEY "$COMP_KEY" >&10
	printf '%-25s %s\n' COMP_LINE "$COMP_LINE" >&10
	printf '%-25s %d\n' COMP_POINT "$COMP_POINT" >&10
	printf '%-25s %d\n' COMP_TYPE "$COMP_TYPE" >&10
	echo >&10
}


complete -F _sunbeam-sh sunbeam.sh
