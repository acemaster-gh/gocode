#!/usr/bin/env bash
# lib/obsidian.sh — gocode v1.0.3

obsidian_project_note() {
    local name="$1" type="$2" repo="${3:-}" netlify="${4:-}" desc="${5:-}"
    [ -z "${OBSIDIAN_VAULT:-}" ] && return
    local dir="$OBSIDIAN_VAULT/gocode-projects"
    mkdir -p "$dir"
    cat > "$dir/${name}.md" <<EOF
# ${name}

**Type:** ${type}
**Created:** $(date +%Y-%m-%d)
**Status:** active
**Description:** ${desc}

## Links
- GitHub: ${repo:-local-only}
- Live: ${netlify:-not-deployed}

## Notes

EOF
}

obsidian_daily_log() {
    [ -z "${OBSIDIAN_VAULT:-}" ] && return
    local daily="$OBSIDIAN_VAULT/daily"
    mkdir -p "$daily"
    local today; today=$(date +%Y-%m-%d)
    local file="$daily/${today}.md"
    [ -f "$file" ] || echo "# ${today}" > "$file"
    echo "- $(date +%H:%M) — $1" >> "$file"
}

obsidian_complete() {
    [ -z "${OBSIDIAN_VAULT:-}" ] && return
    local file="$OBSIDIAN_VAULT/gocode-projects/${1}.md"
    [ -f "$file" ] && sed -i 's/\*\*Status:\*\*.*/\*\*Status:\*\* complete/' "$file"
}

obsidian_quick_note() {
    [ -z "${OBSIDIAN_VAULT:-}" ] && return
    local file="$OBSIDIAN_VAULT/gocode-projects/${1}.md"
    [ -f "$file" ] && echo "- $(date +%H:%M) — $2" >> "$file"
}
