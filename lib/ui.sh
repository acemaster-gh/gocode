#!/usr/bin/env bash
# =============================================================================
# ui.sh — Terminal UI: colors, banner, logging, fzf menus, spinner
# gocode v1.0.2
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
    echo -e "${RESET}${DIM}  Personal Developer Workflow System — v1.0.2${RESET}"
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

# ── Spinner ───────────────────────────────────────────────────────────────────
SPINNER_PID=""

spinner_start() {
    local label="${1:-Working...}"
    local frames=('⠋' '⠙' '⠹' '⠸' '⠼' '⠴' '⠦' '⠧' '⠇' '⠏')
    (
        local i=0
        while true; do
            printf "\r${CYAN}${frames[$i]}${RESET} %s" "$label"
            i=$(( (i+1) % ${#frames[@]} ))
            sleep 0.08
        done
    ) &
    SPINNER_PID=$!
    disown "$SPINNER_PID"
}

spinner_stop() {
    if [[ -n "$SPINNER_PID" ]]; then
        kill "$SPINNER_PID" 2>/dev/null
        wait "$SPINNER_PID" 2>/dev/null
        printf "\r\033[K"
        SPINNER_PID=""
    fi
}

# ── Summary box ───────────────────────────────────────────────────────────────
print_summary_box() {
    local project_name="$1"
    local project_type="$2"
    local project_path="$3"
    local repo_url="$4"
    local netlify_url="$5"

    echo ""
    echo -e "${GREEN}${BOLD}┌─────────────────────────────────────────────────┐${RESET}"
    echo -e "${GREEN}${BOLD}│  ✓ Project Ready                                │${RESET}"
    echo -e "${GREEN}${BOLD}├─────────────────────────────────────────────────┤${RESET}"
    printf "${GREEN}${BOLD}│${RESET}  %-10s ${CYAN}%-36s${RESET} ${GREEN}${BOLD}│${RESET}\n" "Name:"    "$project_name"
    printf "${GREEN}${BOLD}│${RESET}  %-10s ${CYAN}%-36s${RESET} ${GREEN}${BOLD}│${RESET}\n" "Type:"    "$project_type"
    printf "${GREEN}${BOLD}│${RESET}  %-10s ${DIM}%-36s${RESET} ${GREEN}${BOLD}│${RESET}\n" "Folder:"  "$project_path"
    printf "${GREEN}${BOLD}│${RESET}  %-10s ${DIM}%-36s${RESET} ${GREEN}${BOLD}│${RESET}\n" "Repo:"    "$repo_url"
    printf "${GREEN}${BOLD}│${RESET}  %-10s ${DIM}%-36s${RESET} ${GREEN}${BOLD}│${RESET}\n" "Netlify:" "$netlify_url"
    echo -e "${GREEN}${BOLD}└─────────────────────────────────────────────────┘${RESET}"
    echo ""
}

# ── fzf-based project type selector ──────────────────────────────────────────
select_project_type() {
    local selected
    selected=$(printf "Vanilla JS\nTailwind CSS\nReact + Vite" \
        | fzf --ansi \
              --prompt="  Project type: " \
              --pointer="▸" \
              --height=6 \
              --border=rounded \
              --color="prompt:cyan,pointer:green,hl:yellow" \
              --no-info \
              2>/dev/null)

    case "$selected" in
        "Vanilla JS")    echo "vanilla" ;;
        "Tailwind CSS")  echo "tailwind" ;;
        "React + Vite")  echo "react" ;;
        *)               echo "vanilla" ;;
    esac
}

# ── fzf-based commit type selector ───────────────────────────────────────────
select_commit_type() {
    local selected
    selected=$(printf "feat: new feature\nfix: bug fix\ndocs: documentation\nstyle: formatting\nrefactor: code restructure\nperf: performance\ntest: tests\nchore: maintenance" \
        | fzf --ansi \
              --prompt="  Commit type: " \
              --pointer="▸" \
              --height=12 \
              --border=rounded \
              --color="prompt:cyan,pointer:green" \
              --no-info \
              2>/dev/null)

    echo "$selected" | cut -d: -f1
}

# ── fzf-based project selector from list ─────────────────────────────────────
select_from_list() {
    local prompt="$1"
    shift
    local items=("$@")
    printf '%s\n' "${items[@]}" \
        | fzf --ansi \
              --prompt="  $prompt " \
              --pointer="▸" \
              --height=12 \
              --border=rounded \
              --color="prompt:cyan,pointer:green,hl:yellow" \
              --no-info \
              2>/dev/null
}
