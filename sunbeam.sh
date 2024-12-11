#! /usr/bin/env bash


# Download the latest daylight.sh from github
#
# By default, download to ~/tmp/
git-download-latest-daylightsh ()
{
	# shellcheck disable=SC2016
	{ (( $# >= 0 )) && (( $# <= 1 )); } || { printf 'Usage: git-download-latest-daylightsh [$downloadFolder]\n' >&2; return 1; }
	local downloadFolder=${1:-"$HOME/tmp"}
	[[ -d "$downloadFolder" ]] || { echo "Non-existent folder: $downloadFolder" >&2; return 1; }

	url=https://raw.githubusercontent.com/daylight-public/daylight/main/daylight.sh
	curl --silent --remote-name --output-dir "/$downloadFolder/" "$url"
}

# Create a tagname from a root version (eg v1.0.7) and the nightly timestamp
# eg v1.0.7-nightly.20240522152147
gen-nightly-tagname ()
{
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: gen-nightly-tagname $version\n' >&2; return 1; }
	local version=$1

	local ts; ts=$(gen-nightly-timestamp) || return
	local label=$(printf '%s-nightly.%s' "$version" "$ts")

	printf "$label"
}


# Simple function to create a timestamp in the format used for nightly release tags
# eg 20240522152147 (Wed May 22 15:21:47 CDT 2024)
# More separators for the various date fields would be great, but they would break https://semver.org rules
gen-nightly-timestamp ()
{
	date '+%Y%m%d%H%M%S'
}


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


# Use the GitHub API to get the tag of the latest version of a GitHub release,
# for a specified owner+repo
# This can be used to `go install` a specific releae
# Note that GitHub defines 'latest version' as the release that was created most recently,
# unlike https://semver.org, which has complicated rules to define the most recent release.
git-get-latest-release-tag ()
{
	# shellcheck disable=SC2016
	(( $# == 2 )) || { printf 'Usage: git-get-latest-release-tag $owner $repo\n' >&2; return 1; }
	local owner=$1
	local repo=$2
	local tag; tag=$(curl -L --silent "api.github.com/repos/$owner/$repo/releases/latest" | jq -r .tag_name)
	printf '%s' "$tag"
}


# Install the latest daylight.sh on a VM
#
# By default, install from github
# If $scriptPath is specified, install that one instead
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



git-install-latest-dylt ()
{
	local owner=dylt-dev
	local repo=dylt
	local tag=$(git-get-latest-release-tag "$owner" "$repo")
	local release=github.com/$owner/$repo@$tag
	go install "$release"
}	


# Tag the nightly release, push the current commit, and push the tag
git-push-nightly-release ()
{
	# @todo Use [[ $(git status --porcelain) == "" ]] to see if there is uncommited work. If so ask for confirmation
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: git-tag-nightly $version\n' >&2; return 1; }
	local version=$1

	if [[ $(git status --porcelain) != "" ]]; then
		printf '%s\n' "There are uncommitted changes"
		if ! yesno "Push nightly release anyway?"; then 
			return 0
		fi
	fi
	git-tag-nightly "$version" || return
	git push
	git push --tags
}


# Create the nightly tagname, and then a git tag from the name
git-tag-nightly ()
{
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: git-tag-nightly $version\n' >&2; return 1; }
	local version=$1

	local tagname; tagname=$(gen-nightly-tagname $version) || return
	git tag "$tagname"
}


yesno ()
{
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: yesno $prompt\n' >&2; return 1; }
	read -p "$1 " -n 1 -r yn
	if [[ ! $yn =~ ^[Yy]$ ]]; then
		return 1
	fi
}
