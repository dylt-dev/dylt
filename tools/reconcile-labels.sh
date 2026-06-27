#! /usr/bin/env bash

#-------------------------------------------------------------------------------
#
# reconcile-labels.sh
#
# Reconcile labels on GitHub repos to the standard set: bug, feature, task.
# Renames enhancement -> feature, creates task if missing, deletes all others.
# Writes a snapshot of pre-change labels to each repo's misc/gh-labels.json.
#
# Usage:
#   reconcile-labels.sh [--all] [--dry-run] [owner/repo ...]
#
#   --all     auto-detect repos cloned under ~/src/
#   --dry-run show what would change without making API calls
#

set -euo pipefail

ALLOWLIST="bug feature task"
TASK_COLOR="c5def5"
TASK_DESC="An internal improvement or chore"

find_all=false
dry_run=false
repos=()

while [[ $# -gt 0 ]]; do
    case "$1" in
        --all)     find_all=true; shift ;;
        --dry-run) dry_run=true; shift ;;
        --*)       echo "Unknown flag: $1" >&2; exit 1 ;;
        *)         repos+=("$1"); shift ;;
    esac
done

if "$find_all"; then
    if (( ${#repos[@]} > 0 )); then
        echo "error: --all cannot be combined with positional repo args" >&2
        exit 1
    fi
    for d in "$HOME/src"/*/.git; do
        [[ -d "$d" ]] || continue
        remote=$(git -C "$(dirname "$d")" remote get-url origin 2>/dev/null) || continue
        remote=$(echo "$remote" | sed 's/.*github.com[:\/]//; s/\.git$//; s/.*@//')
        [[ -n "$remote" ]] && repos+=("$remote")
    done
fi

if (( ${#repos[@]} == 0 )); then
    echo "Usage: reconcile-labels.sh [--all] [--dry-run] owner/repo ..." >&2
    exit 1
fi

reconcile_one () {
    local repo=$1
    local repo_name=${repo#*/}
    local local_path=""
    local candidate

    for candidate in "$HOME/src/$repo_name" "$HOME/src/$repo"; do
        if [[ -d "$candidate/.git" ]]; then
            local_path=$candidate
            break
        fi
    done

    echo "$repo"

    # Fetch current labels
    local json
    json=$(gh label list --repo "$repo" --json name,color,description 2>/dev/null) || {
        echo "  error: cannot list labels (no access?)" >&2
        return 1
    }

    # Snapshot
    if [[ -n "$local_path" ]] && ! "$dry_run"; then
        mkdir -p "$local_path/misc"
        echo "$json" > "$local_path/misc/gh-labels.json"
        echo "  snapshot -> ${local_path}/misc/gh-labels.json"
    fi

    # Classify labels
    local has_enhancement=false has_feature=false has_task=false
    local -a to_delete=()

    while IFS= read -r name; do
        case "$name" in
            bug)         : ;;
            feature)     has_feature=true ;;
            task)        has_task=true ;;
            enhancement) has_enhancement=true ;;
            *)           to_delete+=("$name") ;;
        esac
    done < <(echo "$json" | jq -r '.[].name')

    # Rename enhancement -> feature
    if "$has_enhancement"; then
        if "$dry_run"; then
            echo "  would rename: enhancement -> feature"
        else
            gh label edit --repo "$repo" enhancement --name feature 2>/dev/null \
                && echo "  renamed: enhancement -> feature" \
                || echo "  FAILED: rename enhancement -> feature" >&2
        fi
        has_feature=true
    fi

    # Create feature if neither exists
    if ! "$has_feature" && ! "$has_enhancement"; then
        if "$dry_run"; then
            echo "  would create: feature"
        else
            gh label create --repo "$repo" feature --color a2eeef \
                --description "New feature or request" 2>/dev/null \
                && echo "  created: feature" \
                || echo "  FAILED: create feature" >&2
        fi
    fi

    # Create task
    if ! "$has_task"; then
        if "$dry_run"; then
            echo "  would create: task"
        else
            gh label create --repo "$repo" task --color "$TASK_COLOR" \
                --description "$TASK_DESC" 2>/dev/null \
                && echo "  created: task" \
                || echo "  FAILED: create task" >&2
        fi
    fi

    # Delete extras
    local label
    for label in "${to_delete[@]+"${to_delete[@]}"}"; do
        if "$dry_run"; then
            echo "  would delete: $label"
        else
            gh label delete --repo "$repo" "$label" --yes 2>/dev/null \
                && echo "  deleted: $label" \
                || echo "  FAILED: delete $label" >&2
        fi
    done
}

for repo in "${repos[@]}"; do
    reconcile_one "$repo"
done
