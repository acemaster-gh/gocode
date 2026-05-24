#!/usr/bin/env bash
# =============================================================================
# lib/netlify.sh — Netlify deploy, link, redeploy
# gocode v1.0.3
# =============================================================================

netlify_deploy() {
    local project_name="$1"
    local project_path="$2"
    local repo_url="${3:-local-only}"

    if ! command -v netlify &>/dev/null; then
        log_warn "Netlify CLI not found — skipping deploy."
        echo "not-deployed"; return
    fi

    # Auth check — avoid false negative on "Already logged in via netlify config"
    local status_out
    status_out=$(netlify status 2>&1)
    if echo "$status_out" | grep -qiE "not logged|not authenticated|please log in|login required"; then
        log_error "Not logged into Netlify. Run: netlify login"
        echo "not-deployed"; return
    fi

    # Step 1 — create site via API (try named, fallback to auto)
    spinner_start "Creating Netlify site..."
    local site_json site_id site_url

    site_json=$(netlify api createSite \
        --data "{\"name\":\"${project_name}\"}" 2>/dev/null || echo "")
    site_id=$(echo  "$site_json" | jq -r '.id      // empty' 2>/dev/null || echo "")
    site_url=$(echo "$site_json" | jq -r '.ssl_url // .url // empty' 2>/dev/null || echo "")

    # If name taken, let Netlify auto-assign
    if [ -z "$site_id" ]; then
        site_json=$(netlify api createSite --data '{}' 2>/dev/null || echo "")
        site_id=$(echo  "$site_json" | jq -r '.id      // empty' 2>/dev/null || echo "")
        site_url=$(echo "$site_json" | jq -r '.ssl_url // .url // empty' 2>/dev/null || echo "")
    fi
    spinner_stop

    if [ -z "$site_id" ]; then
        log_error "Could not create Netlify site. Check: netlify status"
        echo "not-deployed"; return
    fi

    log_ok "Site created: ${site_url}"

    # Step 2 — deploy
    spinner_start "Deploying to Netlify..."
    local deploy_out
    deploy_out=$(cd "$project_path" && \
        netlify deploy --prod --dir . --site "$site_id" 2>&1)
    local deploy_exit=$?
    spinner_stop

    if [ $deploy_exit -ne 0 ] || \
       ! echo "$deploy_out" | grep -qiE "Published|Website URL|Live URL|netlify\.app"; then
        log_warn "Deploy may have failed. Output:"
        echo "$deploy_out" | tail -5 | while IFS= read -r line; do
            echo "    ${line}" >&2
        done
        echo "not-deployed"; return
    fi

    # Extract production URL (last netlify.app URL = prod, not draft)
    local live_url
    live_url=$(echo "$deploy_out" | grep -oP 'https://\S+\.netlify\.app' | tail -1)
    [ -z "$live_url" ] && live_url="$site_url"

    log_ok "Deployed: ${live_url}"

    # Step 3 — link site to project folder so git push auto-deploys
    if [ "$repo_url" != "local-only" ]; then
        spinner_start "Linking Netlify site to project..."
        cd "$project_path"
        netlify link --id "$site_id" &>/dev/null
        spinner_stop

        local state_file="$project_path/.netlify/state.json"
        if [ -f "$state_file" ]; then
            local linked_id; linked_id=$(jq -r '.siteId' "$state_file" 2>/dev/null)
            if [ "$linked_id" = "$site_id" ]; then
                log_ok "Netlify linked — git push will auto-deploy."
            else
                log_warn "Link mismatch. Run 'netlify link' manually."
            fi
        else
            log_warn ".netlify/state.json not found — run 'netlify link' manually."
        fi
    fi

    echo "$live_url"
}

netlify_redeploy() {
    local project_path="$1"
    [ ! -d "$project_path" ] && log_error "Path not found: ${project_path}" && return 1
    cd "$project_path"

    spinner_start "Redeploying to Netlify..."
    local out; out=$(netlify deploy --dir=. --prod 2>&1)
    spinner_stop

    local url
    url=$(echo "$out" | grep -oP 'https://\S+\.netlify\.app' | tail -1)
    if [ -n "$url" ]; then
        log_ok "Redeployed: ${url}"
    else
        log_error "Redeploy failed."; echo "$out" | tail -3 >&2
    fi
}
