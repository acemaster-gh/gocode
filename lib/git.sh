#!/usr/bin/env bash
# =============================================================================
# lib/git.sh — Git and GitHub operations
# gocode v1.0.3
# =============================================================================

git_init_and_commit() {
    log_step "Initialising git..."
    git init -q
    git checkout -b main 2>/dev/null || git branch -M main 2>/dev/null || true
    git add .
    git commit -m "feat: initial commit — created with gocode v1.0.3" -q
    log_ok "Git initialised (branch: main)"
}

# $1=name $2=visibility(public|private) — returns repo URL or "local-only"
github_create_repo() {
    local name="$1" visibility="${2:-public}"

    if ! command -v gh &>/dev/null; then
        log_warn "gh not found — skipping GitHub repo creation."
        echo "local-only"; return
    fi

    if ! gh auth status &>/dev/null; then
        log_error "Not logged into GitHub. Run: gh auth login"
        echo "local-only"; return
    fi

    spinner_start "Creating GitHub repo (${visibility})..."
    local gh_out
    gh_out=$(gh repo create "$name" "--${visibility}" 2>&1)
    local gh_exit=$?
    spinner_stop

    if [ $gh_exit -ne 0 ]; then
        log_error "GitHub create failed: ${gh_out}"
        echo "local-only"; return
    fi

    local url="https://github.com/${GH_USER}/${name}"
    log_ok "Repo created: ${url}"

    # Add remote and push
    git remote add origin "${url}.git" 2>/dev/null || \
        git remote set-url origin "${url}.git"
    git branch -M main 2>/dev/null || true

    spinner_start "Pushing to GitHub..."
    if git push -u origin main -q 2>/dev/null; then
        spinner_stop
        log_ok "Pushed to GitHub."
    else
        spinner_stop
        log_warn "Push failed — run: gocode --push"
    fi

    echo "$url"
}

git_push_final() {
    local repo_url="$1"
    [ "$repo_url" = "local-only" ] && return
    git add . 2>/dev/null
    git diff --cached --quiet || git commit -m "docs: update README" -q
    git push -q 2>/dev/null || true
}

# Interactive commit + push workflow (used by --push)
git_push_workflow() {
    local project_path="$1"
    [ ! -d "$project_path" ] && log_error "Path not found: ${project_path}" && return 1
    cd "$project_path"

    if git diff --quiet && git diff --cached --quiet; then
        log_warn "Nothing to commit."; return 0
    fi

    divider
    local commit_type
    commit_type=$(select_commit_type)
    [ -z "$commit_type" ] && log_warn "No commit type selected." && return 1

    local commit_msg
    read -rp "$(echo -e "${YELLOW}[?]${RESET} Commit message: ")" commit_msg </dev/tty
    [ -z "$commit_msg" ] && log_error "Commit message cannot be empty." && return 1

    git add .
    git commit -m "${commit_type}: ${commit_msg}" -q
    spinner_start "Pushing to GitHub..."
    git push -q 2>/dev/null
    spinner_stop
    log_ok "Pushed: ${commit_type}: ${commit_msg}"
}

github_rename_repo() {
    local old_name="$1" new_name="$2"
    spinner_start "Renaming GitHub repo..."
    if gh repo rename "$new_name" -R "${GH_USER}/${old_name}" --yes 2>/dev/null; then
        spinner_stop
        local url="https://github.com/${GH_USER}/${new_name}"
        log_ok "Renamed → ${url}"
        echo "$url"
    else
        spinner_stop
        log_error "GitHub rename failed."
        echo ""
    fi
}

github_set_visibility() {
    local name="$1" visibility="$2"
    spinner_start "Setting repo to ${visibility}..."
    if gh repo edit "${GH_USER}/${name}" --visibility "$visibility" 2>/dev/null; then
        spinner_stop; log_ok "Repo is now ${visibility}."
    else
        spinner_stop; log_error "Visibility change failed."
    fi
}

github_set_description() {
    local name="$1" desc="$2"
    spinner_start "Updating description..."
    if gh repo edit "${GH_USER}/${name}" --description "$desc" 2>/dev/null; then
        spinner_stop; log_ok "Description updated."
    else
        spinner_stop; log_error "Description update failed."
    fi
}

# Background search for similar repos — shows results in terminal
github_search_similar() {
    local name="$1"
    local topic
    topic=$(echo "$name" | tr '-' ' ' | \
        sed 's/\b\(app\|project\|tool\|my\|the\|a\)\b//gi' | \
        tr -s ' ' | xargs)
    [ -z "$topic" ] && return

    spinner_start "Searching GitHub for similar projects..."
    local results
    results=$(timeout 5s gh search repos "$topic" \
        --language javascript --sort stars --limit 5 \
        --json fullName,stargazersCount,url 2>/dev/null \
        | jq -r '.[] | "⭐ \(.stargazersCount)\t\(.fullName)\t\(.url)"' 2>/dev/null)
    spinner_stop

    if [ -n "$results" ]; then
        echo "" >&2
        echo -e "${BOLD}  Similar projects on GitHub:${RESET}" >&2
        divider
        while IFS=$'\t' read -r stars fname url; do
            echo -e "  ${CYAN}${stars}${RESET}  ${BOLD}${fname}${RESET}" >&2
            echo -e "       ${DIM}${url}${RESET}" >&2
        done <<< "$results"
        divider
        echo "" >&2
    fi
}
