#!/usr/bin/env bash
# =============================================================================
# install.sh — gocode v1.0.3
# =============================================================================
set -euo pipefail

GOCODE_DIR="$HOME/.gocode"
SCRIPTS_DIR="$HOME/.scripts"
SHELL_RC="$HOME/.bashrc"

G='\033[0;32m'; Y='\033[1;33m'; R='\033[0;31m'; C='\033[0;36m'; B='\033[1m'; RS='\033[0m'
ok()   { echo -e "${G}[OK]${RS}    $1"; }
warn() { echo -e "${Y}[WARN]${RS}  $1"; }
err()  { echo -e "${R}[ERROR]${RS} $1"; }
step() { echo -e "${C}[....]${RS}  $1"; }

echo -e "\n${C}${B}  gocode v1.0.3 — installer${RS}\n"

[ "$EUID" -eq 0 ] && err "Do not run as root." && exit 1

INSTALL_SRC="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ── Copy to ~/.gocode ─────────────────────────────────────────────────────────
if [ "$INSTALL_SRC" != "$GOCODE_DIR" ]; then
    step "Copying to ${GOCODE_DIR}..."
    mkdir -p "$GOCODE_DIR"
    cp -r "$INSTALL_SRC/." "$GOCODE_DIR/"
    ok "Copied to ${GOCODE_DIR}"
fi

chmod +x "$GOCODE_DIR/gocode"
chmod +x "$GOCODE_DIR/lib/"*.sh
ok "Permissions set."

# ── Symlink ───────────────────────────────────────────────────────────────────
mkdir -p "$SCRIPTS_DIR"
ln -sf "$GOCODE_DIR/gocode" "$SCRIPTS_DIR/gocode"
ok "Symlinked → ${SCRIPTS_DIR}/gocode"

if ! grep -q 'HOME/.scripts' "$SHELL_RC" 2>/dev/null; then
    echo '' >> "$SHELL_RC"
    echo '# gocode' >> "$SHELL_RC"
    echo 'export PATH="$HOME/.scripts:$PATH"' >> "$SHELL_RC"
    ok "Added ~/.scripts to PATH in ${SHELL_RC}"
fi

# ── System packages ───────────────────────────────────────────────────────────
step "Checking system packages..."
sudo apt-get update -qq

for pkg in git jq fzf curl; do
    if dpkg -s "$pkg" &>/dev/null; then
        ok "${pkg} already installed"
    else
        step "Installing ${pkg}..."
        sudo apt-get install -y "$pkg" -q && ok "${pkg} installed"
    fi
done

# ── GitHub CLI ────────────────────────────────────────────────────────────────
if ! command -v gh &>/dev/null; then
    step "Installing GitHub CLI..."
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
        | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg 2>/dev/null
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] \
https://cli.github.com/packages stable main" \
        | sudo tee /etc/apt/sources.list.d/github-cli.list >/dev/null
    sudo apt-get update -qq && sudo apt-get install -y gh -q
    ok "gh installed"
else
    ok "gh already installed"
fi

# ── Node.js ───────────────────────────────────────────────────────────────────
if ! command -v node &>/dev/null; then
    step "Installing Node.js LTS..."
    curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash - -q 2>/dev/null
    sudo apt-get install -y nodejs -q
    ok "Node.js $(node -v) installed"
else
    ok "Node.js $(node -v) already installed"
fi

# ── Netlify CLI ───────────────────────────────────────────────────────────────
if ! command -v netlify &>/dev/null; then
    step "Installing netlify-cli..."
    npm install -g netlify-cli --silent
    ok "netlify-cli installed"
else
    ok "netlify-cli already installed"
fi

# ── Auth ──────────────────────────────────────────────────────────────────────
echo ""
echo -e "${B}  Authentication${RS}"
echo ""

if gh auth status &>/dev/null; then
    ok "GitHub — already authenticated"
else
    warn "GitHub login required"
    gh auth login
fi

echo ""

if netlify status 2>&1 | grep -qiE "Logged in|Already logged"; then
    ok "Netlify — already authenticated"
else
    warn "Netlify login required"
    netlify login
fi

# ── Install @clack/prompts ────────────────────────────────────────────────────
if [ -f "$GOCODE_DIR/package.json" ]; then
    step "Installing @clack/prompts..."
    cd "$GOCODE_DIR" && npm install --silent
    ok "@clack/prompts installed"
fi

# ── Registry ──────────────────────────────────────────────────────────────────
mkdir -p "$GOCODE_DIR/data"
[ -f "$GOCODE_DIR/data/projects.json" ] || echo '{"projects":[]}' > "$GOCODE_DIR/data/projects.json"
mkdir -p "$HOME/.devsession"
[ -f "$HOME/.devsession/projects.json" ] || echo '{"projects":[]}' > "$HOME/.devsession/projects.json"
ok "Registry initialised"

# ── Done ──────────────────────────────────────────────────────────────────────
echo ""
echo -e "${G}${B}  Installation complete!${RS}"
echo ""
echo -e "  Reload shell:   ${C}source ~/.bashrc${RS}"
echo -e "  Setup wizard:   ${C}gocode --config${RS}"
echo -e "  Create project: ${C}gocode --new${RS}"
echo ""
