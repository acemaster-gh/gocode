#!/usr/bin/env bash
# =============================================================================
# config.sh — User configuration and first-run setup wizard
# Part of gocode v1.0.1
# =============================================================================

CONFIG_FILE="$HOME/.devsession/config.sh"

# ── Load saved config if it exists ───────────────────────────────────────────
load_config() {
    if [ -f "$CONFIG_FILE" ]; then
        # shellcheck source=/dev/null
        source "$CONFIG_FILE"
    else
        run_config_wizard
    fi
}

# ── First-run wizard — writes config.sh once ─────────────────────────────────
run_config_wizard() {
    echo ""
    log_info "First run — setting up gocode..."
    divider

    # GitHub username — auto-detect, no prompt
    GH_USER=$(gh api user --jq '.login' 2>/dev/null)
    if [ -n "$GH_USER" ]; then
        log_ok "GitHub: $GH_USER"
    else
        read -rp "$(echo -e "${YELLOW}[?]${RESET} GitHub username: ")" GH_USER
    fi

    # Projects dir — use default, just confirm
    local default_projects="$HOME/Desktop/js-projects"
    read -rp "$(echo -e "${YELLOW}[?]${RESET} Projects folder [$default_projects]: ")" input_projects
    PROJECTS_DIR="${input_projects:-$default_projects}"

    # Chrome — auto-detect, no prompt
    CHROME_BIN=""
    for bin in google-chrome google-chrome-stable chromium chromium-browser; do
        if command -v "$bin" &>/dev/null; then
            CHROME_BIN="$bin"
            break
        fi
    done
    if [ -n "$CHROME_BIN" ]; then
        log_ok "Chrome: $CHROME_BIN"
    else
        read -rp "$(echo -e "${YELLOW}[?]${RESET} Chrome binary name: ")" CHROME_BIN
    fi

    # All other settings — hardcoded sensible defaults, no questions
    CHROME_PROFILE="Default"
    SLEEP_BEFORE_CHROME=10

    # Write config
    mkdir -p "$HOME/.devsession"
    cat > "$CONFIG_FILE" <<EOF
# gocode config — generated $(date +%Y-%m-%d)
# Edit this file to change settings
# Re-run wizard: gocode2 --config

GH_USER="$GH_USER"
PROJECTS_DIR="$PROJECTS_DIR"
CHROME_BIN="$CHROME_BIN"
CHROME_PROFILE="Default"
SLEEP_BEFORE_CHROME=10

OBSIDIAN_VAULT="$HOME/Desktop/DevMaster-Vault"

BROWSER_TABS=(
    "https://chatgpt.com"
    "https://claude.ai"
    "https://developer.mozilla.org"
    "http://localhost:5500"
)
EOF

    log_ok "Config saved."
    divider
    source "$CONFIG_FILE"
}
