#!/usr/bin/env bash
# =============================================================================
# deps.sh — Dependency checker
# Part of gocode v1.0.1
# =============================================================================

check_dependencies() {
    log_step "Checking dependencies..."
    local missing=()

    for cmd in git gh code jq; do
        if ! command -v "$cmd" &>/dev/null; then
            missing+=("$cmd")
        fi
    done

    # Chrome — check configured binary
    if ! command -v "${CHROME_BIN:-google-chrome}" &>/dev/null; then
        missing+=("${CHROME_BIN:-google-chrome}")
    fi

    # jq — auto install silently
    if [[ " ${missing[*]} " =~ " jq " ]]; then
        log_warn "jq not found. Installing..."
        sudo apt-get install -y jq &>/dev/null
        if command -v jq &>/dev/null; then
            log_ok "jq installed."
            missing=("${missing[@]/jq}")
        fi
    fi

    # Netlify — soft check only
    if ! command -v netlify &>/dev/null; then
        log_warn "Netlify CLI not found — auto-deploy will be skipped."
        log_warn "Install: npm install -g netlify-cli && netlify login"
    fi

    if [ ${#missing[@]} -gt 0 ]; then
        log_error "Missing required dependencies: ${missing[*]}"
        log_error "Install them and re-run gocode."
        exit 1
    fi

    log_ok "All dependencies found."
}
