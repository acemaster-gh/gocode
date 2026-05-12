#!/usr/bin/env bash
# =============================================================================
# ui.sh — Terminal UI: colors, banner, logging, menus
# Part of gocode v1.0.1
# =============================================================================

# ── Colors ────────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'

# ── Banner ────────────────────────────────────────────────────────────────────
print_banner() {
    echo -e "${CYAN}${BOLD}"
    echo "  ██████╗  ██████╗  ██████╗ ██████╗ ██████╗ ███████╗"
    echo " ██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝"
    echo " ██║  ███╗██║   ██║██║     ██║   ██║██║  ██║█████╗  "
    echo " ██║   ██║██║   ██║██║     ██║   ██║██║  ██║██╔══╝  "
    echo " ╚██████╔╝╚██████╔╝╚██████╗╚██████╔╝██████╔╝███████╗"
    echo "  ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝"
    echo -e "${RESET}${DIM}  Personal Developer Workflow System — v1.0.1${RESET}"
    echo ""
}

# ── Log functions ─────────────────────────────────────────────────────────────
log_info()  { echo -e "${BLUE}[INFO]${RESET}  $1"; }
log_ok()    { echo -e "${GREEN}[OK]${RESET}    $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${RESET}  $1"; }
log_error() { echo -e "${RED}[ERROR]${RESET} $1"; }
log_step()  { echo -e "${CYAN}[....] ${BOLD}$1${RESET}"; }
log_done()  { echo -e "${GREEN}[DONE] ${BOLD}$1${RESET}"; }

divider()   { echo -e "${DIM}────────────────────────────────────────────────────${RESET}"; }

# ── Confirm prompt ────────────────────────────────────────────────────────────
confirm() {
    local prompt="$1"
    read -rp "$(echo -e "${YELLOW}[?]${RESET} $prompt [y/N]: ")" answer
    [[ "$answer" =~ ^[Yy]$ ]]
}

# ── Project type selector ─────────────────────────────────────────────────────
select_project_type() {
    echo "" >&2
    echo -e "${BOLD}Select project type:${RESET}" >&2
    echo -e "  ${CYAN}1${RESET}. Vanilla JS" >&2
    echo -e "  ${CYAN}2${RESET}. Tailwind CSS" >&2
    echo -e "  ${CYAN}3${RESET}. React + Vite ${DIM}(coming in v1.0.2)${RESET}" >&2
    echo "" >&2

    local choice
    read -rp "$(echo -e "${YELLOW}[?]${RESET} Choose [1]: ")" choice >&2
    choice="${choice:-1}"

    case "$choice" in
        1) echo "vanilla" ;;
        2) echo "tailwind" ;;
        3) echo "react" ;;
        *) echo "vanilla" ;;
    esac
}
