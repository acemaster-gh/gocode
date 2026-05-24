#!/usr/bin/env bash
# =============================================================================
# lib/readme_gen.sh — Project README generator
# gocode v1.0.3
# =============================================================================

_stack_list() {
    case "$1" in
        vanilla)  printf '%s\n' "- HTML" "- CSS" "- JavaScript (Vanilla)" ;;
        tailwind) printf '%s\n' "- HTML" "- Tailwind CSS" "- JavaScript" ;;
        react)    printf '%s\n' "- React" "- Vite" "- JavaScript" ;;
        *)        printf '%s\n' "- HTML" "- CSS" "- JavaScript" ;;
    esac
}

_ascii_preview() {
    case "$1" in
        vanilla|tailwind) cat <<'EOF'
┌──────────────────────────────────────────┐
│                                          │
│   Title / Header                         │
│                                          │
│   ┌──────────┐   ┌──────────┐           │
│   │  Button  │   │  Button  │           │
│   └──────────┘   └──────────┘           │
│                                          │
│   ┌──────────────────────────────────┐  │
│   │  Main content area               │  │
│   └──────────────────────────────────┘  │
│                                          │
└──────────────────────────────────────────┘
EOF
            ;;
        react) cat <<'EOF'
App
 ├── Header
 │    └── Nav
 ├── Main
 │    ├── Component
 │    └── Component
 └── Footer
EOF
            ;;
    esac
}

generate_project_readme() {
    local name="$1" type="$2" path="$3"
    local repo_url="${4:-}"
    local netlify_url="${5:-}"
    local description="${6:-}"

    local started; started=$(date +%Y-%m-%d)
    local stack; stack=$(_stack_list "$type")
    local ascii; ascii=$(_ascii_preview "$type")

    local live_section
    if [ -n "$netlify_url" ] && [ "$netlify_url" != "not-deployed" ]; then
        live_section="## Live Demo

**${netlify_url}**"
    else
        live_section="## Live Demo

_Not deployed yet._"
    fi

    local repo_section=""
    if [ -n "$repo_url" ] && [ "$repo_url" != "local-only" ]; then
        repo_section="## Repo

**${repo_url}**"
    fi

    local desc_text="${description:-_Add a short description of what this project does._}"

    cat > "$path/README.md" <<EOF
# ${name}

**Started:** ${started} &nbsp;|&nbsp; **Stack:** ${type}

---

## Description

${desc_text}

${live_section}

${repo_section}

## Built With

${stack}

## Preview

\`\`\`
${ascii}
\`\`\`

## What I Learned

_Add notes about what you learned building this._

## Progress

| Date | Update |
|------|--------|
| ${started} | Project created |

---

_Built with [gocode](https://github.com/acemaster-gh/gocode) v1.0.3 ⚡_
EOF

    log_ok "README.md generated."
}

# Append completion row to README progress table
update_readme_on_complete() {
    local project_path="$1"
    local date; date=$(date +%Y-%m-%d)
    local readme="$project_path/README.md"
    [ ! -f "$readme" ] && return
    # Insert after the "Project created" row
    sed -i "/Project created/a | ${date} | ✅ Project completed |" "$readme" 2>/dev/null
    log_ok "README updated with completion date."
}
