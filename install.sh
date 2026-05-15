#!/usr/bin/env bash
# =============================================================================
# install.sh — gocode v1.0.2 installer
# Run once: bash install.sh
# =============================================================================

GOCODE_DIR="$HOME/.gocode"
SCRIPTS_DIR="$HOME/.scripts"
SHELL_RC="$HOME/.bashrc"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

log_ok()    { echo -e "${GREEN}[OK]${RESET}    $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${RESET}  $1"; }
log_error() { echo -e "${RED}[ERROR]${RESET} $1"; }
log_step()  { echo -e "${CYAN}[....] $1${RESET}"; }

echo -e "${CYAN}${BOLD}"
echo "  Installing gocode v1.0.2..."
echo -e "${RESET}"

# ── 1. Verify we're running from the gocode directory ────────────────────────
if [ ! -f "./gocode" ] || [ ! -d "./lib" ]; then
    log_error "Run this script from inside the gocode directory."
    log_error "Example: cd ~/.gocode && bash install.sh"
    exit 1
fi

# ── 2. Copy to ~/.gocode if running from elsewhere ───────────────────────────
INSTALL_SRC="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ "$INSTALL_SRC" != "$GOCODE_DIR" ]; then
    log_step "Copying gocode to $GOCODE_DIR..."
    mkdir -p "$GOCODE_DIR"
    cp -r "$INSTALL_SRC"/. "$GOCODE_DIR/"
    log_ok "Copied to $GOCODE_DIR"
fi

# ── 3. Make scripts executable ───────────────────────────────────────────────
chmod +x "$GOCODE_DIR/gocode"
chmod +x "$GOCODE_DIR/lib/"*.sh
log_ok "Permissions set."

# ── 4. Symlink into ~/.scripts ────────────────────────────────────────────────
mkdir -p "$SCRIPTS_DIR"
ln -sf "$GOCODE_DIR/gocode" "$SCRIPTS_DIR/gocode"
log_ok "Symlinked: gocode → $SCRIPTS_DIR/gocode"

# ── 5. Add ~/.scripts to PATH ────────────────────────────────────────────────
if ! grep -q 'export PATH="$HOME/.scripts:$PATH"' "$SHELL_RC"; then
    echo '' >> "$SHELL_RC"
    echo '# gocode dev tools' >> "$SHELL_RC"
    echo 'export PATH="$HOME/.scripts:$PATH"' >> "$SHELL_RC"
    log_ok "Added ~/.scripts to PATH in $SHELL_RC"
else
    log_warn "PATH entry already exists — skipping."
fi

# ── 6. Install system dependencies ───────────────────────────────────────────
log_step "Checking system dependencies..."

# jq
if ! command -v jq &>/dev/null; then
    log_step "Installing jq..."
    sudo apt-get install -y jq &>/dev/null \
        && log_ok "jq installed." \
        || log_error "jq install failed. Run: sudo apt install jq"
else
    log_ok "jq already installed."
fi

# fzf — required for arrow key menus
if ! command -v fzf &>/dev/null; then
    log_step "Installing fzf (arrow key menus)..."
    sudo apt-get install -y fzf &>/dev/null \
        && log_ok "fzf installed." \
        || log_error "fzf install failed. Run: sudo apt install fzf"
else
    log_ok "fzf already installed."
fi

# ── 7. Init session directory ────────────────────────────────────────────────
mkdir -p "$HOME/.devsession"
if [ ! -f "$HOME/.devsession/projects.json" ]; then
    echo '{"projects":[]}' > "$HOME/.devsession/projects.json"
    log_ok "Session registry initialized."
else
    log_ok "Session registry exists — keeping data."
fi

# ── Done ─────────────────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}${BOLD}Installation complete.${RESET}"
echo ""
echo -e "  Run: ${CYAN}source ~/.bashrc${RESET}"
echo -e "  Then: ${CYAN}gocode${RESET}"
echo ""
