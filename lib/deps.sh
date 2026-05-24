#!/usr/bin/env bash
# lib/deps.sh — gocode v1.0.3

check_dependencies() {
    local missing=()
    for cmd in git gh jq fzf code; do
        command -v "$cmd" &>/dev/null || missing+=("$cmd")
    done

    if ! command -v jq &>/dev/null; then
        log_warn "jq not found — installing..."
        sudo apt-get install -y jq &>/dev/null && log_ok "jq installed." || true
        missing=("${missing[@]/jq}")
    fi

    if ! command -v fzf &>/dev/null; then
        log_warn "fzf not found — installing..."
        sudo apt-get install -y fzf &>/dev/null && log_ok "fzf installed." || true
        missing=("${missing[@]/fzf}")
    fi

    command -v netlify &>/dev/null || log_warn "netlify-cli not found — deploys will be skipped."

    if [ ${#missing[@]} -gt 0 ]; then
        log_error "Missing: ${missing[*]}"
        log_error "Run: bash install.sh"
        exit 1
    fi
}
