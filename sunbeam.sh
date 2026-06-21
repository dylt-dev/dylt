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
# download-dylt-batch()
#
# Download dylt from a GitHub branch or release
#
download-dylt-batch ()
{
    local branch="" release="" latest=0
    local -a pass=()
    local dstFolder=""
    local platform=""
    local extract_flag="" extract_dir="" extract_name=""

    local args=("$@")
    local i=0
    while (( i < $# )); do
        case "${args[i]}" in
            --branch)
                if [[ -n "$release" ]]; then
                    printf 'Error: --branch and --release are incompatible\n' >&2
                    return 1
                fi
                if (( i+1 < $# )) && [[ "${args[i+1]}" != --* ]]; then
                    branch=${args[i+1]}
                    (( i++ ))
                else
                    branch=main
                fi
                ;;
            --release)
                if [[ -n "$branch" ]]; then
                    printf 'Error: --branch and --release are incompatible\n' >&2
                    return 1
                fi
                if (( i+1 < $# )) && [[ "${args[i+1]}" != --* ]]; then
                    release=${args[i+1]}
                    (( i++ ))
                else
                    release=latest
                fi
                ;;
            --latest)
                latest=1
                ;;
            --token)
                if (( i+1 < $# )); then
                    pass+=("${args[i]}" "${args[i+1]}")
                    (( i++ ))
                else
                    printf 'Error: --token requires a value\n' >&2
                    return 1
                fi
                ;;
            --output-dir)
                if (( i+1 < $# )); then
                    pass+=("${args[i]}" "${args[i+1]}")
                    (( i++ ))
                else
                    printf 'Error: --output-dir requires a value\n' >&2
                    return 1
                fi
                ;;
            --remote-name)
                pass+=("${args[i]}")
                ;;
            --extract)
                extract_flag=1
                ;;
            --extract-dir)
                if (( i+1 < $# )); then
                    extract_dir=${args[i+1]}
                    (( i++ ))
                else
                    printf 'Error: --extract-dir requires a value\n' >&2
                    return 1
                fi
                ;;
            --extract-name)
                if (( i+1 < $# )); then
                    extract_name=${args[i+1]}
                    (( i++ ))
                else
                    printf 'Error: --extract-name requires a value\n' >&2
                    return 1
                fi
                ;;
            --)
                (( i++ ))
                break
                ;;
            --*)
                printf 'Unknown flag: %s\n' "${args[i]}" >&2
                return 1
                ;;
            *)
                if [[ -z "$dstFolder" ]]; then
                    dstFolder=${args[i]}
                else
                    platform=${args[i]}
                fi
                ;;
        esac
        (( i++ ))
    done

    [[ -n "$dstFolder" ]] || { printf 'Usage: download-dylt-batch [--branch [<name>]] [--release [<tag>]] [--latest] [--token <pat>] [--output-dir <dir>] [--remote-name] [--extract] [--extract-dir <dir>] [--extract-name <name>] [--] <dstFolder> [<platform>]\n' >&2; return 1; }
    [[ -d "$dstFolder" ]] || { printf 'Non-existent folder: %s\n' "$dstFolder" >&2; return 1; }

    if (( latest )) && [[ -z "$release" ]]; then
        printf 'Error: --latest requires --release\n' >&2
        return 1
    fi
    if [[ -z "$branch" && -z "$release" ]]; then
        branch=main
    fi

    [[ -z "$platform" ]] && platform=$(detect-platform) || :
    [[ -z "$extract_dir" ]] && extract_dir=$dstFolder

    if [[ -n "$release" ]]; then
        source-github-utils || return

        local tag json assetName checksumName tmpDir releasePath checksumPath

        if [[ "$release" == "latest" ]]; then
            tag=$(github-release-get-latest-tag "${pass[@]}" dylt-dev dylt) || return
        else
            tag=$release
        fi

        tmpDir=$(mktemp -d) || return

        json=$(github-release-get-data --version "$tag" "${pass[@]}" dylt-dev dylt) || {
            rm -rf "$tmpDir"
            return 1
        }

        assetName=$(printf '%s' "$json" | jq -r '.assets[] | select(.name | endswith(".tar.gz") and startswith("dylt_'"$platform"'")) | .name' | head -1) || :
        if [[ -z "$assetName" ]]; then
            printf 'No tar.gz asset found for platform %s in release %s\n' "$platform" "$tag" >&2
            rm -rf "$tmpDir"
            return 1
        fi

        checksumName=$(printf '%s' "$json" | jq -r '.assets[] | select(.name | endswith("checksums.txt")) | .name' | head -1) || :

        local release_args=(--version "$tag")
        if [[ -n "$extract_flag" || -n "$extract_dir" || -n "$extract_name" ]]; then
            release_args+=(--extract-dir "$extract_dir" --extract-name "${extract_name:-dylt}")
        fi
        releasePath=$(github-release-download "${release_args[@]}" "${pass[@]}" dylt-dev dylt "$assetName" "$tmpDir") || {
            rm -rf "$tmpDir"
            return 1
        }

        if [[ -n "$checksumName" ]]; then
            checksumPath=$(github-release-download --version "$tag" "${pass[@]}" dylt-dev dylt "$checksumName" "$tmpDir") || {
                printf 'Warning: could not download %s — skipping checksum verification\n' "$checksumName" >&2
                checksumPath=""
            }
        fi

        if [[ -n "$checksumPath" ]]; then
            if ! (cd "$tmpDir" && grep -F "$assetName" "$checksumName" | sha256sum -c -); then
                printf 'Checksum verification failed for %s\n' "$assetName" >&2
                rm -rf "$tmpDir"
                return 1
            fi
        fi

        mv "$releasePath" "$dstFolder/${extract_name:-$(basename "$releasePath")}" || {
            rm -rf "$tmpDir"
            return 1
        }
        rm -rf "$tmpDir"
    else
        local url="https://raw.githubusercontent.com/dylt-dev/dylt/$branch/sunbeam.sh"
        local target="$extract_dir/${extract_name:-sunbeam.sh}"
        mkdir -p "$extract_dir"
        curl --location --silent --fail --output "$target" "$url" || return
    fi
}


#-------------------------------------------------------------------------------
#
# download-dylt()
#
# Download dylt with optional interactive prompts
#
download-dylt ()
{
    download-dylt-batch "$@" || return
}


#-------------------------------------------------------------------------------
#
# download-daylight()
#
# Download daylight.sh with optional interactive prompts
#
download-daylight ()
{
    local gen_completions=""
    local completions_path=""
    local -a pass_args=()
    local args=("$@")

    local i=0
    while (( i < $# )); do
        case "${args[i]}" in
            --gen-bash-completions)
                gen_completions=1
                if (( i+1 < $# )) && [[ "${args[i+1]}" != --* ]]; then
                    completions_path=${args[i+1]}
                    (( i++ ))
                fi
                ;;
            *)
                pass_args+=("${args[i]}")
                ;;
        esac
        (( i++ ))
    done

    download-daylight-batch "${pass_args[@]}" || return

    if [[ -z "$gen_completions" ]] && [[ -t 0 ]]; then
        printf 'Generate bash completions for daylight.sh? [y/N] '
        local reply
        read -r reply
        [[ "$reply" =~ ^[yY] ]] || return 0
    fi

    # Find dstFolder from pass_args (last positional arg)
    local dstFolder="${pass_args[${#pass_args[@]}-1]}"

    [[ -n "$dstFolder" ]] || { printf 'error: could not determine destination folder\n' >&2; return 1; }

    local scriptPath="$dstFolder/daylight.sh"
    [[ -f "$scriptPath" ]] || { printf 'error: %s not found after download\n' "$scriptPath" >&2; return 1; }

    local compPath="${completions_path:-$HOME/bash-completion.d/daylight.sh}"
    mkdir -p "$(dirname "$compPath")" || return

    bash "$scriptPath" gen-completion-script daylight.sh < <(bash "$scriptPath" list-bash-funcs < "$scriptPath") > "$compPath" || return
    printf 'Bash completions written to %s\n' "$compPath" >&2
}


#-------------------------------------------------------------------------------
#
# download-daylight-batch()
#
# Download daylight.sh from a GitHub branch or release
#
download-daylight-batch ()
{
    local branch="" release="" latest=0
    local -a pass=()
    local dstFolder=""
    local extract_flag="" extract_dir="" extract_name=""

    local args=("$@")
    local i=0
    while (( i < $# )); do
        case "${args[i]}" in
            --branch)
                if [[ -n "$release" ]]; then
                    printf 'Error: --branch and --release are incompatible\n' >&2
                    return 1
                fi
                if (( i+1 < $# )) && [[ "${args[i+1]}" != --* ]]; then
                    branch=${args[i+1]}
                    (( i++ ))
                else
                    branch=main
                fi
                ;;
            --release)
                if [[ -n "$branch" ]]; then
                    printf 'Error: --branch and --release are incompatible\n' >&2
                    return 1
                fi
                if (( i+1 < $# )) && [[ "${args[i+1]}" != --* ]]; then
                    release=${args[i+1]}
                    (( i++ ))
                else
                    release=latest
                fi
                ;;
            --latest)
                latest=1
                ;;
            --token)
                if (( i+1 < $# )); then
                    pass+=("${args[i]}" "${args[i+1]}")
                    (( i++ ))
                else
                    printf 'Error: --token requires a value\n' >&2
                    return 1
                fi
                ;;
            --output-dir)
                if (( i+1 < $# )); then
                    pass+=("${args[i]}" "${args[i+1]}")
                    (( i++ ))
                else
                    printf 'Error: --output-dir requires a value\n' >&2
                    return 1
                fi
                ;;
            --remote-name)
                pass+=("${args[i]}")
                ;;
            --extract)
                extract_flag=1
                ;;
            --extract-dir)
                if (( i+1 < $# )); then
                    extract_dir=${args[i+1]}
                    (( i++ ))
                else
                    printf 'Error: --extract-dir requires a value\n' >&2
                    return 1
                fi
                ;;
            --extract-name)
                if (( i+1 < $# )); then
                    extract_name=${args[i+1]}
                    (( i++ ))
                else
                    printf 'Error: --extract-name requires a value\n' >&2
                    return 1
                fi
                ;;
            --)
                (( i++ ))
                break
                ;;
            --*)
                printf 'Unknown flag: %s\n' "${args[i]}" >&2
                return 1
                ;;
            *)
                if [[ -z "$dstFolder" ]]; then
                    dstFolder=${args[i]}
                else
                    printf 'Unexpected argument: %s\n' "${args[i]}" >&2
                    return 1
                fi
                ;;
        esac
        (( i++ ))
    done

    [[ -n "$dstFolder" ]] || { printf 'Usage: download-daylight-batch [--branch [<name>]] [--release [<tag>]] [--latest] [--token <pat>] [--output-dir <dir>] [--remote-name] [--extract] [--extract-dir <dir>] [--extract-name <name>] [--] <dstFolder>\n' >&2; return 1; }
    [[ -d "$dstFolder" ]] || { printf 'Non-existent folder: %s\n' "$dstFolder" >&2; return 1; }

    if (( latest )) && [[ -z "$release" ]]; then
        printf 'Error: --latest requires --release\n' >&2
        return 1
    fi
    if [[ -z "$branch" && -z "$release" ]]; then
        branch=main
    fi

    [[ -z "$extract_dir" ]] && extract_dir=$dstFolder

    if [[ -n "$release" ]]; then
        source-github-utils || return

        local tag json assetName tmpDir releasePath checksumFile
        local checksumName=SHA256SUMS

        if [[ "$release" == "latest" ]]; then
            tag=$(github-release-get-latest-tag "${pass[@]}" daylight-public daylight) || return
        else
            tag=$release
        fi

        tmpDir=$(mktemp -d) || return

        json=$(github-release-get-data --version "$tag" "${pass[@]}" daylight-public daylight) || {
            rm -rf "$tmpDir"
            return 1
        }

        assetName=$(printf '%s' "$json" | jq -r '.assets[] | select(.name | endswith(".tar.gz")) | .name' | head -1) || :
        if [[ -z "$assetName" ]]; then
            printf 'No tar.gz asset found in release %s\n' "$tag" >&2
            rm -rf "$tmpDir"
            return 1
        fi

        releasePath=$(github-release-download --version "$tag" "${pass[@]}" --extract-dir "$extract_dir" --extract-name "${extract_name:-daylight.sh}" daylight-public daylight "$assetName" "$tmpDir") || {
            rm -rf "$tmpDir"
            return 1
        }

        checksumFile=$(github-release-download --version "$tag" "${pass[@]}" daylight-public daylight "$checksumName" "$tmpDir") || {
            printf 'SHA256SUMS not found in release %s — cannot verify integrity\n' "$tag" >&2
            rm -rf "$tmpDir"
            return 1
        }

        if ! (cd "$tmpDir" && grep -F "$assetName" "$checksumName" | sha256sum -c -); then
            printf 'Checksum verification failed for %s\n' "$assetName" >&2
            rm -rf "$tmpDir"
            return 1
        fi

        # releasePath now points to the extracted file from github-release-download
        rm -rf "$tmpDir"
    else
        local url="https://raw.githubusercontent.com/daylight-public/daylight/$branch/daylight.sh"
        local target="$extract_dir/${extract_name:-daylight.sh}"
        mkdir -p "$extract_dir"
        curl --location --silent --fail --output "$target" "$url" || return
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

    printf 'NOTE: Create %s/%s/env with GITHUB_TOKEN=... before timer fires\n' "$svcDir" "$instance"
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
# trigger-nightly-release-batch()
#
# Trigger a GHA workflow for a given repo via workflow_dispatch.
# Requires --workflow; accepts --token (falls back to GITHUB_TOKEN env var)
# and --label. No interactivity, no token inference.
#
trigger-nightly-release-batch ()
{
    source-github-utils || return

    local -A argmap=()
    local nargs=0
    github-curl-parse-args argmap nargs "$@" || return
    shift "$nargs"
    # shellcheck disable=SC2016
    (( $# == 1 )) || { printf 'Usage: trigger-nightly-release-batch --workflow <name> [--token <pat>] [--label <label>] $owner/$repo\n' >&2; return 1; }
    local repo=$1

    local workflow=${argmap[workflow]}
    [[ -n "$workflow" ]] || { printf 'error: --workflow is required\n' >&2; return 1; }

    local token=${argmap[token]:-${GITHUB_TOKEN:?error: --token not given and GITHUB_TOKEN not set}}

    local wf_name=${workflow%.yml}
    wf_name=${wf_name%.yaml}
    if ! curl -sf -o /dev/null \
        "https://api.github.com/repos/$repo/actions/workflows/${wf_name}.yml"; then
        printf 'error: workflow "%s" not found in %s\n' "$workflow" "$repo" >&2
        return 1
    fi

    source-github-utils || return

    local -a flags=(--token "$token")
    local data
    if [[ -n "${argmap[label]}" ]]; then
        local label
        label=$(sanitize-label "${argmap[label]}") || return
        data=$(printf '{"ref":"main","inputs":{"label":"%s"}}' "$label")
    else
        data='{"ref":"main"}'
    fi
    flags+=(--data "$data")

    github-curl "${flags[@]}" "/repos/$repo/actions/workflows/${wf_name}.yml/dispatches" || return
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
    local -A argmap=()
    local nargs=0
    github-curl-parse-args argmap nargs "$@" || return
    shift "$nargs"
    # shellcheck disable=SC2016
    (( $# == 1 )) || { printf 'Usage: trigger-nightly-release [--workflow <name>] [--token <pat>] [--label <label>] $owner/$repo\n' >&2; return 1; }
    local repo=$1

    local workflow=${argmap[workflow]:-nightly-release}
    local label=${argmap[label]:-''}

    local token
    if [[ -v argmap[token] ]]; then
        token=${argmap[token]}
    elif [[ -n "${GITHUB_TOKEN-}" ]]; then
        token=$GITHUB_TOKEN
    elif [[ -n "${GH_TOKEN-}" ]]; then
        token=$GH_TOKEN
    elif type gh &>/dev/null; then
        token=$(gh auth token 2>/dev/null) || token=''
    fi

    local -a batch_args=(--workflow "$workflow" --token "$token")
    [[ -n "$label" ]] && batch_args+=(--label "$label")
    trigger-nightly-release-batch "${batch_args[@]}" "$repo" || return
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
            download-dylt)                            download-dylt "$@";;
            download-dylt-batch)                      download-dylt-batch "$@";;
            download-daylight)                        download-daylight "$@";;
            download-daylight-batch)                  download-daylight-batch "$@";;
            git-download-latest-daylightsh)           git-download-latest-daylightsh "$@";;
            git-get-latest-release-spec)              git-get-latest-release-spec "$@";;
            git-get-latest-release-tag)               git-get-latest-release-tag "$@";;
            git-get-latest-release-version)           git-get-latest-release-version "$@";;
            git-install-latest-daylightsh)            git-install-latest-daylightsh "$@";;
            git-install-latest-dylt)                  git-install-latest-dylt "$@";;
	    install-build-dylt-svc)                   install-build-dylt-svc "$@";;
            sanitize-label)                           sanitize-label "$@";;
            trigger-nightly-release)                  trigger-nightly-release "$@";;
            trigger-nightly-release-batch)            trigger-nightly-release-batch "$@";;
            yesorno)                                  yesorno "$@";;
            *) printf 'Unknown command: %s \n' "$cmd";;
        esac
    fi
}


if ! (return 0 2>/dev/null); then
    main "$@"
fi
