#! /usr/bin/env bash


#------------------------------------------------------------------------------
#
# add-to-bashrc()
#
# Add a sunbeam() function to bashrc that calls the sunbeam script with "$@"
#
add-to-bashrc ()
{
    local bashrc="$HOME/.bashrc"
    local funcName='sunbeam'
    local sunbeamPath="$HOME/.local/bin/sunbeam.sh"

    if grep -q "^${funcName}()" "$bashrc" 2>/dev/null; then
        printf '`%s` function already exists in %s — nothing to do\n' "$funcName" "$bashrc"
        return 0
    fi

    cat >> "$bashrc" << EOF

# added by opencode in loyal service to master
$funcName()
{
    $sunbeamPath "\$@"
}
EOF

    printf 'Added `%s` function to %s\n' "$funcName" "$bashrc"
    printf 'Restart your shell or run: source %s\n' "$bashrc"
    printf 'Then type: %s --help\n' "$funcName"
}


#-------------------------------------------------------------------------------
#
# detect-platform()
#
# Detect the OS and architecture as a platform string
#
detect-platform ()
{
    (( $# == 0 )) || { printf 'Usage: detect-platform\n' >&2; return 1; }
    local os arch

    case "$(uname -s)" in
        Linux)                     os="linux" ;;
        Darwin)                    os="darwin" ;;
        MINGW*|MSYS*|CYGWIN*)     os="windows" ;;
        *) printf 'Unsupported OS: %s\n' "$(uname -s)" >&2; return 1 ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)             arch="amd64" ;;
        aarch64|arm64)            arch="arm64" ;;
        armv7l|armv6l)            arch="arm" ;;
        *) printf 'Unsupported architecture: %s\n' "$(uname -m)" >&2; return 1 ;;
    esac

    printf '%s-%s' "$os" "$arch"
}


#-------------------------------------------------------------------------------
#
# download-dylt()
#
# Download latest dylt release
#
download-dylt ()
{
    # basic validation before loading remote functions
    # shellcheck disable=SC2016
    (( $# >= 1 )) || { printf 'Usage: download-dylt $dstFolder [$platform]\n' >&2; return 1; }
    source-github-utils || return
    # parse github args
    local -A argmap=()
    local nargs=0
    github-parse-args argmap nargs "$@" || return
    local -a flags=()
    github-create-flags argmap flags token
    shift "$nargs"
    local dstFolder=$1
    local platform=$2
    [[ -d "$dstFolder" ]] || { echo "Non-existent folder: $dstFolder" >&2; return 1; }

    if [[ -z "$platform" ]]; then
        platform=$(detect-platform) || return
    fi

    local -a flags=()
    [[ -n "${argmap[token]+exists}" ]] && flags+=(--token "${argmap[token]}") 

    local version; version=$(github-release-get-latest-tag "${flags[@]}" dylt-dev dylt) || return

    local releaseName="dylt_${platform}.tar.gz"
    local legacyName="dylt_$(dylt-legacy-platform "$platform").tar.gz"
    local url="https://github.com/dylt-dev/dylt/releases/download/$version/$releaseName"

    if curl --fail --location --head "$url" >/dev/null 2>&1; then
        github-release-download-latest "${flags[@]}" dylt-dev dylt "$releaseName" "$dstFolder" || return
    else
        github-release-download-latest "${flags[@]}" dylt-dev dylt "$legacyName" "$dstFolder" || return
    fi
}


#-------------------------------------------------------------------------------
#
# download-github-utils()
#
# Download github-utils.sh from daylight and write to stdout
# Safe to pipe to source or redirect to a file
# Returns 1 on failure
#
download-github-utils ()
{
    (( $# == 0 )) || { printf 'Usage: download-github-utils\n' >&2; return 1; }
    safe-curl "https://raw.githubusercontent.com/daylight-public/daylight/main/github-utils.sh"
}


#------------------------------------------------------------------------------
#
# gen-nightly-tagname()
#
# Create a tagname from a root version (eg v1.0.7) and the nightly timestamp
# eg v1.0.7-nightly.20240522152147
#
gen-nightly-tagname ()
{
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: gen-nightly-tagname $version\n' >&2; return 1; }
	local version=$1

	local ts; ts=$(gen-nightly-timestamp) || return
	local label=$(printf '%s-nightly.%s' "$version" "$ts")

	printf "$label"
}


#------------------------------------------------------------------------------
#
# gen-nightly-timestamp()
#
# Simple function to create a timestamp in the format used for nightly release tags
# eg 20240522152147 (Wed May 22 15:21:47 CDT 2024)
# More separators for the various date fields would be great, but they would break https://semver.org rules
#
gen-nightly-timestamp ()
{
	date '+%Y%m%d%H%M%S'
}


#------------------------------------------------------------------------------
#
# git-do-nightly-release()
#
# Tag the nightly release, push the current commit, and push the tag
#
git-do-nightly-release ()
{
	# @todo Use [[ $(git status --porcelain) == "" ]] to see if there is uncommited work. If so ask for confirmation
	# shellcheck disable=SC2016
	# shellcheck disable=SC2016
	{ (( $# >= 0 )) && (( $# <= 1 )); } || { printf 'Usage: git-do-nightly-release [$version]\n' >&2; return 1; }
	local version=${1:-''}

	if [[ -z "$version" ]]; then
		version=$(git-get-latest-release-version dylt-dev dylt) || return
		printf '$version=%s\n' "$version"
	fi
	if [[ $(git status --porcelain) != "" ]]; then
		printf '%s\n' "There are uncommitted changes"
		if ! yesorno yn "Push nightly release anyway? "; then 
			return 0 
		fi
	fi
	local tag
	git-tag-nightly tag "$version" || return
# 	day ssh h0 GOBIN=/opt/bin /usr/local/go/bin/go install github.com/dylt-dev/dylt@$tag
#	day ssh h1 GOBIN=/opt/bin /usr/local/go/bin/go install github.com/dylt-dev/dylt@$tag
#	day ssh h2 GOBIN=/opt/bin /usr/local/go/bin/go install github.com/dylt-dev/dylt@$tag
}


#------------------------------------------------------------------------------
#
# git-download-latest-daylightsh()
#
# Download the latest daylight.sh from github
#
# By default, download to ~/tmp/
#
git-download-latest-daylightsh ()
{
	# shellcheck disable=SC2016
	{ (( $# >= 0 )) && (( $# <= 1 )); } || { printf 'Usage: git-download-latest-daylightsh [$downloadFolder]\n' >&2; return 1; }
	local downloadFolder=${1:-"$HOME/tmp"}
	[[ -d "$downloadFolder" ]] || { echo "Non-existent folder: $downloadFolder" >&2; return 1; }

	url=https://raw.githubusercontent.com/daylight-public/daylight/main/daylight.sh
	curl --silent --remote-name --output-dir "/$downloadFolder/" "$url"
}


#------------------------------------------------------------------------------
#
# git-get-latest-release-spec()
#
# get the lastest release spec, ie github.com/$owner/$repo/$latestTag
#
git-get-latest-release-spec ()
{
	# shellcheck disable=SC2016
	(( $# == 2 )) || { printf 'Usage: git-get-latest-release-tag $owner $repo\n' >&2; return 1; }
	local owner=$1
	local repo=$2
	local tag; tag=$(git-get-latest-release-tag "$owner" "$repo") || return
	local spec; spec=$(printf '%s/%s/%s/%s' "github.com" "$owner" "$repo" "$tag") || return
	printf '%s' "$spec"
}


#------------------------------------------------------------------------------
#
# git-get-latest-release-tag()
#
# Use the GitHub API to get the tag of the latest version of a GitHub release,
# for a specified owner+repo
# This can be used to `go install` a specific releae
# Note that GitHub defines 'latest version' as the release that was created most recently,
# unlike https://semver.org, which has complicated rules to define the most recent release.
#
git-get-latest-release-tag ()
{
	# shellcheck disable=SC2016
	(( $# == 2 )) || { printf 'Usage: git-get-latest-release-tag $owner $repo\n' >&2; return 1; }
	local owner=$1
	local repo=$2
	local tag; tag=$(curl -L --silent "api.github.com/repos/$owner/$repo/releases/latest" | jq -r .tag_name)
	printf '%s' "$tag"
}


#------------------------------------------------------------------------------
#
# git-get-latest-release-version()
#
# Get the lastest release version, which is derived from the getting latest tag
# and then dropping everything after the hyphen
#
git-get-latest-release-version ()
{
	# shellcheck disable=SC2016
	(( $# == 2 )) || { printf 'Usage: git-get-latest-release-tag $owner $repo\n' >&2; return 1; }
	local owner=$1
	local repo=$2

	local releaseTag; releaseTag=$(git-get-latest-release-tag "$owner" "$repo") || return
	local releaseVersion=${releaseTag%%-*}
	printf '%s' "$releaseVersion" || return

}


#------------------------------------------------------------------------------
#
# git-install-latest-daylightsh()
#
# Install the latest daylight.sh on a VM
#
# By default, install from github
# If $scriptPath is specified, install that one instead
#
git-install-latest-daylightsh ()
{
	# shellcheck disable=SC2016
	# shellcheck disable=SC2016
	{ (( $# >= 1 )) && (( $# <= 2 )); } || { printf 'Usage: git-install-latest-daylightsh $remoteHost [$scriptPath]\n' >&2; return 1; }
	local remoteHost=$1
	ssh ubuntu@$remoteHost -- mkdir -p /opt/bin
	ssh ubuntu@$remoteHost -- 'if [[ -f /opt/bin/daylight.sh ]]; then cp /opt/bin/daylight.sh /opt/bin/daylight.sh.bk; fi'
	local scriptPath
	if (( $# == 2 )); then
		scriptPath=$2
		[[ -f "$scriptPath" ]] || { echo "Non-existent path: $scriptPath" >&2; return 1; }
	else
		downloadFolder="$HOME/tmp"
		git-download-latest-daylightsh "$downloadFolder"
		scriptPath=/tmp/daylight.sh
	fi
	scp "$scriptPath" "ubuntu@$remoteHost:/opt/bin/daylight.sh"
}


#------------------------------------------------------------------------------
#
# git-install-latest-dylt()
#
# Install the latest dylt based on its release tag
#
git-install-latest-dylt ()
{
	local owner=dylt-dev
	local repo=dylt
	local tag=$(git-get-latest-release-tag "$owner" "$repo")
	local release=github.com/$owner/$repo@$tag
	go install "$release"
}	


#------------------------------------------------------------------------------
#
# git-tag-nightly()
#
# Create the nightly tagname, and then a git tag from the name
#
git-tag-nightly ()
{
	# shellcheck disable=SC2016
	(( $# == 2 )) || { printf 'Usage: git-tag-nightly varname $version\n' >&2; return 1; }
	# shellcheck disable=SC2178
	[[ $1 != tagname ]] && { local -n tagname; tagname=$1; }
	local version=$2

	tagname=$(gen-nightly-tagname $version) || return
	git tag "$tagname"
	git push
	git push --tags
	local releaseSpec="github.com/dylt-dev/dylt@$tagname"
	go list -m "$releaseSpec"
	echo "$tagname"
}


# #-------------------------------------------------------------------------------
# #
# # github-app-get-client-id()
# #
# # Get the OAuth client ID for a GitHub App
# #
# github-app-get-client-id ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 1 )) || { printf 'Usage: github-app-get-id $appSlug\n' >&2; return 1; }
#     local appSlug=$1
# 
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
#     local -A info
#     github-app-get-info "${flags[@]}" info "$appSlug" || return
#     local clientId=${info[client_id]}
#     printf '%s' "$clientId"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-app-get-data()
# #
# # Get GitHub App installation data from the API
# #
# github-app-get-data ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 1 )) || { printf 'Usage: github-app-get-data $appSlug\n' >&2; return 1; }
#     local appSlug=$1
# 
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
#     github-curl "${flags[@]}" "/apps/$appSlug" || return
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-app-get-id()
# #
# # Get the ID of a GitHub App
# #
# github-app-get-id ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 1 )) || { printf 'Usage: github-app-get-id $appSlug\n' >&2; return 1; }
#     local appSlug=$1
# 
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
#     local -A info
#     github-app-get-info "${flags[@]}" info "$appSlug" || return
#     local id=${info[id]}
#     printf '%s' "$id"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-app-get-info()
# #
# # Get detailed info about a GitHub App
# #
# github-app-get-info ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 2 )) || { printf 'Usage: github-app-get-data $infovar $appSlug\n' >&2; return 1; }
#     local -n _info=$1
#     local appSlug=$2
# 
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
#     local tmpCurl; tmpCurl=$(mktemp --tmpdir curl.XXXXXX) || return
#     github-app-get-data "${flags[@]}" "$appSlug" >"$tmpCurl" || return
#     local tmpJq; tmpJq=$(mktemp --tmpdir jq.XXXXXX) || return
#     jq -r '[.id, .client_id, .slug] | @tsv' <"$tmpCurl" >"$tmpJq" || return
#     read -r -a args < "$tmpJq" || return
# 
#     _info[id]=${args[0]}
#     _info[client_id]=${args[1]}
#     # shellcheck disable=SC2154
#     _info[slug]=${args[2]}
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-create-flags()
# #
# # Create curl flags from a parsed argument map
# #
# github-create-flags ()
# {
#     # shellcheck disable=SC2016
#     (( $# >=2 )) || { printf 'Usage: github-create-flags argmap flags [$flag1 $flag2 ... $flagn]\n' >&2; return 1; }
#     # Check that argmap is either an assoc array or a nameref to an assoc array
#     [[ $1 != argmap ]] && { local -n argmap; argmap=$1; }
#     [[ $(declare -p argmap 2>/dev/null) == "declare -A"* ]] \
#     || [[ $(declare -p "${!argmap}" 2>/dev/null) == "declare -A"* ]] \
#     || { printf "%s is not an associative array, and it's not a nameref to an associative array either\n" "argmap" >&2; return 1; }
#     # Check that flags is either an array or a nameref to an array
#     [[ $2 != flags ]] && { local -n flags; argmap=$2; }
#     [[ $(declare -p flags 2>/dev/null) == "declare -a"* ]] \
#     || [[ $(declare -p "${!flags}" 2>/dev/null) == "declare -a"* ]] \
#     || { printf "%s is not an array, and it's not a nameref to an array either\n" "flags" >&2; return 1; }
# 
#     flags=()
#     local argname arg
#     shift 2
#     if (( $# == 0 )); then
#         for argname in "${!argmap[@]}"; do
#             arg=${argmap["$argname"]}
#             flags+=("$argname" "$arg")
#         done
#     else
#         while (( $# > 0 )); do
#             argname=$1
#             if [[ -v argmap["$argname"] ]]; then
#                 arg=${argmap["$argname"]}
#                 flags+=("--${argname}" "$arg")
#             fi
#             shift
#         done
#     fi
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-create-url()
# #
# # Create a full GitHub API URL from a path
# #
# github-create-url ()
# {
#     local urlPath=$1
#     local urlBase=${2:-'https://api.github.com'}
# 
#     # Trim leading slash
#     if [[ $urlPath == /* ]]; then
#         urlPath=${urlPath:1}
#     fi
#     # concatenate urlBase and Path
#     local url="$urlBase/$urlPath"
#     printf '%s' "$url" || return
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-create-user-access-token()
# #
# # Create a GitHub user access token via API
# #
# github-create-user-access-token ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 2 )) || { printf 'Usage: github-create-user-access-token tokenvar $appslug\n' >&2; return 1; }
#     # shellcheck disable=SC2178
#     [[ $1 != tokenvar ]] && { local -n tokenvar; tokenvar=$1; }
#     local appSlug=$2
# 
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
#     
#     # Get the clientId for the dylt-cli GitHub App CLI, which must be installed 
#     local clientId; clientId=$(github-app-get-client-id "${flags[@]}" "$appSlug") || return
# 
#     # Use client id to invoke device code flow
#     flags+=(--data '')
#     urlPath="/login/device/code?client_id=$clientId"
#     urlBase="https://github.com"
#     local -a args
#     read -r -a args < <(github-curl "${flags[@]}" "$urlPath" "$urlBase" \
#                         | jq -r '[.device_code, .user_code, .verification_uri] | @tsv') \
#                         || { printf 'Call failed: github-curl()\n'; return; }
#     local deviceCode=${args[0]}
#     local userCode=${args[1]}
#     local verificationUri=${args[2]}
# 
#     # Prompt user to do stuff in the browser
#     echo
#     printf '%-40s%s\n' "User Code" "$userCode"
#     printf '%-40s%s\n' "Verification Uri" "$verificationUri"
#     if command -v pbcopy >/dev/null; then
#         printf '%s' "$userCode" | pbcopy
#     fi
#     if command -v open >/dev/null; then
#         echo
#         read -r -p "Hit <Enter> to open $verificationUri in your browser ..." _
#         open "$verificationUri"
#     fi
#     echo
#     
#     # Post to the thing and grab the access token
#     local prompt; prompt=$(printf 'Go to %s and enter %s. Then return here and press <Enter> ...' "$verificationUri" "$userCode") || return
#     read -r -p  "$prompt" _
#     local grantType='urn:ietf:params:oauth:grant-type:device_code'
#     urlPath="$(printf '/login/oauth/access_token?client_id=%s&device_code=%s&grant_type=%s' "$clientId" "$deviceCode" "$grantType")"
#     urlBase="https://github.com"
#     read -r -a args < <(github-curl "${flags[@]}" "$urlPath" "$urlBase" \
#                         | jq -r '[.access_token] | @tsv') \
#                         || return
#     # return the access token
#     # shellcheck disable=SC2034
#     tokenvar=${args[0]}
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-curl()
# #
# # Make an authenticated request to the GitHub API
# #
# github-curl ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@"
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# >= 1 && $# <= 2 )) || { printf 'Usage: github-curl [flags] $urlPath [$urlBase]\n' >&2; return 1; }
#     local urlPath=${1##/} # Trim leading slash if necessary
#     local urlBase=${2:-'https://api.github.com'}
# 
#     # Set headers to flag values or default
#     local acceptDefault='application/vnd.github+json'
#     local accept=${argmap[accept]:-$acceptDefault}
#     local outputDefault='-'
#     local output=${argmap[output]:-$outputDefault}
#     # Set url and token, if present
#     local url="$urlBase/$urlPath"
#     # Can't really parameterize on token -- we need separate curl calls for with token, and without
#     local -a flags=(--fail-with-body --location --silent)
#     flags+=(--header "Accept: $accept")
#     flags+=(--output "$output")
#     [[ -v argmap[data] ]] && flags+=(--data "$(printf "'%s'" "${argmap[data]}")")
#     local tokenVal
#     if [[ -v argmap[token] ]]; then
#         tokenVal=${argmap[token]}
#     elif [[ -n "${GITHUB_TOKEN-}" ]]; then
#         tokenVal=$GITHUB_TOKEN
#     fi
#     [[ -n "$tokenVal" ]] && flags+=(--header "Authorization: Bearer $tokenVal")
#     curl "${flags[@]}" "$url" \
#         || { printf 'curl failed inside github-curl\n' >&2; return 1; }
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-curl-post()
# #
# # @deprecated
# # Use github-curl with --data 'your-data'
# #
# github-curl-post ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@"
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     { (( $# >= 2 )) && (( $# <= 3 )); } || { printf 'Usage: github-curl-post $urlPath $postData [$urlBase]\n' >&2; return 1; }
#     local urlPath=${1##/} # Trim leading slash if necessary
#     local postData=$2
#     local urlBase=${3:-'https://api.github.com'}
# 
#     local acceptDefault='application/vnd.github+json'
#     local accept=${argmap[accept]:-$acceptDefault}
#     local outputDefault='-'
#     local output=${argmap[output]:-$outputDefault}
#     # Set url and token, if present
#     local url="$urlBase/$urlPath"
#     local token=${argmap[token]}
#     # Can't really parameterize on token -- we need separate curl calls for with token, and without
#     if [[ -n $token ]]; then
#         curl --fail-with-body \
#              --location \
#              --silent \
#              --data "'$postData'" \
#              --header "Accept: $accept" \
#              --header "Authorization: Token $token" \
#              --output "$output" \
#              "$url" \
#         || return
#     else
#         curl --fail-with-body \
#              --location \
#              --silent \
#              --data "'$postData'" \
#              --header "Accept: $accept" \
#              --output "$output" \
#              "$url" \
#         || return
#     fi
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # dylt-legacy-platform()
# #
# # Translate current platform spec into a legacy platform spec
# #
# dylt-legacy-platform ()
# {
#     (( $# == 1 )) || { printf 'Usage: dylt-legacy-platform $canonical_platform\n' >&2; return 1; }
#     local platform=$1
#     case $platform in
#         linux-amd64)  printf 'Linux_x86_64' ;;
#         linux-arm64)  printf 'Linux_arm64'  ;;
#         linux-386)    printf 'Linux_i386'   ;;
#         linux-arm)    printf 'Linux_armv7l' ;;
#         darwin-amd64) printf 'Darwin_x86_64' ;;
#         darwin-arm64) printf 'Darwin_arm64'  ;;
#         windows-amd64) printf 'Windows_x86_64' ;;
#         windows-arm64) printf 'Windows_arm64'  ;;
#         *)            printf '%s' "$platform" ;;
#     esac
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-download-latest-release()
# #
# # Download the latest release asset from a GitHub repository
# #
# github-download-latest-release ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@"
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 4 )) || { printf 'Usage: download-latest-release $org $repo $name $downloadFolder\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local name=$3
#     local downloadFolder=$4
# 
#     # Get release package data as assoc array
#     local -A releaseInfo
#     github-get-release-package-info releaseInfo "$org" "$repo" "$name" || return
#     local url=${releaseInfo[url]}
#     local accept='Accept: application/octet-stream'
#     local output="$downloadFolder/$name"
#     local token=${argmap[token]}
#     if [[ -n $token ]]; then
#         github-curl --token "$token" --accept "$accept" --output "$output"
#     else
#         github-curl --accept "$accept" --output "$output"
#     fi
#     printf '%s' "$releasePath"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-get-release-data()
# #
# # @deprecated
# # Use github-release-get-data
# #
# github-get-release-data ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@"
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     { (( $# >= 2 )) && (( $# <= 4 )); } || { printf 'Usage: github-get-release-data [flags] $org $repo [$releaseTag [$platform]]\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local tag=${3:-""}
#     
#     local urlPath; urlPath="$(github-get-releases-url-path "$org" "$repo" "$tag")" || return
# 	local tmpCurl; tmpCurl=$(create-temp-file github.get.release.data.json) || return
#     # build argstring for github-curl
#     local argstring=''
#     [[ -n ${argmap[token]} ]] && argstring+="--token ${argmap[token]}"
#     # github-curl -- note $argstring is unquoted
#     github-curl "$argstring" "$urlPath" >"$tmpCurl" || return
# 	printf '%s' "$tmpCurl"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-get-release-name-list()
# #
# # @deprecated
# # Use github-release-get-name-list
# #
# github-get-release-name-list ()
# {
#     # shellcheck disable=SC2016
#     { (( $# >= 3 )) && (( $# <= 4 )); } || { printf 'Usage: github-get-release-name-list listVar $org $repo [$tag]\n' >&2; return 1; }
#     local -n listVar; listVar=$1
#     # ${@:2} skips the first two args, which are $0 and the $listVar nameref 
#     local tmpCurl; tmpCurl=$(github-get-release-data "${@:2}") || return
#     # shellcheck disable=SC2034
# 	local tmpJq; tmpJq=$(create-temp-file jq.get.release.name.list.txt) || return
# 	jq -r '[.assets[].name] | sort | @tsv' \
#         <"$tmpCurl" \
#         >"$tmpJq" \
#         || return
# 
# 	# shellcheck disable=SC2034
#     read -r -a listVar <"$tmpJq" || return
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-get-release-package-data()
# #
# # @deprecated
# # Use github-release-get-package-data
# #
# github-get-release-package-data ()
# {
#     # shellcheck disable=SC2016
#     (( $# == 3 )) || { printf 'Usage: github-release-package-data $org $repo $name\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local name=$3
# 
#     local urlPath; urlPath="$(github-get-releases-url-path "$org" "$repo")" || return
#     local tmpCurl; tmpCurl="$(create-temp-file 'curl.release')" || return
#     github-curl "$urlPath" >"$tmpCurl" || return
#     local tmpJq; tmpJq="$(create-temp-file 'jq.release')" || return
#     jq -r --arg name "$name" \
#        '.assets[]
#         | select(.name == $name)' \
#       </"$tmpCurl" >/"$tmpJq" \
#       || return 
#     printf '%s' "$tmpJq"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-get-release-package-info()
# #
# # @deprecated
# # Use github-release-get-package-info
# #
# github-get-release-package-info ()
# {
#     # shellcheck disable=SC2016
#     (( $# == 4 )) || { printf 'Usage: github-get-release-package-info infovar $org $repo $name\n' >&2; return 1; }
#     local -n info=$1
#     local releaseDataPath; releaseDataPath=$(github-get-release-package-data "${@:2}") || return
#     local tmpJq; tmpJq=$(create-temp-file 'jq.release.info') || return
#     jq -r '[.id, .url, .browser_download_url] | @tsv' \
#       <"$releaseDataPath" \
#       >"$tmpJq" \
#       || return
#     local -a args
#     read -r -a args <"$tmpJq" || return
#     info[id]=${args[0]}
#     info[url]=${args[1]}
#     local browser_download_url=${args[2]}
#     local filename=${browser_download_url##*/}
#     info[browser_download_url]=$browser_download_url
#     info[filename]=$filename
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-parse-args()
# #
# # Parse common GitHub API arguments into an associative array
# #
# github-parse-args ()
# {
#     # shellcheck disable=SC2016
#     (( $# >= 2 )) || { printf 'Usage: github-parse-args infovar nargs [$args]\n' >&2; return 1; }
#     # shellcheck disable=SC2178
#     [[ $1 != argmap ]] && { local -n argmap; argmap=$1; }
#     # Check that argmap is either an assoc array or a nameref to an assoc array
#     [[ $(declare -p argmap 2>/dev/null) == "declare -A"* ]] \
#     || [[ $(declare -p "${!argmap}" 2>/dev/null) == "declare -A"* ]] \
#     || { printf "%s is not an associative array, and it's not a nameref to an associative array either\n" "argmap" >&2; return 1; }
#     # shellcheck disable=SC2178
#     [[ $2 != nargs ]] && { local -n nargs; nargs=$2; }
# 
#     nargs=0
#     shift 2
#     while (( $# > 0 )); do
#         case $1 in
#             '--accept'   |\
#             '--data'     |\
#             '--output'   |\
#             '--token'    |\
#             '--platform' |\
#             '--version' \
#             )
#                 (( $# >= 2 )) || { printf -- '%s specified but no value provided.\n' "$1" >&2; return 1; }
#                 argmap["${1##--}"]=$2
#                 ((nargs+=2))
#                 shift 2
#                 ;;
#             '--')
#                 shift
#                 ((nargs++))
#                 break
#                 ;;
#             *)
#                 break
#                 ;;
#         esac
#     done
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-create-url-path()
# #
# # Create a URL path for a GitHub release
# #
# github-release-create-url-path ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# >= 2 )) || { printf 'Usage: github-release-create-url-path $org $repo\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
# 
#     local tag=${argmap[version]:-''}
#     local urlPath
#     if [[ -n "$tag" ]]; then
#         local urlPath="/repos/$org/$repo/releases/tags/$tag"
#     else
#         local urlPath="/repos/$org/$repo/releases/latest"
#     fi
# 
#     # printf with \n if interactive
#     if [[ -t 0 ]]; then
#         printf '%s\n' "$urlPath"
#     else
#         printf '%s' "$urlPath"
#     fi
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-download()
# #
# # Download a release asset from a GitHub repository
# #
# github-release-download ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 4 )) || { printf 'Usage: github-release-download $org $repo $releaseName $downloadFolder\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local name=$3
#     local downloadFolder=${4%%/}
# 
#     # Get release info
#     local -a flags=()
#     github-create-flags argmap flags token version
#     local -A releaseInfo
#     github-release-get-package-info "${flags[@]}" releaseInfo "$org" "$repo" "$name" || return
#     # download release file using releaseInfo data
#     local urlPath=${releaseInfo[urlPath]}
#     local filename=${releaseInfo[filename]}
#     local accept='Accept: application/octet-stream'
#     local output="$downloadFolder/$filename"
#     flags+=(--accept "$accept" --output "$output")
#     github-curl "${flags[@]}" "$urlPath" || return
#     printf '%s' "$output"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-download-latest()
# #
# # Download the latest release asset for a named package
# #
# github-release-download-latest ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 4 )) || { printf 'Usage: github-release-download-latest [$flags] $org $repo $name $downloadFolder\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local name=$3
#     local downloadFolder=${4%%/}
# 
#     local -a flags
#     github-create-flags argmap flags token || return
#     local version; version=$(github-release-get-latest-tag "${flags[@]}" "$org" "$repo") || return
#     flags+=(--version "$version")
#     github-release-download "${flags[@]}" "$org" "$repo" "$name" "$downloadFolder" || return
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-get-data()
# #
# # Get release data from a GitHub repository
# #
# github-release-get-data ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@"
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     { (( $# >= 2 )) } || { printf 'Usage: github-release-get-data [flags] $org $repo\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     
#     local -a flags
#     github-create-flags argmap flags version || return
#     local urlPath; urlPath=$(github-release-create-url-path "${flags[@]}" "$org" "$repo") || return
#     # build argstring for github-curl
#     github-create-flags argmap flags token || return
#     github-curl "${flags[@]}" "$urlPath" || return
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-get-latest-tag()
# #
# # Get the latest release tag from a GitHub repo
# #
# github-release-get-latest-tag ()
# {
#     command -v "jq" >/dev/null || { printf '%s is required, but was not found.\n' "jq" >&2; return 1; }
#     # shellcheck disable=SC2016
#     (( $# == 2 )) || { printf 'Usage: github-get-latest-version [flags] $org $repo\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     
#     releasesUrlPath=$(github-release-create-url-path "$org" "$repo")
#     # build flags for github-curl
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
#     # local VER; VER=$(github-curl "${flags[@]}" "$releasesUrlPath" \
#                     #  | jq -r .tag_name)
#     local tmpCurl; tmpCurl=$(mktemp --tmpdir curl.latest.tag.XXXXXX) || return
#     github-curl "${flags[@]}" "$releasesUrlPath" >"$tmpCurl" || return
#     local tmpJq; tmpJq=$(mktemp --tmpdir jq.latest.tag.XXXXXX) || return
#     jq -r '.tag_name' <"$tmpCurl" >"$tmpJq" || return
#     read -r tag < "$tmpJq" || return    
#     
#     printf '%s' "$tag"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-get-package-data()
# #
# # Get raw asset data for a named release package
# #
# github-release-get-package-data ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 3 )) || { printf 'Usage: github-release-get-package-data $org $repo $name\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local name=$3
# 
#     local -a flags
#     github-create-flags argmap flags version token || return
#     local urlPath; urlPath=$(github-release-create-url-path "${flags[@]}" "$org" "$repo") || return
#     github-create-flags argmap flags token || return
#     local tmpCurl; tmpCurl=$(mktemp --tmpdir curl.release.XXXXXX) || return
#     github-curl "${flags[@]}" "$urlPath" >"$tmpCurl" || return
#     jq -r --arg name "$name" \
#        '.assets[]
#         | select(.name == $name)' \
#       <"$tmpCurl" \
#       || return 
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-get-package-info()
# #
# # Get structured package info for a named release asset
# #
# github-release-get-package-info ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 4 )) || { printf 'Usage: github-get-release-package-info infovar $org $repo $name\n' >&2; return 1; }
#     # shellcheck disable=SC2178
#     [[ $1 != info ]] && { local -n info; info=$1; }
#     local org=$2
#     local repo=$3
#     local name=$4
# 
#     # Call github-release-get-package-data and create/parse the necesary fields
#     local -a flags=()
#     github-create-flags argmap flags token version
#     local -a fields=()
#     read -r -a fields < <(github-release-get-package-data "${flags[@]}" "$org" "$repo" "$name" \
#     | jq -r '
#         [.browser_download_url,
#          .content_type,
#          (.browser_download_url | match(".*/(.*)").captures[0].string),
#          .id,
#          .name,
#          .url,
#          (.url | match("https://api.github.com/(.*)").captures[0].string)
#         ] | @tsv' \
#       || return)
#     # Package fields into the info assoc array
#     info[browser_download_url]=${fields[0]}
#     info[content_type]=${fields[1]}
#     info[filename]=${fields[2]}
#     info[id]=${fields[3]}
#     info[name]=${fields[4]}
#     info[url]=${fields[5]}
#     info[urlPath]=${fields[6]}
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-install()
# #
# # Download and install a GitHub release asset
# #
# github-release-install ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# >= 4 && $# <= 5 )) || { printf 'Usage: github-install-latest-release $org $repo $releaseName $installFolder [$downloadFolder]\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local name=$3
#     local installFolder=$4
# 	local downloadFolder=${5:-$(create-temp-folder)}
# 	[[ -d "$downloadFolder" ]] || { echo "Non-existent folder: $downloadFolder" >&2; return 1; }
#     local -a flags=()
#     github-create-flags argmap flags token version
#     local releasePath; releasePath=$(github-release-download "${flags[@]}" "$org" "$repo" "$name" "$downloadFolder") || return
#     case "$releasePath" in
#         *.tgz|*.tar.gz)
#             tar --strip-components=1 -C "$installFolder" -xzf "$releasePath";;
#         *)
#             printf "Unsupported file type - can't install (%s)\n" "$releasePath" >&2
#             return 1;;
#     esac
# 	printf '%s' "$installFolder"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-install-latest()
# #
# # Install the latest release from a GitHub repo
# #
# # @note - github-release-install will install the latest by default, if you don't specify a version
# #
# github-release-install-latest ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# >= 4 && $# <= 5 )) || { printf 'Usage: github-release-install-latest $org $repo $releaseName $installFolder [$downloadFolder]\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local name=$3
#     local installFolder=$4
# 	local downloadFolder=${5:-''}
# 
#     local -a flags
#     github-create-flags argmap flags token
#     local version; version=$(github-release-get-latest-tag "${flags[@]}" "$org" "$repo") || return    
#     flags+=(--version "$version")
#     github-release-install "${flags[@]}" "$org" "$repo" "$releaseName" "$installFolder" "$downloadFolder"
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-list()
# #
# # List releases for a GitHub repository
# #
# github-release-list ()
# {
# 	# parse github args
# 	local -A argmap=()
# 	local nargs=0
# 	github-parse-args argmap nargs "$@"
# 	shift "$nargs"
# 	# shellcheck disable=SC2016
# 	(( $# == 2 )) || { printf 'Usage: github-release-list [flags] $org $repo\n' >&2; return 1; }
# 	local org=$1
# 	local repo=$2
# 
# 	# get release name list, using token if provided
#     local -a flags=()
# 	[[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
# 	github-release-get-data "${flags[@]}" "$org" "$repo" \
# 	| jq -r '[.assets[].name] | sort | @tsv' \
# 	|| return
# 
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-list-platforms()
# #
# # List available platforms for a release
# #
# github-release-list-platforms ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
# 	# shellcheck disable=SC2016
# 	(( $# = 2 )) || { printf 'Usage: github-release-list [flags] $org $repo\n' >&2; return 1; }
# 	local org=$1
# 	local repo=$2
# 
# 	# get release name list, using token if provided
#     readarray -t -d $'\t' releases < <(github-release-list "$@")
#     local platform
#     for release in "${releases[@]}"; do
#         if [[ ! "$release" =~ checksums.txt ]]; then
#             platform="${release##"${repo}"_}"
#             platform="${platform%%.*}"
#         	printf '%s\n' "$platform"
#         fi
#     done
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-select()
# #
# # Select a release from a list of choices
# #
# github-release-select ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
# 	# shellcheck disable=SC2016
# 	(( $# == 3 )) || { printf 'Usage: github-release-select [flags] name $org $repo\n' >&2; return 1; }
# 	[[ $1 != 'name' ]] && local -n name=$1
# 	local org=$2
# 	local repo=$3
# 
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
# 	IFS=$'\t' read -r -a names < <(github-release-list "${flags[@]}" "$org" "$repo") || return
# 	select name in "${names[@]}"; do break; done
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-release-select-platform()
# #
# # Select a platform for a release download
# #
# github-release-select-platform ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
# 	# shellcheck disable=SC2016
# 	(( $# = 2 )) || { printf 'Usage: github-release-select-platforms [flags] $org $repo' >&2; return 1; }
#     local platforms
#     readarray -t -d $'\n' platforms < <(github-release-list-platforms "$@")
# 	select platform in "${platforms[@]}"; do
#         printf '%s' "$platform" || return
#         break
#     done
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-test-repo()
# #
# # Test if a GitHub repo exists
# #
# # Simple attempt to get info for a repo
# # If it does not succeed, it could mean the org or repo are nonexistent or misspelled
# # But it could also mean that the repo is non-public and requires a token for authentication
# # The Github API returns 404s for all of the above, so the error status doesn't tell us anything
# #
# github-test-repo ()
# {
#     # parse github args
#     local -A argmap=()
#     local nargs=0
#     github-parse-args argmap nargs "$@" || return
#     shift "$nargs"
#     # shellcheck disable=SC2016
#     (( $# == 2 )) || { printf 'Usage: github-test-repo $org $repo\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
# 
#     local urlPath="/repos/$org/$repo"
#     local -a flags=()
#     [[ -v argmap[token] ]] && flags+=(--token "${argmap[token]}")
#     # We don't care about the info, just if we can successfully call the endpoint
#     github-curl "${flags[@]}" --output /dev/null "$urlPath" || return
# }
# 
# 
# #-------------------------------------------------------------------------------
# #
# # github-test-repo-with-auth()
# #
# # @deprecated - use github-test-repo and pass a token
# # Test if a GitHub repo exists with authentication
# #
# # Simple attempt to get info for a repo
# # If it does not succeed, it could mean the org or repo are nonexistent or misspelled
# # But it could also mean that the repo is non-public and requires a token for authentication
# # The Github API returns 404s for all of the above, so the error status doesn't tell us anything
# #
# github-test-repo-with-auth ()
# {
#     # shellcheck disable=SC2016
#     (( $# == 3 )) || { printf 'Usage: github-test-repo-with-auth $org $repo $token\n' >&2; return 1; }
#     local org=$1
#     local repo=$2
#     local token=$3
# 
#     local urlPath="/repos/$org/$repo"
#     # We don't care about the info, just if we can successfully call the endpoint
#     github-curl --output /dev/null --token "$token" "$urlPath" || return
# }


#------------------------------------------------------------------------------
#
# install-build-dylt-svc ()
#
install-build-dylt-svc ()
{
    local repo=https://raw.githubusercontent.com/dylt-dev/dylt/main
    mkdir -p /opt/svc/build-dylt/bin
    chown -R rayray:rayray /opt/svc/build-dylt
    curl --silent --remote-name --output-dir /opt/svc/build-dylt "$repo/svc/build-dylt/build-dylt.service"
    curl --silent --remote-name --output-dir /opt/svc/build-dylt "$repo/svc/build-dylt/build-dylt.timer"
    curl --silent --remote-name --output-dir /opt/svc/build-dylt/bin "$repo/svc/build-dylt/bin/run.sh"
    chmod 777 /opt/svc/build-dylt/bin/run.sh
    systemctl enable /opt/svc/build-dylt/build-dylt.service
    systemctl enable /opt/svc/build-dylt/build-dylt.timer
    systemctl start build-dylt.timer
}


#-------------------------------------------------------------------------------
#
# safe-curl()
#
# Download a URL silently with curl, capturing errors to a temp file
# On success, the response body is written to stdout with no other output
# On failure, the temp file path is printed in the error message
# Additional arguments are passed through to curl
#
safe-curl ()
{
    (( $# >= 1 )) || { printf 'Usage: safe-curl $url [$curlArgs]\n' >&2; return 1; }
    local url=$1
    shift
    local errfile
    errfile=$(mktemp 2>/dev/null) || errfile=/tmp/safe-curl.err
    curl --fail --location --silent --show-error "$@" "$url" 2>"$errfile" || {
        printf 'error: curl failed for %s\n' "$url" >&2
        printf '  details saved to %s\n' "$errfile" >&2
        return 1
    }
    rm -f "$errfile"
}


#-------------------------------------------------------------------------------
#
# source-github-utils()
#
# Ensure github-utils.sh is available by sourcing it idempotently
# If github-curl is already defined, do nothing
# Otherwise, download github-utils.sh from daylight and source it
# Returns 1 if download or source fails, or if expected functions are missing
#
source-github-utils ()
{
    (( $# == 0 )) || { printf 'Usage: source-github-utils\n' >&2; return 1; }
    type github-curl &>/dev/null && return 0
    source-remote-script "https://raw.githubusercontent.com/daylight-public/daylight/main/github-utils.sh" || return
    type github-curl &>/dev/null || {
        printf 'error: sourced script does not define github-curl\n' >&2
        return 1
    }
}


#-------------------------------------------------------------------------------
#
# source-remote-script()
#
# Download a script from a URL to a temp file and source it
# Returns 1 if curl fails, the temp file cannot be created, or source fails
# After sourcing, the temp file is removed
#
source-remote-script ()
{
    (( $# == 1 )) || { printf 'Usage: source-remote-script $url\n' >&2; return 1; }
    local url=$1 tmp
    tmp=$(mktemp) || { printf 'error: could not create temp file\n' >&2; return 1; }
    safe-curl "$url" > "$tmp" || { rm -f "$tmp"; return 1; }
    source "$tmp" || {
        printf 'error: could not source %s\n' "$url" >&2
        rm -f "$tmp"
        return 1
    }
    rm -f "$tmp"
}


#------------------------------------------------------------------------------
#
# yesorno()
#
# From @day-sh/app-funcs.sh
#
yesorno ()
{
    # shellcheck disable=SC2016
    (( $# == 2 )) || { printf 'Usage: yesorno varname $prompt\n' >&2; return 1; }
    [[ $1 != val ]] && local -n val=$1
	val=''
    local prompt=$2

    local s
    while [[ ! "${s,,}" =~ y|n ]]; do
        read -r -n1 -p "$prompt" s || return
        echo
        if [[ $s == '?' ]]; then
            printf 'Please enter Y, y, N, or n\n'
        fi
    done
    # shellcheck disable=SC2034
    val=$s
	case "$val" in
		Y|y) return 0 ;;
		N|n) return 1 ;;
		*) printf '*** Something is wrong ($val=%s)'  "$val"; return 1 ;;
	esac
}


#------------------------------------------------------------------------------
#
# main()
#
# Dispatch command lines arguments to the appropriate function
#
main ()
{
    if (($# >= 1)); then
        local cmd=$1
        shift
        case "$cmd" in
            add-to-bashrc)                            add-to-bashrc "$@";;
            gen-nightly-tagname)                      gen-nightly-tagname "$@";;
            gen-nightly-timestamp)                    gen-nightly-timestamp "$@";;
            git-do-nightly-release)                   git-do-nightly-release "$@";;
            git-download-latest-daylightsh)           git-download-latest-daylightsh "$@";;
            git-get-latest-release-spec)              git-get-latest-release-spec "$@";;
            git-get-latest-release-tag)               git-get-latest-release-tag "$@";;
            git-get-latest-release-version)           git-get-latest-release-version "$@";;
            git-install-latest-daylightsh)            git-install-latest-daylightsh "$@";;
            git-install-latest-dylt)                  git-install-latest-dylt "$@";;
            git-tag-nightly)                          git-tag-nightly "$@";;
            github-app-get-client-id)                 github-app-get-client-id "$@";;
            github-app-get-data)                      github-app-get-data "$@";;
            github-app-get-id)                        github-app-get-id "$@";;
            github-app-get-info)                      github-app-get-info "$@";;
            github-create-flags)                      github-create-flags "$@";;
            github-create-url)                        github-create-url "$@";;
            github-create-user-access-token)          github-create-user-access-token "$@";;
            github-curl)                              github-curl "$@";;
            github-curl-post)                         github-curl-post "$@";;
            github-download-latest-release)           github-download-latest-release "$@";;
            github-get-release-data)                  github-get-release-data "$@";;
            github-get-release-name-list)             github-get-release-name-list "$@";;
            github-get-release-package-data)          github-get-release-package-data "$@";;
            github-get-release-package-info)          github-get-release-package-info "$@";;
            github-parse-args)                        github-parse-args "$@";;
            github-release-create-url-path)           github-release-create-url-path "$@";;
            github-release-download)                  github-release-download "$@";;
            github-release-download-latest)           github-release-download-latest "$@";;
            github-release-get-data)                  github-release-get-data "$@";;
            github-release-get-latest-tag)            github-release-get-latest-tag "$@";;
            github-release-get-package-data)          github-release-get-package-data "$@";;
            github-release-get-package-info)          github-release-get-package-info "$@";;
            github-release-install)                   github-release-install "$@";;
            github-release-install-latest)            github-release-install-latest "$@";;
            github-release-list)                      github-release-list "$@";;
            github-release-list-platforms)            github-release-list-platforms "$@";;
            github-release-select)                    github-release-select "$@";;
            github-release-select-platform)           github-release-select-platform "$@";;
            github-test-repo)                         github-test-repo "$@";;
            github-test-repo-with-auth)               github-test-repo-with-auth "$@";;
	    install-build-dylt-svc)                   install-build-dylt-svc "$@";;
            yesorno)                                  yesorno "$@";;
            *) printf 'Unknown command: %s \n' "$cmd";;
        esac
    fi
}


if ! (return 0 2>/dev/null); then
    main "$@"
fi
