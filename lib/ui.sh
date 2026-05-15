#!/usr/bin/env bash
# =============================================================================
# ui.sh — Terminal UI: colors, banner, logging, fzf menus, spinner
# gocode v1.0.2
# =============================================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'

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

log_info()  { echo -e "${BLUE}[INFO]${RESET}  $1"; }
log_ok()    { echo -e "${GREEN}[OK]${RESET}    $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${RESET}  $1"; }
log_error() { echo -e "${RED}[ERROR]${RESET} $1"; }
log_step()  { echo -e "${CYAN}[....] ${BOLD}$1${RESET}"; }
log_done()  { echo -e "${GREEN}[DONE] ${BOLD}$1${RESET}"; }
divider()   { echo -e "${DIM}────────────────────────────────────────────────────${RESET}"; }

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
    local project_name="$1" project_type="$2"
    local project_path="$3" repo_url="$4" netlify_url="$5"
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

# ── fzf wrapper — writes to tempfile, avoids $() TTY swallow ─────────────────
_fzf_pick() {
    local prompt="$1"
    local height="${2:-10}"
    local header="${3:-}"
    local color="${4:-prompt:cyan,pointer:green,hl:yellow}"
    local tmpfile
    tmpfile=$(mktemp)

    local fzf_args=(
        --prompt="  $prompt "
        --pointer="▸"
        --height="$height"
        --border=rounded
        --color="$color"
        --no-info
    )
    [ -n "$header" ] && fzf_args+=(--header="  $header")

    fzf "${fzf_args[@]}" > "$tmpfile"
    local result
    result=$(cat "$tmpfile")
    rm -f "$tmpfile"
    echo "$result"
}

# ── Project type selector ─────────────────────────────────────────────────────
select_project_type() {
    local selected
    selected=$(printf "Vanilla JS\nTailwind CSS\nReact + Vite" \
        | _fzf_pick "Project type:" 6)
    case "$selected" in
        "Vanilla JS")   echo "vanilla" ;;
        "Tailwind CSS") echo "tailwind" ;;
        "React + Vite") echo "react" ;;
        *)              echo "vanilla" ;;
    esac
}

# ── Commit type selector ──────────────────────────────────────────────────────
select_commit_type() {
    local selected
    selected=$(printf "feat: new feature\nfix: bug fix\ndocs: documentation\nstyle: formatting\nrefactor: code restructure\nperf: performance\ntest: tests\nchore: maintenance" \
        | _fzf_pick "Commit type:" 12)
    echo "$selected" | cut -d: -f1
}

# ── Generic list picker ───────────────────────────────────────────────────────
select_from_list() {
    local prompt="$1"
    local header="$2"
    shift 2
    printf '%s\n' "$@" | _fzf_pick "$prompt" 15 "$header"
}
