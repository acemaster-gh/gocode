#!/usr/bin/env bash
# =============================================================================
# lib/git.sh — Git, GitHub, push workflow
# gocode v1.0.2
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

    spinner_start "Creating GitHub repo..."
    local gh_output
    gh_output=$(gh repo create "$project_name" --public 2>&1)
    spinner_stop

    if echo "$gh_output" | grep -q "github.com"; then
        local repo_url
        repo_url=$(gh repo view "$project_name" --json url -q '.url' 2>/dev/null)
        git remote add origin "https://github.com/${GH_USER}/${project_name}.git" 2>/dev/null
        git branch -M main 2>/dev/null
        git push -u origin main -q 2>/dev/null
        log_ok "GitHub repo: $repo_url" >&2
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
        git diff --cached --quiet || git commit -m "docs: update readme" -q
        git push -q 2>/dev/null
    fi
}

# ── gocode --push workflow ────────────────────────────────────────────────────
git_push_workflow() {
    # Find active project path
    local project_name
    project_name=$(jq -r '.projects[] | select(.status == "active") | .name' \
        "$PROJECTS_REGISTRY" 2>/dev/null \
        | _fzf_pick "Select project:" 8)

    if [ -z "$project_name" ]; then
        log_warn "No project selected."
        return 1
    fi

    local project_path
    project_path=$(jq -r --arg n "$project_name" \
        '.projects[] | select(.name == $n) | .path' \
        "$PROJECTS_REGISTRY" 2>/dev/null)

    if [ ! -d "$project_path" ]; then
        log_error "Project folder not found: $project_path"
        return 1
    fi

    cd "$project_path" || return 1

    # Check if there's anything to commit
    if git diff --quiet && git diff --cached --quiet; then
        log_warn "Nothing to commit in $project_name."
        return 0
    fi

    divider
    echo -e "${BOLD}  Committing: $project_name${RESET}"
    divider

    # Arrow key commit type selection
    local commit_type
    commit_type=$(select_commit_type)

    if [ -z "$commit_type" ]; then
        log_warn "No commit type selected."
        return 1
    fi

    # Commit message
    local commit_msg
    read -rp "$(echo -e "${YELLOW}[?]${RESET} Commit message: ")" commit_msg

    if [ -z "$commit_msg" ]; then
        log_error "Commit message cannot be empty."
        return 1
    fi

    local full_message="${commit_type}: ${commit_msg}"

    # Add, commit, push
    git add .
    git commit -m "$full_message" -q
    spinner_start "Pushing to GitHub..."
    git push -q 2>/dev/null
    spinner_stop

    log_ok "Pushed: $full_message"
}

# ── GitHub similar projects search ───────────────────────────────────────────
github_search_similar() {
    local project_name="$1"
    local topic
    topic=$(echo "$project_name" | tr '-' ' ' | \
        sed 's/\b\(app\|project\|tool\|my\|the\|a\)\b//gi' | \
        tr -s ' ' | xargs)

    [ -z "$topic" ] && return

    spinner_start "Searching GitHub for similar projects..."
    local results
    results=$(timeout 5s gh search repos "$topic" \
        --language javascript \
        --sort stars \
        --limit 5 \
        --json fullName,stargazersCount,url \
        2>/dev/null \
        | jq -r '.[] | "⭐ \(.stargazersCount)\t\(.fullName)\t\(.url)"' 2>/dev/null)
    spinner_stop

    if [ -n "$results" ]; then
        echo ""
        echo -e "${BOLD}Similar projects on GitHub:${RESET}"
        divider
        while IFS=$'\t' read -r stars name url; do
            echo -e "  ${CYAN}${stars}${RESET}  ${BOLD}${name}${RESET}"
            echo -e "       ${DIM}${url}${RESET}"
        done <<< "$results"
        divider
        echo ""
    fi
}
