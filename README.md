<div align="center">

# ⚡ gocode

**Terminal-based developer workflow system for Linux**

A single command that scaffolds, manages, and deploys web projects — git, GitHub, Netlify, VS Code, and your browser all automated from one interactive menu.

![Version](https://img.shields.io/badge/version-1.0.3-6366f1?style=flat-square)
![Platform](https://img.shields.io/badge/platform-Linux-22d3ee?style=flat-square)
![Shell](https://img.shields.io/badge/shell-bash-114a8b?style=flat-square)
![UI](https://img.shields.io/badge/UI-@clack%2Fprompts-green?style=flat-square)

</div>

---

## What it does

Type `gocode` and get a full interactive menu. Pick **New Project** and it:

- Asks 5 questions with a beautiful @clack UI (same library Vite uses, but gocode's own style)
- Creates the project folder with a full boilerplate
- Generates a unique HSL colour palette and injects it into CSS
- Runs `git init` and makes the first commit on `main`
- Creates a GitHub repo (public or private) — optional
- Deploys to Netlify — optional
- Generates a README with description, live link, built-with, and folder preview
- Logs everything to your Obsidian vault
- Opens VS Code and your dev browser tabs

All in under 60 seconds.

---

## Installation

```bash
git clone https://github.com/acemaster-gh/gocode ~/.gocode
cd ~/.gocode
bash install.sh
source ~/.bashrc
gocode --config
```

`install.sh` handles everything — installs `git`, `jq`, `fzf`, `gh`, `node`, `netlify-cli`, runs `gh auth login`, `netlify login`, links `gocode` to your PATH, and installs `@clack/prompts`.

---

## Usage

```bash
gocode          # opens the main interactive menu
gocode --new    # jump straight to project creation
gocode --help   # list all commands
```

### Main menu

Running `gocode` with no arguments shows:

```
  ██████╗  ██████╗  ██████╗ ██████╗ ██████╗ ███████╗
 ██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝
 ██║  ███╗██║   ██║██║     ██║   ██║██║  ██║█████╗
 ██║   ██║██║   ██║██║     ██║   ██║██║  ██║██╔══╝
 ╚██████╔╝╚██████╔╝╚██████╗╚██████╔╝██████╔╝███████╗
  ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝

◆  What do you want to do?
│  ⚡ New Project      scaffold, push, deploy a new project
│  ▶ Resume           pick up where you left off
│  ↑ Push             commit and push to GitHub
│  🚀 Deploy          push to Netlify
│  ▦ All Projects     see every project and status
│  ◉ Status           detailed info on a project
│  ✎ Edit             rename · description · visibility
│  ✕ Delete           remove a project everywhere
│  📝 Quick Note      add a note to Obsidian
│  ⚙ Settings         config · browser · doctor · backup
└
```

---

## Commands

| Command | What it does |
|---------|-------------|
| `gocode` | Main interactive menu |
| `--new` | Create a new project |
| `--push` | Select commit type + message, push |
| `--deploy` | Deploy to Netlify |
| `--status` | Show project details |
| `--list` | All projects with status table |
| `--log` | Last 10 commits |
| `--diff` | Uncommitted changes |
| `--edit` | Edit name / description / visibility everywhere |
| `--delete` | Delete folder, GitHub repo, registry, Obsidian note |
| `--rename` | Rename project in all places |
| `--visibility` | Toggle public / private on GitHub |
| `--desc` | Set GitHub repo description |
| `--complete` | Mark project complete, update README |
| `--archive` | Archive a project |
| `--note` | Quick note to Obsidian |
| `--open-github` | Open GitHub repo in browser |
| `--open-live` | Open live Netlify site |
| `--open-netlify` | Open Netlify dashboard |
| `--open-folder` | Open project in file manager |
| `--stats` | Project counts by status |
| `--backup` | Back up registry and config |
| `--doctor` | Check all dependencies and registry health |
| `--set-browser` | Redetect and pick browser |
| `--config` | Re-run setup wizard |
| `--update` | Pull latest gocode |
| `--help` | Show all commands |

---

## Project types

**Vanilla JS** — HTML + CSS with full CSS variable system, dark mode support, colour palette auto-generated and injected per project.

**Tailwind CSS** — Same structure with Tailwind CDN. No build step needed.

**React + Vite** — Full `npm create vite@latest` scaffold with dependencies installed and ready to `npm run dev`.

---

## UI

gocode uses `@clack/prompts` for the interactive layer — the same library powering Vite, create-remix, and create-astro — but with gocode's own identity:

- Full ASCII art banner responsive to terminal width
- Step counter `[1/5]` on every prompt
- Folder structure preview before building
- `p.note()` summary showing exactly what will be created
- Final outro with next steps, GitHub link, and live URL
- Colour-coded commit types in `--push`
- `p.multiselect()` for browser tab configuration

---

## How the registry works

Every project is saved to `~/.devsession/projects.json`:

```json
{
  "projects": [
    {
      "name": "modal-system",
      "type": "vanilla",
      "path": "/home/ace/Desktop/js-projects/modal-system",
      "repo_url": "https://github.com/acemaster-gh/modal-system",
      "netlify_url": "https://modal-system.netlify.app",
      "description": "A modal system with smooth animations",
      "status": "active",
      "created": "2026-05-17",
      "last_opened": "2026-05-24"
    }
  ]
}
```

Registry writes are atomic — written to a temp file then moved — so a crash never corrupts your data.

---

## Obsidian integration

If you have an Obsidian vault configured, gocode automatically creates a note per project at `vault/gocode-projects/project-name.md`, logs every create and resume to your daily note, and marks the note complete when you run `--complete`.

---

## File structure

```
~/.gocode/
├── gocode              # main entry point
├── install.sh          # one-shot installer
├── package.json        # @clack/prompts dependency
├── lib/
│   ├── prompt.js       # @clack/prompts UI layer
│   ├── ui.sh           # ASCII banner, colors, spinner, log functions
│   ├── config.sh       # setup wizard, browser detection
│   ├── registry.sh     # projects.json CRUD, atomic writes
│   ├── git.sh          # git and GitHub operations
│   ├── netlify.sh      # Netlify deploy and site management
│   ├── browser.sh      # multi-browser detection
│   ├── obsidian.sh     # vault integration
│   ├── colors.sh       # HSL palette generator
│   ├── readme_gen.sh   # project README generator
│   └── deps.sh         # dependency checker
├── boilerplates/
│   ├── vanilla/
│   ├── tailwind/
│   └── react/
└── plugins/            # drop .sh files here to extend gocode
```

---

## Changelog

### v1.0.3
- Full `@clack/prompts` UI — step counters, folder preview, summary note, next-steps outro
- Main interactive menu on `gocode` with no arguments
- Optional GitHub repo creation with public/private choice
- Optional Netlify deployment
- Browser auto-detection — Chrome, Brave, Firefox, Edge, Chromium, Opera, Vivaldi
- `--delete` — removes folder, GitHub repo, registry entry, Obsidian note
- `--edit` — rename/description/visibility updated in all places at once
- `--list` — all projects in a status table
- `--note` — quick Obsidian note from terminal
- `--open-folder`, `--open-github`, `--open-live`, `--open-netlify`
- Atomic registry writes — no corruption on crash
- `set -euo pipefail` throughout
- ASCII art banner responsive to terminal width
- HSL colour palette generator — unique colours per project
- `github_search_similar` — shows related repos on project creation

### v1.0.2
- Initial public release
- Vanilla, Tailwind, React project types
- GitHub and Netlify integration
- Obsidian vault logging
- fzf-based menus

---

## Roadmap (v1.0.4)

Branch management · PR creation · GitHub issues · Netlify rollback · env vars · project templates · import existing projects · Pomodoro timer · time tracking · changelog generation · GitHub releases · portfolio generator · AI commit messages · terminal dashboard

---

## License

MIT

---

<div align="center">
Built by <a href="https://github.com/acemaster-gh">acemaster-gh</a>
</div>
