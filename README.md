# gocode ⚡

> A personal developer workflow automation system for Linux — built to eliminate the repetitive overhead of starting, managing, and tracking JavaScript projects.

---

## What is gocode?

Every time you start a new project you do the same 15 things manually — create a folder, set up files, initialize git, create a GitHub repo, open VS Code, open browser tabs, set up Netlify. Every single time.

`gocode` automates all of it into a single command. One word in your terminal and your entire coding environment is ready in under 30 seconds.

It's not a framework. It's not a package. It's a shell script that lives in `~/.scripts/` and works as a lightweight personal developer operating system — purpose built for developers who want to spend time writing code, not setting up code.

---

## Features

### Project Creation
- Creates project folder with auto-validated, duplicate-protected naming
- Injects production-ready boilerplate into `index.html`, `style.css`, `script.js`, `README.md`
- CSS reset + `min-height: 100vh` centering + `const DOM = {}` pattern out of the box
- Git init → initial commit → GitHub repo creation → push — fully automated
- Netlify deploy on first run + links Netlify to GitHub for auto-deploy on every future push

### Session Management
- Tracks all projects in `~/.devsession/projects.json`
- On launch shows all active projects — resume any with a single keypress
- Mark projects complete from inside the menu
- Reopens VS Code workspace and browser tabs automatically on resume

### Browser Automation
- Opens ChatGPT, Claude, MDN, and localhost:5500 in your Chrome profile after VS Code loads
- 10 second delay so VS Code is ready before browser opens

### Obsidian Integration
- Writes a daily note to `05-DAILY-NOTES/YYYY-MM-DD.md` on every project create or resume
- Creates an individual project note in `01-PROJECTS/Active/` with repo URL, Netlify URL, goals, and progress log
- Marks notes as completed when you mark a project done
- `gocode --status` opens today's Obsidian daily note automatically

### Dependency System
- Checks for `git`, `gh`, `code`, `google-chrome`, `jq` on every run
- Auto-installs `jq` if missing
- Silent Netlify check — warns but never blocks

---

## Installation

### Requirements
- Ubuntu Linux (or any Debian-based distro)
- `git` — version control
- `gh` — GitHub CLI, authenticated via `gh auth login`
- `code` — VS Code with Live Server extension
- `google-chrome` — with a Default profile
- `jq` — JSON processor (auto-installed if missing)
- `netlify-cli` — install via `npm install -g netlify-cli` then `netlify login`
- `node` + `npm` — for Netlify CLI

### Steps

```bash
# 1. Clone the repo
git clone https://github.com/acemaster-gh/gocode.git
cd gocode

# 2. Run the installer
bash install.sh

# 3. Reload your shell
source ~/.bashrc

# 4. Run
gocode
```

The installer:
- Copies `gocode` to `~/.scripts/gocode`
- Makes it executable
- Adds `~/.scripts` to your PATH in `~/.bashrc`
- Initializes `~/.devsession/` session directory

---

## Usage

```bash
gocode              # Start or resume a project
gocode --status     # Project dashboard + open Obsidian daily note
gocode --complete   # Mark an active project as completed
gocode --help       # Show all commands
```

### First run — new project
```
[?] Project name (e.g. weather-app): tip-calculator
[....] Creating folder structure...
[OK]   Starter files created.
[....] Initializing git...
[OK]   Git initialized and committed.
[....] Creating GitHub repo...
[OK]   GitHub repo created: https://github.com/you/tip-calculator
[....] Deploying to Netlify...
[OK]   Netlify deployed: https://tip-calculator.netlify.app
[OK]   Netlify linked to GitHub — auto-deploy on every push enabled.
[....] Writing Obsidian notes...
[OK]   Obsidian daily note updated: 2026-05-09.md
[OK]   Obsidian project note created: tip-calculator.md
[DONE] Project 'tip-calculator' is ready.
[....] Opening VS Code...
[....] Opening browser tabs in 10s...
```

### Second run — active projects exist
```
────────────────────────────────────────────────────
Active Projects:
────────────────────────────────────────────────────
  1. tip-calculator (last opened: 2026-05-09)
     https://github.com/you/tip-calculator
  N. Start a new project
  C. Mark a project as completed
  S. Status dashboard
────────────────────────────────────────────────────
[?] Choice:
```

---

## Project Structure Generated

Every new project gets this structure:

```
project-name/
├── index.html      # HTML5 boilerplate with linked CSS and JS
├── style.css       # CSS reset + centered body layout
├── script.js       # const DOM = {} pattern + init function
└── README.md       # Pre-filled with repo URL and Netlify link
```

---

## Obsidian Vault Structure Expected

```
DevMaster-Vault/
├── 01-PROJECTS/
│   └── Active/           # One note per project created here
└── 05-DAILY-NOTES/       # One note per day, appended on each session
```

If your vault is in a different location, update `OBSIDIAN_VAULT` at the top of the `gocode` script.

---

## Configuration

All config lives at the top of `~/.scripts/gocode`:

```bash
PROJECTS_DIR="$HOME/Desktop/js-projects"     # Where projects are created
OBSIDIAN_VAULT="$HOME/Desktop/DevMaster-Vault" # Your Obsidian vault
CHROME_PROFILE="Default"                      # Chrome profile name
SLEEP_BEFORE_CHROME=10                        # Seconds to wait before opening browser
```

Change any of these to match your setup. The GitHub username is pulled automatically from `gh api user` — no hardcoding needed.

---

## Roadmap

These are features actively planned for future versions:

### Near Term
- **First-run config wizard** — detects your username, vault path, projects dir, and Chrome profile automatically. Makes gocode portable across any Linux machine with zero manual edits.
- **Project type selector** — choose between `vanilla`, `tailwind`, or `react` on project create. Each injects the correct boilerplate and build setup.
- **Tailwind support** — CDN-based for quick projects, Vite-based for production builds.
- **React + Vite support** — full scaffold with `npm create vite`, folder structure, and Netlify build settings auto-configured.

### Medium Term
- **GitHub similar projects search** — on project create, pulls top repos matching your project type from GitHub API and displays them in terminal as inspiration and reference.
- **Session timer** — tracks how long VS Code is open per session and logs it to your Obsidian daily note.
- **`gocode --open`** — reopen last active project instantly without going through the menu.
- **Auto git push on complete** — when you mark a project done, runs a final `git add . && git commit && git push` automatically.

### Long Term
- **Framework-aware Netlify config** — auto-generates `netlify.toml` with correct build commands per project type.
- **Snippet injector** — pick from your local code snippets folder and inject into the current project.
- **Multi-machine sync** — sync session state and project registry across machines via a private GitHub gist.
- **Obsidian weekly review auto-fill** — at the end of each week, generates a summary of projects worked on, commits made, and time logged.
- **Template system** — define custom project templates (landing page, dashboard, API-connected app) and scaffold them on create.

---

## Contributing

Contributions are welcome. This project is intentionally kept as a single-file shell script — please keep it that way. No external dependencies beyond what's listed in requirements.

### How to contribute

```bash
# Fork the repo, clone your fork
git clone https://github.com/your-username/gocode.git
cd gocode

# Create a branch
git checkout -b feat/your-feature-name

# Make your changes to the gocode script
# Test on a fresh Linux environment if possible

# Commit using conventional commits
git add .
git commit -m "feat: description of what you added"
git push origin feat/your-feature-name

# Open a pull request on GitHub
```

### Contribution guidelines
- Keep everything in the single `gocode` file — no splitting into lib files
- New features go in as clearly named, clearly commented functions
- Follow the existing section structure (UTILITIES → DEPENDENCIES → REGISTRY → etc)
- Test that existing features still work after your change
- Update this README if you add a new command or config option
- Use conventional commit messages: `feat:`, `fix:`, `docs:`, `refactor:`

### Good first contributions
- Adding support for a new Linux distro's Chrome binary name
- Adding `zsh` support alongside bash
- Writing a test script that validates gocode works on a fresh install
- Improving error messages

---

## Author

Built by **Ace** — a developer building toward freelance independence through daily JavaScript practice and obsessive workflow optimization.

- GitHub: [@acemaster-gh](https://github.com/acemaster-gh)

---

## License

MIT — use it, fork it, build on it.
