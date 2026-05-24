#!/usr/bin/env bash
# =============================================================================
# lib/prompts.sh — Input prompts: text, confirm, select (fzf-backed)
# gocode v1.0.3
# =============================================================================

# Text input — accepts placeholder as default on empty Enter
prompt_text() {
    local message="$1"
    local placeholder="${2:-}"
    local result=""

    while true; do
        if [ -n "$placeholder" ]; then
            read -rp "$(echo -e "${CYAN}[?]${RESET} ${message} ${DIM}[${placeholder}]${RESET}: ")" result </dev/tty
            [ -z "$result" ] && result="$placeholder"
        else
            read -rp "$(echo -e "${CYAN}[?]${RESET} ${message}: ")" result </dev/tty
        fi
        [ -z "$result" ] && { echo -e "${RED}[!]${RESET} Cannot be empty." >&2; continue; }
        break
    done

    echo "$result"
}

# Yes/No — returns 0 (yes) or 1 (no)
prompt_confirm() {
    local message="$1"
    local default="${2:-n}"
    local hint
    [ "$default" = "y" ] && hint="Y/n" || hint="y/N"
    local answer
    read -rp "$(echo -e "${CYAN}[?]${RESET} ${message} [${hint}]: ")" answer </dev/tty
    answer="${answer:-$default}"
    [[ "$answer" =~ ^[Yy]$ ]]
}

# Select from list — wraps fzf
prompt_select() {
    local message="$1"; shift
    printf '%s\n' "$@" | _fzf_pick "${message}" 12 "↑↓ navigate  Enter select"
}
