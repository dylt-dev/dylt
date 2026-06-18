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
# Install the nightly release service for dylt using the indexed template from
# daylight-public/daylight.  This replaces the old per-repo service layout with
# the shared template pattern at /opt/svc/nightly-release/<instance>/.
#
install-build-dylt-svc ()
{
    local daylightBase="https://raw.githubusercontent.com/daylight-public/daylight/main"
    local instance=dylt-dev-dylt
    local svcDir=/opt/svc/nightly-release
    local repoDir=$svcDir/$instance/repo

    mkdir -p "$svcDir/$instance/bin"
    chown -R rayray:rayray "$svcDir"

    curl -s --remote-name --output-dir /etc/systemd/system \
        "$daylightBase/svc/nightly-release/nightly-release@.service"
    curl -s --remote-name --output-dir /etc/systemd/system \
        "$daylightBase/svc/nightly-release/nightly-release@.timer"

    cat > "$svcDir/$instance/bin/run.sh" <<'RUNEOF'
#!/usr/bin/env bash
main ()
{
    local svcDir=/opt/svc/nightly-release
    cd "$svcDir/dylt-dev-dylt/repo" || exit 1
    git pull --ff-only origin main || exit 1
    source ./sunbeam.sh || exit 1
    trigger-nightly-release "dylt-dev/dylt" || exit 1
}
main "$@"
RUNEOF
    chmod 755 "$svcDir/$instance/bin/run.sh"

    systemctl enable "nightly-release@$instance.service"
    systemctl enable "nightly-release@$instance.timer"
    systemctl start "nightly-release@$instance.timer"

    printf 'NOTE: Create %s/%s/env with GITHUB_PAT=... before timer fires\n' "$svcDir" "$instance"
    printf 'NOTE: Clone the repo with: git clone git@github.com:dylt-dev/dylt.git %s\n' "$repoDir"
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
# sanitize-label()
#
# Sanitize a human-readable label for use in a git tag.
# Spaces become dashes; non-alphanumeric, non-dot, non-underscore,
# non-hyphen characters are stripped.
#
sanitize-label ()
{
    (( $# == 1 )) || { printf 'Usage: sanitize-label $label\n' >&2; return 1; }
    local label=$1
    label=${label// /-}
    label=${label//[^a-zA-Z0-9._-]/}
    printf '%s' "$label"
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
# trigger-nightly-release()
#
# Trigger a nightly-release GHA workflow for a given repo via workflow_dispatch.
# Requires $GITHUB_PAT in the environment.
#
trigger-nightly-release ()
{
    (( $# == 1 )) || { printf 'Usage: trigger-nightly-release $owner/$repo\n' >&2; return 1; }
    local repo=$1
    local token=${GITHUB_PAT:?error: GITHUB_PAT not set}
    curl --fail --silent --show-error -X POST \
        "https://api.github.com/repos/$repo/actions/workflows/nightly-release.yml/dispatches" \
        -H "Authorization: Bearer $token" \
        -H "Accept: application/vnd.github.v3+json" \
        -d '{"ref": "main"}' || return
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
            git-download-latest-daylightsh)           git-download-latest-daylightsh "$@";;
            git-get-latest-release-spec)              git-get-latest-release-spec "$@";;
            git-get-latest-release-tag)               git-get-latest-release-tag "$@";;
            git-get-latest-release-version)           git-get-latest-release-version "$@";;
            git-install-latest-daylightsh)            git-install-latest-daylightsh "$@";;
            git-install-latest-dylt)                  git-install-latest-dylt "$@";;
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
            sanitize-label)                           sanitize-label "$@";;
            trigger-nightly-release)                  trigger-nightly-release "$@";;
            yesorno)                                  yesorno "$@";;
            *) printf 'Unknown command: %s \n' "$cmd";;
        esac
    fi
}


if ! (return 0 2>/dev/null); then
    main "$@"
fi
