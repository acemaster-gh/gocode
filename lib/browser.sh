#!/usr/bin/env bash
# lib/browser.sh — gocode v1.0.3

browser_open_tabs() {
    local live_url="${1:-}"
    if ! command -v "${CHROME_BIN:-}" &>/dev/null; then
        log_warn "Browser not found: ${CHROME_BIN:-unset} — skipping."
        return
    fi
    log_step "Opening browser in ${SLEEP_BEFORE_CHROME:-8}s..."
    sleep "${SLEEP_BEFORE_CHROME:-8}"
    local tabs=("${BROWSER_TABS[@]:-}")
    [ -n "$live_url" ] && tabs+=("$live_url")
    "${CHROME_BIN}" --profile-directory="${CHROME_PROFILE:-Default}" \
        "${tabs[@]}" &>/dev/null &
    log_ok "Browser tabs opened."
}

browser_open_url() {
    local url="$1"
    if command -v "${CHROME_BIN:-}" &>/dev/null; then
        "${CHROME_BIN}" --profile-directory="${CHROME_PROFILE:-Default}" \
            "$url" &>/dev/null &
    else
        xdg-open "$url" &>/dev/null &
    fi
}

set_browser_interactive() {
    log_step "Detecting browsers..."
    local bin; bin=$(select_browser_interactive)
    [ -z "$bin" ] && log_warn "No browser selected." && return
    if grep -q "^CHROME_BIN=" "$CONFIG_FILE" 2>/dev/null; then
        sed -i "s|^CHROME_BIN=.*|CHROME_BIN=\"${bin}\"|" "$CONFIG_FILE"
    else
        echo "CHROME_BIN=\"${bin}\"" >> "$CONFIG_FILE"
    fi
    export CHROME_BIN="$bin"
    log_ok "Browser set → ${bin}"
}
