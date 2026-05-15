# gocode ⚡

> A modular, single-command developer workflow system for Linux — built to eliminate the repetitive overhead of starting, managing, and tracking JavaScript projects.

![version](https://img.shields.io/badge/version-1.0.2-cyan)
![platform](https://img.shields.io/badge/platform-Linux-blue)
![shell](https://img.shields.io/badge/shell-bash-green)
![license](https://img.shields.io/badge/license-MIT-yellow)

---

## What is gocode?

Every time you start a new project you do the same things manually — create a folder, set up files, initialize git, create a GitHub repo, open VS Code, open browser tabs, deploy to Netlify. Every single time.

`gocode` automates all of it into one command. One word in your terminal and your entire coding environment is ready in under 30 seconds — with a unique color palette, dark mode, structured boilerplate, and live deployment already set up.

---

## Architecture

```
~/.gocode/
├── gocode                  # Main entry point — thin router
├── install.sh              # One-time installer
├── lib/
│   ├── config.sh           # First-run wizard + user settings
│   ├── ui.sh               # Colors, banner, fzf menus, spinner
│   ├── deps.sh             # Dependency checker
│   ├── registry.sh         # projects.json read/write
│   ├── git.sh              # Git, GitHub, push workflow
│   ├── browser.sh          # Chrome tab launcher
│   ├── netlify.sh          # Deploy + auto-link + verify
│   ├── obsidian.sh         # Daily note + project note writer
│   ├── colors.sh           # Random cohesive palette generator
│   └── readme_gen.sh       # Auto README generator per project
├── boilerplates/
│   ├── vanilla/            # HTML + CSS (full :root vars + dark mode) + JS
│   ├── tailwind/           # Tailwind CDN + dark mode toggle
│   └── react/              # React + Vite scaffold
└── plugins/                # Drop a .sh file here — auto-loaded on run
```

---

## Features

### Project Creation
- Arrow key project type selector — Vanilla JS, Tailwind CSS, React + Vite
- Validates name, blocks duplicates
- Injects production-ready boilerplate — full `:root` design tokens, dark mode toggle wired into HTML + CSS + JS
- **Random cohesive color palette** — every project gets a unique but professional color scheme generated from HSL
- Git init → commit → GitHub repo → push — fully automated
- Netlify deploy + auto-link to GitHub verified via `.netlify/state.json`
- VS Code opens immediately after git init — no waiting for GitHub/Netlify
- Auto README generated per project — editable, updates on complete

### Session Management
- Arrow key project selector on launch — resume any active project
- Tracks all projects in `~/.devsession/projects.json`
- `gocode --open` — reopen last project instantly, no menu
- `gocode --complete` — mark project done, updates Obsidian + README
- `gocode --delete` — removes local folder, GitHub repo, and registry in one command

### Git Workflow
- `gocode --push` — arrow key commit type selector + message prompt + push
- Commit types: feat, fix, docs, style, refactor, perf, test, chore

### Netlify
- Auto-deploy on create with verified GitHub link
- `gocode --deploy` — manual redeploy of current project

### GitHub Search
- On new project create, searches GitHub for similar JS repos
- Shows stars + repo URL — inspiration and reference in seconds

### Browser Automation
- Opens ChatGPT, Claude, MDN, localhost:5500 in your Chrome profile
- 10 second delay after VS Code opens

### Obsidian Integration
- Daily note in `05-DAILY-NOTES/YYYY-MM-DD.md` — appended on every create/resume
- Individual project note in `01-PROJECTS/Active/` with repo URL, Netlify URL, task list, session log
- Notes marked complete when project is marked done

### Plugin System
- Drop any `.sh` file into `plugins/` — gocode sources it automatically on next run

---

## Installation

### Requirements

| Tool | Purpose | Install |
|------|---------|---------|
| `git` | Version control | `sudo apt install git` |
| `gh` | GitHub CLI | `sudo apt install gh` then `gh auth login` |
| `code` | VS Code | [code.visualstudio.com](https://code.visualstudio.com) |
| `google-chrome` | Browser automation | [google.com/chrome](https://google.com/chrome) |
| `jq` | JSON processor | Auto-installed by installer |
| `fzf` | Arrow key menus | Auto-installed by installer |
| `node` + `npm` | For Netlify CLI | [nodejs.org](https://nodejs.org) |
| `netlify-cli` | Deploy + link | `npm install -g netlify-cli` then `netlify login` |

### Steps

```bash
# 1. Clone into ~/.gocode
git clone https://github.com/acemaster-gh/gocode.git ~/.gocode

# 2. Run installer
cd ~/.gocode && bash install.sh

# 3. Reload shell
source ~/.bashrc

# 4. Run
gocode
```

The installer:
- Makes all scripts executable
- Symlinks `gocode` into `~/.scripts/`
- Adds `~/.scripts` to PATH
- Auto-installs `jq` and `fzf` via apt
- Initializes `~/.devsession/` session registry

---

## Usage

```bash
gocode               # Start or resume a project
gocode --push        # Commit and push with type selector
gocode --open        # Reopen last active project instantly
gocode --delete      # Delete a project everywhere
gocode --deploy      # Redeploy last project to Netlify
gocode --status      # Dashboard + open Obsidian daily note
gocode --complete    # Mark a project as completed
gocode --config      # Re-run setup wizard
gocode --update      # Pull latest version from GitHub
gocode --help        # Show all commands
```

---

## What Gets Created Per Project

```
project-name/
├── index.html      # HTML5, dark mode data-theme, theme toggle button
├── style.css       # Full :root tokens, randomized palette, dark mode
├── script.js       # const DOM={}, state, functions, event listeners, init()
└── README.md       # Pre-filled with repo URL, Netlify URL, progress log
```

**Every project has a different color palette** — generated from a random HSL hue with consistent saturation and lightness so it always looks professional. Dark mode works out of the box.

---

## Boilerplate Preview

**style.css** includes:
- Full `:root` design tokens — colors, spacing, typography, radius, transitions
- Dark mode via `[data-theme="dark"]` — unique per project
- Card, button, input, utility classes ready to use

**script.js** structure:
```js
// ── DOM Objects ──
const DOM = { app: document.getElementById('app'), ... };

// ── State ──
const state = { theme: localStorage.getItem('theme') || 'light' };

// ── Functions ──
function init() { ... }

// ── Event Listeners ──
DOM.themeToggle.addEventListener('click', toggleTheme);

// ── Init ──
init();
```

---

## Configuration

First run triggers a setup wizard that auto-detects your GitHub username and Chrome binary. Config is saved to `~/.devsession/config.sh` — edit it anytime:

```bash
GH_USER="your-github-username"
PROJECTS_DIR="$HOME/Desktop/js-projects"
CHROME_BIN="google-chrome"
CHROME_PROFILE="Default"
SLEEP_BEFORE_CHROME=10
OBSIDIAN_VAULT="$HOME/Desktop/DevMaster-Vault"

BROWSER_TABS=(
    "https://chatgpt.com"
    "https://claude.ai"
    "https://developer.mozilla.org"
    "http://localhost:5500"
)
```

Re-run wizard anytime: `gocode --config`

---

## Roadmap

### v1.0.3
- Vite-style terminal UI using `@clack/prompts` — replaces fzf
- GNOME Pomodoro session timer integration
- `gocode --note` — quick Obsidian note without opening full gocode
- zsh support

### v1.1.0
- Session timer — log coding time to Obsidian daily note
- Weekly Obsidian summary auto-generated
- `gocode --stats` — total projects, commits, time per project
- Streak tracker — consecutive coding days

### v2.0.0
- Multi-machine sync via private GitHub gist
- Snippet injector from local code snippets folder
- Template system — save any project as a reusable template
- `gocode --doctor` — full system health check

---

## Contributing

Keep the modular architecture — one concern per lib file. Don't merge unrelated logic.

```bash
# Fork → clone → branch
git checkout -b feat/your-feature

# Edit the relevant lib file
# Test on a clean install if possible

# Commit with conventional commits
git commit -m "feat: description"
git push origin feat/your-feature
# Open a pull request
```

### Guidelines
- One lib file per concern — don't mix git logic into browser.sh
- All output via `log_ok / log_warn / log_error / log_step`
- New commands go in main `gocode` entrypoint as a `case` branch
- New lib files must be sourced in main `gocode` and added to `install.sh` chmod
- Update this README when adding commands or config options
- Conventional commits: `feat:` `fix:` `docs:` `refactor:` `chore:`

### Good First Contributions
- Arrow key menus for more selectors
- zsh support alongside bash
- Install script for Arch Linux / Fedora
- Test script that validates gocode on a fresh install
- New boilerplate type (API project, dashboard, landing page)

---

## Version History

| Version | Highlights |
|---------|-----------|
| v1.0.2 | fzf arrow key menus, random color palettes, `--push` `--delete` `--open` `--deploy`, React+Vite, Netlify link verification, spinner, summary box |
| v1.0.1 | Modular architecture, Tailwind, dark mode boilerplates, GitHub search, Obsidian integration, plugin system |
| v1.0.0 | Single-file — project creation, GitHub, Netlify, session tracking |

---

## Author

Built by **Ace** — developer building toward freelance independence through daily JavaScript practice and obsessive workflow optimization.

- GitHub: [@acemaster-gh](https://github.com/acemaster-gh)

---

## License

MIT — use it, fork it, build on it.
