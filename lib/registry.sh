#!/usr/bin/env bash
# =============================================================================
# registry.sh — Project registry: read/write ~/.devsession/projects.json
# Part of gocode v1.0.1
# =============================================================================

SESSION_DIR="$HOME/.devsession"
PROJECTS_REGISTRY="$SESSION_DIR/projects.json"

init_registry() {
    mkdir -p "$SESSION_DIR"
    if [ ! -f "$PROJECTS_REGISTRY" ]; then
        echo '{"projects":[]}' > "$PROJECTS_REGISTRY"
    fi
}

get_incomplete_projects() {
    jq -r '.projects[] | select(.status == "active") | .name' "$PROJECTS_REGISTRY" 2>/dev/null
}

get_project_field() {
    local name="$1"
    local field="$2"
    jq -r --arg n "$name" --arg f "$field" \
        '.projects[] | select(.name == $n) | .[$f]' \
        "$PROJECTS_REGISTRY" 2>/dev/null
}

add_project_to_registry() {
    local name="$1"
    local path="$2"
    local type="$3"
    local repo_url="$4"
    local netlify_url="$5"
    local date
    date=$(date +%Y-%m-%d)

    local updated
    updated=$(jq \
        --arg name "$name" \
        --arg path "$path" \
        --arg type "$type" \
        --arg repo "$repo_url" \
        --arg netlify "$netlify_url" \
        --arg date "$date" \
        '.projects += [{
            "name": $name,
            "path": $path,
            "type": $type,
            "status": "active",
            "created": $date,
            "last_opened": $date,
            "repo_url": $repo,
            "netlify_url": $netlify
        }]' "$PROJECTS_REGISTRY")
    echo "$updated" > "$PROJECTS_REGISTRY"
}

mark_project_complete() {
    local name="$1"
    local updated
    updated=$(jq \
        --arg name "$name" \
        '(.projects[] | select(.name == $name) | .status) = "completed"' \
        "$PROJECTS_REGISTRY")
    echo "$updated" > "$PROJECTS_REGISTRY"
    log_ok "Marked '$name' as completed."
}

update_last_opened() {
    local name="$1"
    local date
    date=$(date +%Y-%m-%d)
    local updated
    updated=$(jq \
        --arg name "$name" \
        --arg date "$date" \
        '(.projects[] | select(.name == $name) | .last_opened) = $date' \
        "$PROJECTS_REGISTRY")
    echo "$updated" > "$PROJECTS_REGISTRY"
}

project_exists() {
    local name="$1"
    local exists
    exists=$(jq -r \
        --arg name "$name" \
        '.projects[] | select(.name == $name) | .name' \
        "$PROJECTS_REGISTRY" 2>/dev/null)
    [ -n "$exists" ]
}

count_by_status() {
    local status="$1"
    jq --arg s "$status" \
        '[.projects[] | select(.status == $s)] | length' \
        "$PROJECTS_REGISTRY" 2>/dev/null
}
