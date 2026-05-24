#!/usr/bin/env bash
# =============================================================================
# lib/ui.sh — Terminal UI, colors, banner, logging, menus, spinner
# gocode v1.0.3
# =============================================================================

# =============================================================================
# TERMINAL DETECTION
# =============================================================================
detect_terminal() {
    TERM_COLS=$(tput cols  2>/dev/null || echo 80)
    TERM_LINES=$(tput lines 2>/dev/null || echo 24)
    local ncolors; ncolors=$(tput colors 2>/dev/null || echo 8)
    [ "$ncolors" -ge 256 ] && COLOR_SUPPORT=256 || COLOR_SUPPORT=8
    _setup_colors
}

_setup_colors() {
    if [ "${COLOR_SUPPORT:-8}" -ge 256 ]; then
        RED='\e[38;5;203m'; GREEN='\e[38;5;114m'
        YELLOW='\e[38;5;221m'; BLUE='\e[38;5;75m'
        CYAN='\e[38;5;117m'; MAGENTA='\e[38;5;183m'
    else
        RED='\033[0;31m'; GREEN='\033[0;32m'
        YELLOW='\033[1;33m'; BLUE='\033[0;34m'
        CYAN='\033[0;36m'; MAGENTA='\033[0;35m'
    fi
    BOLD='\033[1m'; DIM='\033[2m'; RESET='\033[0m'
}

detect_terminal

# =============================================================================
# BANNER — ASCII art on wide terminals, compact on narrow
# =============================================================================
print_banner() {
    echo ""
    if [ "${TERM_COLS:-80}" -ge 70 ]; then
        echo -e "${CYAN}${BOLD}"
        echo "  ██████╗  ██████╗  ██████╗ ██████╗ ██████╗ ███████╗"
        echo " ██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝"
        echo " ██║  ███╗██║   ██║██║     ██║   ██║██║  ██║█████╗  "
        echo " ██║   ██║██║   ██║██║     ██║   ██║██║  ██║██╔══╝  "
        echo " ╚██████╔╝╚██████╔╝╚██████╗╚██████╔╝██████╔╝███████╗"
        echo "  ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝"
        echo -e "${RESET}${DIM}  Developer Workflow System · v1.0.3${RESET}"
    else
        echo -e "${CYAN}${BOLD}  gocode v1.0.3${RESET}"
    fi
    echo ""
}

# =============================================================================
# LOG FUNCTIONS — all write to stderr so $() captures only return values
# =============================================================================
log_info()  { echo -e "${BLUE}[INFO]${RESET}  $1"       >&2; }
log_ok()    { echo -e "${GREEN}[OK]${RESET}    $1"       >&2; }
log_warn()  { echo -e "${YELLOW}[WARN]${RESET}  $1"      >&2; }
log_error() { echo -e "${RED}[ERROR]${RESET} $1"         >&2; }
log_step()  { echo -e "${CYAN}[....] ${BOLD}$1${RESET}"  >&2; }
log_done()  { echo -e "${GREEN}[DONE] ${BOLD}$1${RESET}" >&2; }
divider()   { echo -e "${DIM}────────────────────────────────────────────────────${RESET}" >&2; }

# =============================================================================
# CONFIRM  — simple y/N prompt
# =============================================================================
confirm() {
    local prompt="$1"
    local answer
    read -rp "$(echo -e "${YELLOW}[?]${RESET} ${prompt} [y/N]: ")" answer </dev/tty
    [[ "$answer" =~ ^[Yy]$ ]]
}

# =============================================================================
# SPINNER — writes to /dev/tty so \r works regardless of stdout capture
# =============================================================================
SPINNER_PID=""

spinner_start() {
    local label="${1:-Working...}"
    local frames=('⠋' '⠙' '⠹' '⠸' '⠼' '⠴' '⠦' '⠧' '⠇' '⠏')
    {
        local i=0
        while true; do
            printf "\r${CYAN}${frames[$i]}${RESET} %s" "$label" >/dev/tty
            i=$(( (i+1) % ${#frames[@]} ))
            sleep 0.08
        done
    } &
    SPINNER_PID=$!
    disown "$SPINNER_PID"
}

spinner_stop() {
    if [ -n "${SPINNER_PID:-}" ] && kill -0 "$SPINNER_PID" 2>/dev/null; then
        kill "$SPINNER_PID" 2>/dev/null
        wait "$SPINNER_PID" 2>/dev/null || true
    fi
    SPINNER_PID=""
    printf "\r\033[K" >/dev/tty
}

# =============================================================================
# SUMMARY BOX — no width padding math (avoids ANSI length bugs)
# =============================================================================
print_summary_box() {
    local name="$1" type="$2" path="$3"
    local repo="${4:-local-only}" netlify="${5:-not-deployed}"
    echo ""
    echo -e "${GREEN}${BOLD}╔══════════════════════════════════════════════════════╗${RESET}"
    echo -e "${GREEN}${BOLD}║${RESET}  ${GREEN}${BOLD}✓ Project Ready: ${name}${RESET}"
    echo -e "${GREEN}${BOLD}║${RESET}"
    echo -e "${GREEN}${BOLD}║${RESET}  ${DIM}Type:${RESET}    ${type}"
    echo -e "${GREEN}${BOLD}║${RESET}  ${DIM}Path:${RESET}    ${path}"
    [ "$repo" != "local-only" ]      && echo -e "${GREEN}${BOLD}║${RESET}  ${DIM}GitHub:${RESET}  ${repo}"
    [ "$netlify" != "not-deployed" ] && echo -e "${GREEN}${BOLD}║${RESET}  ${DIM}Live:${RESET}    ${netlify}"
    echo -e "${GREEN}${BOLD}╚══════════════════════════════════════════════════════╝${RESET}"
    echo ""
}

# =============================================================================
# FZF WRAPPER — writes to tempfile, avoids $() TTY swallow
# =============================================================================
_fzf_pick() {
    local prompt="$1"
    local height="${2:-10}"
    local header="${3:-}"
    local tmpfile; tmpfile=$(mktemp)

    local fzf_args=(
        --prompt="  ${prompt} "
        --pointer="▸"
        --height="$height"
        --border=rounded
        --color="prompt:cyan,pointer:green,hl:yellow"
        --no-info
    )
    [ -n "$header" ] && fzf_args+=(--header="  ${header}")

    fzf "${fzf_args[@]}" > "$tmpfile" 2>/dev/tty
    local result; result=$(cat "$tmpfile")
    rm -f "$tmpfile"
    echo "$result"
}

# =============================================================================
# SELECTORS
# =============================================================================
select_project_type() {
    local selected
    selected=$(printf "Vanilla JS\nTailwind CSS\nReact + Vite" \
        | _fzf_pick "Project type:" 6 "↑↓ to move  Enter to select")
    case "$selected" in
        "Vanilla JS")   echo "vanilla"  ;;
        "Tailwind CSS") echo "tailwind" ;;
        "React + Vite") echo "react"    ;;
        *)              echo "vanilla"  ;;
    esac
}

select_visibility() {
    local selected
    selected=$(printf "Public\nPrivate" \
        | _fzf_pick "Visibility:" 5 "GitHub repo visibility")
    echo "${selected,,}"   # lowercase
}

select_commit_type() {
    local selected
    selected=$(printf \
        "feat: new feature\nfix: bug fix\ndocs: documentation\nstyle: formatting\nrefactor: restructure\nperf: performance\ntest: tests\nchore: maintenance" \
        | _fzf_pick "Commit type:" 12 "Select commit type")
    echo "$selected" | cut -d: -f1
}

select_from_list() {
    local prompt="$1" header="$2"
    shift 2
    printf '%s\n' "$@" | _fzf_pick "$prompt" 15 "$header"
}
