# gocode ⚡

> A modular, single-command developer workflow system for Linux — built to eliminate repetitive overhead when starting, managing, and tracking JavaScript projects.

---

## What is gocode?

Every time you start a new project you do the same 15 things manually — create a folder, set up files, initialize git, create a GitHub repo, open VS Code, open browser tabs, deploy to Netlify. Every single time.

`gocode` automates all of it into one command. One word in your terminal and your entire coding environment is ready in under 30 seconds.

It's not a framework. It's not a package. It's a modular bash system that lives in `~/.gocode/` and works as a lightweight personal developer operating system — purpose built for developers who want to spend time writing code, not setting up code.

---

## Architecture

```
~/.gocode/
├── gocode                  # Main entry point — thin router
├── install.sh              # One-time installer
├── lib/
│   ├── config.sh           # First-run wizard + user settings
│   ├── ui.sh               # Colors, banner, menus, logging
│   ├── deps.sh             # Dependency checker
│   ├── registry.sh         # projects.json read/write
│   ├── git.sh              # Git, GitHub, similar project search
│   ├── browser.sh          # Chrome tab launcher
│   ├── netlify.sh          # Deploy + auto-link to GitHub
│   └── obsidian.sh         # Daily note + project note writer
├── boilerplates/
│   ├── vanilla/            # HTML + CSS (full :root vars + dark mode) + JS
│   ├── tailwind/           # Tailwind CDN + dark mode toggle
│   └── react/              # Coming in v1.0.2
└── plugins/                # Drop a .sh file here — auto-loaded on run
```

---

## Features

### Project Creation
- Validates name, blocks duplicates
- Project type selector: Vanilla JS or Tailwind CSS
- Injects production-ready boilerplate — CSS reset, full `:root` design tokens, dark mode toggle wired into HTML + CSS + JS
- Git init → commit → GitHub repo → push — fully automated
- Netlify deploy + auto-link to GitHub (every `git push` redeploys)

### Session Management
- Tracks all projects in `~/.devsession/projects.json`
- On launch: shows active projects — resume any with a keypress
- Mark projects complete from the menu
- Reopens VS Code + browser tabs on resume

### GitHub Similar Projects Search
- On new project create, searches GitHub for top JS repos matching your project name
- Shows stars + repo URL in terminal for inspiration and reference
- 5 second timeout — never blocks you

### Browser Automation
- Opens ChatGPT, Claude, MDN, localhost:5500 in your Chrome profile
- 10 second delay after VS Code opens

### Obsidian Integration
- Daily note in `05-DAILY-NOTES/YYYY-MM-DD.md` — appended on every create/resume
- Individual project note in `01-PROJECTS/Active/` with repo URL, Netlify URL, task list, session log
- Notes marked complete when project is marked done

### Plugin System
- Drop any `.sh` file into `plugins/` — gocode sources it automatically on next run
- Add features without touching core files

---

## Installation

### Requirements
- Ubuntu Linux (or any Debian-based distro)
- `git` — version control
- `gh` — GitHub CLI (`gh auth login` must be done)
- `code` — VS Code
- `google-chrome` — with a Default profile
- `jq` — JSON processor (auto-installed if missing)
- `node` + `npm` — for Netlify CLI
- `netlify-cli` — `npm install -g netlify-cli` then `netlify login`

### Steps

```bash
# 1. Clone the repo
git clone https://github.com/acemaster-gh/gocode.git ~/.gocode

# 2. Run installer
cd ~/.gocode && bash install.sh

# 3. Reload shell
source ~/.bashrc

# 4. Run
gocode
```

---

## Usage

```bash
gocode               # Start or resume a project
gocode --status      # Dashboard + open Obsidian daily note
gocode --complete    # Mark a project as completed
gocode --config      # Re-run setup wizard
gocode --update      # Pull latest version from GitHub
gocode --help        # Show all commands
```

---

## Boilerplate — What Gets Injected

**Vanilla:**
- `index.html` — HTML5, dark mode `data-theme` attribute, theme toggle button
- `style.css` — Full `:root` design tokens (colors, spacing, typography, radius, transitions), dark mode overrides, card, button, input, utility classes
- `script.js` — `const DOM = {}`, `state`, `functions`, `event listeners`, `init()` — all sectioned with comments. Dark mode toggle wired and working out of the box.

**Tailwind:**
- `index.html` — Tailwind CDN, dark mode class strategy, theme toggle button
- `style.css` — Minimal extras only
- `script.js` — Same structure as vanilla, dark mode uses Tailwind `dark` class

---

## Roadmap

### v1.0.2
- Arrow key navigation in menus (like Vite CLI)
- React + Vite scaffolding
- Improved terminal UI

### v1.1.0
- Session timer — logs coding time to Obsidian
- `gocode --open` — reopen last project instantly
- Auto git push on project complete
- Weekly Obsidian summary auto-generated

### v2.0.0
- Multi-machine sync via private GitHub gist
- Netlify environment variable management
- Template system — define custom project structures
- Snippet injector from local code snippets folder

---

## Contributing

Contributions are welcome. Keep everything modular — one concern per lib file. No giant scripts.

```bash
# Fork → clone → branch
git checkout -b feat/your-feature

# Edit the relevant lib file only
# Test on a clean install if possible

# Commit
git commit -m "feat: description"
git push origin feat/your-feature
# Open a pull request
```

### Guidelines
- One lib file per concern — don't mix git logic into browser.sh etc
- All display output via `log_ok / log_warn / log_error / log_step`
- New commands go in main `gocode` entrypoint as a `case` branch
- Update this README if you add a command or config option
- Conventional commits: `feat:` `fix:` `docs:` `refactor:`

### Good first contributions
- Arrow key menu navigation
- zsh support
- Test script for fresh install validation
- New boilerplate type

---

## Author

Built by **Ace** — developer building toward freelance independence through daily JavaScript practice and obsessive workflow optimization.

- GitHub: [@acemaster-gh](https://github.com/acemaster-gh)

---

## Version History

| Version | Notes |
|---------|-------|
| v1.0.1  | Modular architecture, Tailwind support, dark mode boilerplates, GitHub search, Netlify auto-link, Obsidian integration, plugin system |
| v1.0.0  | Single-file system — project creation, GitHub, Netlify, session tracking |

---

## License

MIT — use it, fork it, build on it.
