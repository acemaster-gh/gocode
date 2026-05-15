#!/usr/bin/env bash
# =============================================================================
# lib/netlify.sh — Netlify deploy and link operations
# gocode v1.0.2
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

    spinner_start "Deploying to Netlify..."
    cd "$project_path" || return 1

    local deploy_output
    deploy_output=$(netlify deploy --dir=. --prod 2>&1)
    local deploy_exit=$?
    spinner_stop

    if [[ $deploy_exit -ne 0 ]]; then
        log_warn "Netlify deploy failed: $deploy_output"
        NETLIFY_URL="not-deployed"
        export NETLIFY_URL
        return 0
    fi

    # Extract live URL
    NETLIFY_URL=$(echo "$deploy_output" | grep -oE 'https://[a-zA-Z0-9-]+\.netlify\.app' | head -1)
    NETLIFY_URL="${NETLIFY_URL:-not-deployed}"
    log_ok "Deployed: $NETLIFY_URL"

    # ── Get site ID — 3 methods ───────────────────────────────────────────────

    local site_id=""

    # Method 1: UUID pattern from deploy output
    site_id=$(echo "$deploy_output" \
        | grep -oE '[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}' \
        | head -1)

    # Method 2: site name from URL → query API with jq
    if [[ -z "$site_id" && "$NETLIFY_URL" != "not-deployed" ]]; then
        local site_name
        site_name=$(echo "$NETLIFY_URL" | sed 's|https://||' | sed 's|\.netlify\.app.*||')
        site_id=$(netlify api listSites 2>/dev/null \
            | jq -r --arg name "$site_name" \
                '.[] | select(.name == $name) | .id' 2>/dev/null \
            | head -1)
    fi

    # Method 3: most recently created site
    if [[ -z "$site_id" ]]; then
        site_id=$(netlify api listSites 2>/dev/null \
            | jq -r 'sort_by(.created_at) | reverse | .[0].id' 2>/dev/null)
    fi

    # ── Link site to GitHub ───────────────────────────────────────────────────
    if [[ "$repo_url" != "local-only" && -n "$site_id" ]]; then
        spinner_start "Linking Netlify to GitHub..."
        netlify link --id "$site_id" &>/dev/null
        spinner_stop

        # ── Verify link via .netlify/state.json ──────────────────────────────
        local state_file="$project_path/.netlify/state.json"
        if [[ -f "$state_file" ]]; then
            local linked_id
            linked_id=$(jq -r '.siteId' "$state_file" 2>/dev/null)
            if [[ "$linked_id" == "$site_id" ]]; then
                log_ok "Netlify linked and verified — git push will auto-deploy."
            else
                log_warn "Netlify link state mismatch. Run 'netlify link' manually."
            fi
        else
            log_warn "Netlify link may have failed — .netlify/state.json not found."
            log_warn "Run 'netlify link' manually inside the project folder."
        fi
    elif [[ "$repo_url" != "local-only" ]]; then
        log_warn "Could not detect site ID. Run 'netlify link' manually."
    fi

    export NETLIFY_URL
}

# ── Manual redeploy ───────────────────────────────────────────────────────────
netlify_redeploy() {
    local project_path="$1"

    if ! command -v netlify &>/dev/null; then
        log_error "Netlify CLI not found."
        return 1
    fi

    if [ ! -d "$project_path" ]; then
        log_error "Project folder not found: $project_path"
        return 1
    fi

    cd "$project_path" || return 1

    spinner_start "Redeploying to Netlify..."
    local output
    output=$(netlify deploy --dir=. --prod 2>&1)
    spinner_stop

    local url
    url=$(echo "$output" | grep -oE 'https://[a-zA-Z0-9-]+\.netlify\.app' | head -1)

    if [[ -n "$url" ]]; then
        log_ok "Redeployed: $url"
    else
        log_error "Redeploy failed: $output"
    fi
}
