#!/usr/bin/env bash
# =============================================================================
# lib/obsidian.sh — Obsidian daily note and project note writing
# gocode v1.0.1
# =============================================================================

obsidian_daily_log() {
    local project_name="$1"
    local repo_url="${2:-local-only}"
    local netlify_url="${3:-not-deployed}"

    [[ -z "$OBSIDIAN_VAULT" || ! -d "$OBSIDIAN_VAULT" ]] && return 0

    local notes_dir="${OBSIDIAN_VAULT}/05-DAILY-NOTES"
    mkdir -p "$notes_dir"

    local today_file="${notes_dir}/$(date +%Y-%m-%d).md"
    local now
    now=$(date +%H:%M)

    if [[ ! -f "$today_file" ]]; then
        cat > "$today_file" <<EOF
# $(date +%Y-%m-%d) — Dev Log

## Summary


---

## Sessions

EOF
    fi

    cat >> "$today_file" <<EOF

### 🚀 ${project_name} — started at ${now}

- **Repo:** ${repo_url}
- **Live:** ${netlify_url}

EOF

    log_ok "Obsidian daily note updated: $(date +%Y-%m-%d).md"
}

obsidian_project_note() {
    local project_name="$1"
    local repo_url="${2:-local-only}"
    local netlify_url="${3:-not-deployed}"
    local project_type="${4:-vanilla}"

    [[ -z "$OBSIDIAN_VAULT" || ! -d "$OBSIDIAN_VAULT" ]] && return 0

    local projects_dir="${OBSIDIAN_VAULT}/01-PROJECTS/Active"
    mkdir -p "$projects_dir"

    local note_file="${projects_dir}/${project_name}.md"

    cat > "$note_file" <<EOF
---
project: ${project_name}
status: active
started: $(date +%Y-%m-%d)
repo: ${repo_url}
live: ${netlify_url}
stack: ${project_type}
tags: [javascript, project, active]
---

# ${project_name}

## 🔗 Links
- **GitHub:** ${repo_url}
- **Live:** ${netlify_url}

## 🎯 Goal


## 📋 Tasks
- [ ] 

## 🧠 Notes


## 📅 Session Log

| Date | Notes |
|------|-------|
| $(date +%Y-%m-%d) | Project created |

EOF

    log_ok "Obsidian project note created: ${project_name}.md"
}

obsidian_complete() {
    local project_name="$1"
    [[ -z "$OBSIDIAN_VAULT" || ! -d "$OBSIDIAN_VAULT" ]] && return 0

    local note_file="${OBSIDIAN_VAULT}/01-PROJECTS/Active/${project_name}.md"

    if [[ -f "$note_file" ]]; then
        sed -i "s/^status: active/status: completed/" "$note_file"
        sed -i "/^started:/a completed: $(date +%Y-%m-%d)" "$note_file"
        sed -i "s/tags: \[javascript, project, active\]/tags: [javascript, project, completed]/" "$note_file"
        printf "\n| $(date +%Y-%m-%d) | ✅ Project marked complete |\n" >> "$note_file"
        log_ok "Obsidian note marked as complete."
    fi
}
