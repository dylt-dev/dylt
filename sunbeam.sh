#! /usr/bin/env bash



gen-nightly-tagname ()
{
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: gen-nightly-tagname $version\n' >&2; return 1; }
	local version=$1

	local ts; ts=$(gen-nightly-timestamp) || return
	local label=$(printf '%s-nightly.%s' "$version" "$ts")

	printf "$label"
}


gen-nightly-timestamp ()
{
	date '+%Y%m%d%H%M%S'
}


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
	local version=github.com/$owner/$repo@$tag
	go install "$version"
}	


git-tag-nightly ()
{
	# shellcheck disable=SC2016
	(( $# == 1 )) || { printf 'Usage: git-tag-nightly $version\n' >&2; return 1; }
	local version=$1

	local tagname; tagname=$(gen-nightly-tagname $version) || return
	git tag "$tagname"
}
