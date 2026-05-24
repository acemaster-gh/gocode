#!/usr/bin/env bash
# =============================================================================
# lib/config.sh — User configuration and first-run setup wizard
# gocode v1.0.3
# =============================================================================

CONFIG_FILE="$HOME/.devsession/config.sh"

load_config() {
    if [ -f "$CONFIG_FILE" ]; then
        # shellcheck source=/dev/null
        source "$CONFIG_FILE"
    else
        run_config_wizard
    fi
}

# =============================================================================
# BROWSER DETECTION — no duplicates
# =============================================================================
detect_browsers() {
    local found=() seen=()
    declare -A candidates=(
        ["Google Chrome"]="google-chrome google-chrome-stable"
        ["Brave"]="brave-browser brave"
        ["Firefox"]="firefox"
        ["Chromium"]="chromium-browser chromium"
        ["Microsoft Edge"]="microsoft-edge"
        ["Opera"]="opera"
        ["Vivaldi"]="vivaldi-stable vivaldi"
    )
    for name in "${!candidates[@]}"; do
        for bin in ${candidates[$name]}; do
            if command -v "$bin" &>/dev/null; then
                local real; real=$(command -v "$bin")
                local dup=false
                for s in "${seen[@]:-}"; do [ "$s" = "$real" ] && dup=true && break; done
                if [ "$dup" = false ]; then
                    found+=("${name}|${real}")
                    seen+=("$real")
                fi
                break
            fi
        done
    done
    [ ${#found[@]} -eq 0 ] && return 1
    printf '%s\n' "${found[@]}"
}

select_browser_interactive() {
    local -a browsers
    mapfile -t browsers < <(detect_browsers 2>/dev/null)

    if [ ${#browsers[@]} -eq 0 ]; then
        log_warn "No browsers detected. Set CHROME_BIN manually in config."
        echo ""; return
    fi

    if [ ${#browsers[@]} -eq 1 ]; then
        local bin; bin=$(echo "${browsers[0]}" | cut -d'|' -f2)
        local name; name=$(echo "${browsers[0]}" | cut -d'|' -f1)
        log_ok "Browser: ${name}"
        echo "$bin"; return
    fi

    local names=()
    for b in "${browsers[@]}"; do names+=("$(echo "$b" | cut -d'|' -f1)"); done
    names+=("None / skip")

    local chosen
    chosen=$(printf '%s\n' "${names[@]}" | _fzf_pick "Choose browser:" 10 "Select your default browser")

    [ "$chosen" = "None / skip" ] && echo "" && return
    for b in "${browsers[@]}"; do
        local n; n=$(echo "$b" | cut -d'|' -f1)
        [ "$n" = "$chosen" ] && echo "$(echo "$b" | cut -d'|' -f2)" && return
    done
    echo ""
}

# =============================================================================
# WIZARD
# =============================================================================
run_config_wizard() {
    print_banner
    log_info "First run — setting up gocode..."
    divider

    # GitHub username
    local GH_USER=""
    GH_USER=$(gh api user --jq '.login' 2>/dev/null || echo "")
    if [ -n "$GH_USER" ]; then
        log_ok "GitHub: ${GH_USER}"
    else
        read -rp "$(echo -e "${YELLOW}[?]${RESET} GitHub username: ")" GH_USER </dev/tty
    fi

    # Projects dir
    local default_dir="$HOME/Desktop/js-projects"
    read -rp "$(echo -e "${YELLOW}[?]${RESET} Projects folder [${default_dir}]: ")" input_dir </dev/tty
    local PROJECTS_DIR="${input_dir:-$default_dir}"
    mkdir -p "$PROJECTS_DIR"

    # Obsidian vault — auto-detect
    local OBSIDIAN_VAULT="$HOME/Desktop/DevMaster-Vault"
    local found_vault
    found_vault=$(find "$HOME" -maxdepth 4 -name ".obsidian" -type d 2>/dev/null | head -1)
    if [ -n "$found_vault" ]; then
        OBSIDIAN_VAULT=$(dirname "$found_vault")
        log_ok "Obsidian vault: ${OBSIDIAN_VAULT}"
    fi

    # Browser
    log_step "Detecting browsers..."
    local CHROME_BIN=""
    CHROME_BIN=$(select_browser_interactive)
    [ -n "$CHROME_BIN" ] && log_ok "Browser: ${CHROME_BIN}"

    # Write config
    mkdir -p "$HOME/.devsession"
    cat > "$CONFIG_FILE" <<EOF
# gocode v1.0.3 config — $(date +%Y-%m-%d)
GH_USER="${GH_USER}"
PROJECTS_DIR="${PROJECTS_DIR}"
CHROME_BIN="${CHROME_BIN}"
CHROME_PROFILE="Default"
SLEEP_BEFORE_CHROME=8
OBSIDIAN_VAULT="${OBSIDIAN_VAULT}"

BROWSER_TABS=(
    "https://claude.ai"
    "https://chatgpt.com"
    "https://developer.mozilla.org"
    "http://localhost:5500"
)
EOF

    log_ok "Config saved → ${CONFIG_FILE}"
    divider
    source "$CONFIG_FILE"
}
