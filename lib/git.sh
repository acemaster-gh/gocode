#!/usr/bin/env bash
# =============================================================================
# git.sh — Git and GitHub functions
# Part of gocode v1.0.1
# =============================================================================

git_init_and_commit() {
    log_step "Initializing git..."
    git init -q
    git add .
    git commit -m "feat: initial commit" -q
    log_ok "Git initialized and committed."
}

github_create_repo() {
    local project_name="$1"
    log_step "Creating GitHub repo..." >&2

    local gh_output
    gh_output=$(gh repo create "$project_name" --public 2>&1)

    if echo "$gh_output" | grep -q "github.com"; then
        local repo_url
        repo_url=$(gh repo view "$project_name" --json url -q '.url' 2>/dev/null)
        git remote add origin "https://github.com/${GH_USER}/${project_name}.git" 2>/dev/null
        git branch -M main 2>/dev/null
        git push -u origin main -q 2>/dev/null
        log_ok "GitHub repo created: $repo_url" >&2
        log_ok "Initial commit pushed." >&2
        echo "$repo_url"
    else
        log_error "GitHub repo creation failed: $gh_output" >&2
        log_warn "Continuing with local git only." >&2
        echo "local-only"
    fi
}

git_push_final() {
    local repo_url="$1"
    if [ "$repo_url" != "local-only" ]; then
        git add . 2>/dev/null
        git diff --cached --quiet || git commit -m "docs: update readme with links" -q
        git push -q 2>/dev/null
    fi
}

# ── GitHub similar projects search ───────────────────────────────────────────
github_search_similar() {
    local project_name="$1"
    local topic
    topic=$(echo "$project_name" | tr '-' ' ' | \
        sed 's/\b\(app\|project\|tool\|my\|the\|a\)\b//gi' | \
        tr -s ' ' | xargs)

    [ -z "$topic" ] && return

    log_step "Searching GitHub for similar projects..."

    local results
    results=$(timeout 5s gh api \
        "search/repositories?q=${topic}+language:javascript&sort=stars&order=desc&per_page=5" \
        --jq '.items[] | "  ⭐ \(.stargazers_count)\t\(.full_name)\t\(.html_url)"' \
        2>/dev/null)

    if [ -n "$results" ]; then
        echo -e "${BOLD}Similar projects on GitHub:${RESET}"
        divider
        while IFS=$'\t' read -r stars name url; do
            echo -e "  ${CYAN}${stars}${RESET}  ${BOLD}${name}${RESET}"
            echo -e "       ${DIM}${url}${RESET}"
        done <<< "$results"
        divider
    else
        log_warn "GitHub search unavailable — skipping."
    fi
    echo ""
}
