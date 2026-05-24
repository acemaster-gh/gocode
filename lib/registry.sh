#!/usr/bin/env bash
# =============================================================================
# lib/registry.sh — Project registry CRUD
# gocode v1.0.3 — atomic writes, backward-compatible field names
# =============================================================================

SESSION_DIR="$HOME/.devsession"
PROJECTS_REGISTRY="$SESSION_DIR/projects.json"

init_registry() {
    mkdir -p "$SESSION_DIR"
    [ ! -f "$PROJECTS_REGISTRY" ] && echo '{"projects":[]}' > "$PROJECTS_REGISTRY"
    _fix_null_types
}

# Atomic write — never corrupts on crash
_registry_write() {
    local tmp; tmp=$(mktemp)
    echo "$1" > "$tmp" && mv "$tmp" "$PROJECTS_REGISTRY"
}

# Migrate old records with null type
_fix_null_types() {
    local updated
    updated=$(jq '(.projects[] | select(.type == null) | .type) = "vanilla"' \
        "$PROJECTS_REGISTRY" 2>/dev/null) && _registry_write "$updated"
}

add_project_to_registry() {
    local name="$1" path="$2" type="$3"
    local repo_url="${4:-local-only}" netlify_url="${5:-not-deployed}"
    local description="${6:-}"
    local date; date=$(date +%Y-%m-%d)
    local updated
    updated=$(jq \
        --arg n  "$name"        --arg p "$path"         \
        --arg t  "$type"        --arg r "$repo_url"     \
        --arg nl "$netlify_url" --arg d "$date"         \
        --arg ds "$description" \
        '.projects += [{
            "name":$n,"path":$p,"type":$t,"status":"active",
            "created":$d,"last_opened":$d,
            "repo_url":$r,"netlify_url":$nl,"description":$ds
        }]' "$PROJECTS_REGISTRY")
    _registry_write "$updated"
}

update_project_field() {
    local name="$1" field="$2" value="$3"
    local updated
    updated=$(jq --arg n "$name" --arg f "$field" --arg v "$value" \
        '(.projects[] | select(.name==$n) | .[$f]) = $v' "$PROJECTS_REGISTRY")
    _registry_write "$updated"
}

get_project_field() {
    jq -r --arg n "$1" --arg f "$2" \
        '.projects[] | select(.name==$n) | .[$f] // "unknown"' \
        "$PROJECTS_REGISTRY" 2>/dev/null
}

get_incomplete_projects() {
    jq -r '.projects[] | select(.status=="active") | .name' \
        "$PROJECTS_REGISTRY" 2>/dev/null
}

get_last_active_project() {
    jq -r '[.projects[] | select(.status=="active")] | sort_by(.last_opened) | reverse | .[0].name' \
        "$PROJECTS_REGISTRY" 2>/dev/null
}

project_exists() {
    local n
    n=$(jq -r --arg n "$1" '.projects[] | select(.name==$n) | .name' \
        "$PROJECTS_REGISTRY" 2>/dev/null)
    [ -n "$n" ]
}

mark_project_complete() {
    update_project_field "$1" "status" "completed"
    log_ok "Marked '${1}' as completed."
}

archive_project() {
    update_project_field "$1" "status" "archived"
}

update_last_opened() {
    update_project_field "$1" "last_opened" "$(date +%Y-%m-%d)"
}

count_by_status() {
    jq --arg s "$1" '[.projects[] | select(.status==$s)] | length' \
        "$PROJECTS_REGISTRY" 2>/dev/null
}
