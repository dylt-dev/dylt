#! /usr/bin/env bash



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