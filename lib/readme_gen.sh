#!/usr/bin/env bash
# =============================================================================
# lib/readme_gen.sh — Auto README generator for projects
# gocode v1.0.2
# =============================================================================

# ── Generate README for a project folder ──────────────────────────────────────
generate_project_readme() {
    local project_name="$1"
    local project_type="$2"
    local repo_url="$3"
    local netlify_url="${4:-not-deployed}"
    local date
    date=$(date +%Y-%m-%d)

    local live_line=""
    if [[ "$netlify_url" != "not-deployed" ]]; then
        live_line="[$project_name Live Demo]($netlify_url)"
    else
        live_line="Not deployed yet"
    fi

    local stack_line="HTML · CSS · JavaScript (Vanilla)"
    if [[ "$project_type" == "tailwind" ]]; then
        stack_line="HTML · Tailwind CSS · JavaScript (Vanilla)"
    elif [[ "$project_type" == "react" ]]; then
        stack_line="React · Vite · JavaScript"
    fi

    cat > README.md <<EOF
# $project_name

> **Stack:** $stack_line | **Started:** $date

## 🔗 Links

| Resource | URL |
|----------|-----|
| GitHub   | $repo_url |
| Live     | $live_line |

## 📖 Description

> Add a short description of what this project does.

## ✨ Features

- [ ] Feature 1
- [ ] Feature 2
- [ ] Feature 3

## 🚀 Getting Started

\`\`\`bash
git clone $repo_url
cd $project_name
# Open with VS Code Live Server
\`\`\`

## 🧠 What I Learned

> Add notes about what you learned building this.

## 📅 Progress

| Date | Update |
|------|--------|
| $date | Project created |

---

*Built with [gocode](https://github.com/acemaster-gh/gocode) ⚡*
EOF

    log_ok "README.md generated for $project_name"
}

# ── Update README when project completes ──────────────────────────────────────
update_readme_on_complete() {
    local project_path="$1"
    local date
    date=$(date +%Y-%m-%d)

    if [ -f "$project_path/README.md" ]; then
        # Append completion entry to progress table
        sed -i "/| $date | Project created |/a | $date | ✅ Project completed |" \
            "$project_path/README.md" 2>/dev/null
        log_ok "README updated with completion date."
    fi
}
