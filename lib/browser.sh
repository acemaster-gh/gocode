#!/usr/bin/env bash
# =============================================================================
# lib/browser.sh — Chrome launcher and tab management
# gocode v1.0.1
# =============================================================================

browser_open_tabs() {
    if ! command -v "${CHROME_BIN:-google-chrome}" &>/dev/null; then
        log_warn "Chrome not found at: ${CHROME_BIN} — skipping browser launch."
        return 0
    fi

    log_step "Opening browser tabs in ${SLEEP_BEFORE_CHROME:-10}s..."
    sleep "${SLEEP_BEFORE_CHROME:-10}"

    "${CHROME_BIN}" \
        --profile-directory="${CHROME_PROFILE:-Default}" \
        "${BROWSER_TABS[@]}" \
        &>/dev/null &

    log_ok "Browser tabs opened."
}

browser_open_resume() {
    browser_open_tabs
}

browser_open_url() {
    local url="$1"
    if command -v "${CHROME_BIN:-google-chrome}" &>/dev/null; then
        "${CHROME_BIN}" \
            --profile-directory="${CHROME_PROFILE:-Default}" \
            "$url" \
            &>/dev/null &
    fi
}
