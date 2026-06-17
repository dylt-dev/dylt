#!/usr/bin/env bash
# install-sunbeam.sh — Bootstrap sunbeam.sh
#
# sunbeam.sh is a utility script for dylt development tasks.
# It lives at ~/.local/bin/sunbeam.sh and provides commands for
# nightly releases, tagging, and other dev workflows.
#
# ## Usage
#
#   curl -fsSL https://raw.githubusercontent.com/dylt-dev/dylt/main/install-sunbeam.sh | bash
#
# ## Ubuntu / Debian
#
#   curl -fsSL https://raw.githubusercontent.com/dylt-dev/dylt/main/install-sunbeam.sh | bash
#
# ## macOS
#
#   curl -fsSL https://raw.githubusercontent.com/dylt-dev/dylt/main/install-sunbeam.sh | bash
#
# ## Windows (Git Bash / WSL)
#
#   curl -fsSL https://raw.githubusercontent.com/dylt-dev/dylt/main/install-sunbeam.sh | bash
#
# ## What it does
#
#   1. Downloads sunbeam.sh to ~/.local/bin/sunbeam.sh
#   2. Makes it executable
#
# ## After install
#
#   Run `sunbeam add-to-bashrc` to add a `sunbeam()` shell function
#   to ~/.bashrc so you can type `sunbeam` from anywhere.
#
#   Then restart your shell or run: source ~/.bashrc


main ()
{
    local installPath="${1:-$HOME/.local/bin/sunbeam.sh}"

    printf 'Installing sunbeam.sh to %s ...\n' "$installPath"
    mkdir -p "$(dirname "$installPath")" || { printf 'Failed to create directory\n' >&2; exit 1; }

    curl -fsSL -o "$installPath" https://raw.githubusercontent.com/dylt-dev/dylt/main/sunbeam.sh || {
        printf 'Failed to download sunbeam.sh\n' >&2
        exit 1
    }

    chmod 755 "$installPath" || { printf 'Failed to chmod\n' >&2; exit 1; }

    printf '\nDone — installed to %s\n' "$installPath"
    printf 'Next step: run  %s add-to-bashrc  to add the `sunbeam` shell function\n' "$installPath"
    printf 'Then:   source ~/.bashrc\n'
    printf 'Then:   sunbeam --help\n'
}


main "$@"
