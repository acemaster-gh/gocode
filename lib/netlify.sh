#!/usr/bin/env bash
# =============================================================================
# lib/netlify.sh — Netlify deploy and link operations
# gocode v1.0.1
# =============================================================================

netlify_deploy() {
    local project_name="$1"
    local project_path="$2"
    local repo_url="$3"

    if ! command -v netlify &>/dev/null; then
        log_warn "Netlify CLI not found — skipping deploy."
        NETLIFY_URL="not-deployed"
        export NETLIFY_URL
        return 0
    fi

    log_step "Deploying to Netlify..."
    cd "$project_path" || return 1

    local deploy_output
    deploy_output=$(netlify deploy --dir=. --prod 2>&1)

    if [[ $? -ne 0 ]]; then
        log_warn "Netlify deploy failed: $deploy_output"
        NETLIFY_URL="not-deployed"
        export NETLIFY_URL
        return 0
    fi

    # Extract live URL
    NETLIFY_URL=$(echo "$deploy_output" | grep -oE 'https://[a-zA-Z0-9-]+\.netlify\.app' | head -1)
    NETLIFY_URL="${NETLIFY_URL:-not-deployed}"
    log_ok "Deployed: $NETLIFY_URL"

    # Method 1: extract site ID directly from deploy output line
    local site_id
    site_id=$(echo "$deploy_output" | grep -i 'site id' | grep -oE '[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}' | head -1)

    # Method 2: extract site name from URL using sed, query API
    if [[ -z "$site_id" && "$NETLIFY_URL" != "not-deployed" ]]; then
        local site_name
        site_name=$(echo "$NETLIFY_URL" | sed 's|https://||' | sed 's|\.netlify\.app.*||')
        site_id=$(netlify api listSites 2>/dev/null \
            | jq -r --arg name "$site_name" '.[] | select(.name == $name) | .id' 2>/dev/null | head -1)
    fi

    # Method 3: get most recently created site from API
    if [[ -z "$site_id" ]]; then
        site_id=$(netlify api listSites 2>/dev/null \
            | jq -r 'sort_by(.created_at) | reverse | .[0].id' 2>/dev/null)
    fi

    # Link using site ID
    if [[ "$repo_url" != "local-only" && -n "$site_id" ]]; then
        log_step "Linking Netlify to GitHub for auto-deploy..."
        netlify link --id "$site_id" &>/dev/null \
            && log_ok "Netlify linked — git push will auto-deploy." \
            || log_warn "Netlify link failed. Run 'netlify link' manually."
    elif [[ "$repo_url" != "local-only" ]]; then
        log_warn "Could not detect site ID. Run 'netlify link' manually."
    fi

    export NETLIFY_URL
}
