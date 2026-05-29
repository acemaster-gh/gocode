<div align="center">
<br>

```
  ██████╗  ██████╗  ██████╗ ██████╗ ██████╗ ███████╗
  ██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝
  ██║  ███╗██║   ██║██║     ██║   ██║██║  ██║█████╗
  ██║   ██║██║   ██║██║     ██║   ██║██║  ██║██╔══╝
  ╚██████╔╝╚██████╔╝╚██████╗╚██████╔╝██████╔╝███████╗
   ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝
```

<h3>Developer workflow automation.<br>One binary. Every command you run every day.</h3>

<br>

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Version](https://img.shields.io/badge/Version-v2.0.0-6366f1?style=for-the-badge)](https://github.com/acemaster-gh/gocode/releases/tag/v2.0.0)
[![License](https://img.shields.io/badge/License-MIT-22c55e?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-808080?style=for-the-badge)](#install)

<br>

</div>

---

<br>

## The problem

Every new project starts the same way.

Open terminal. `mkdir project`. `cd project`. `git init`. Create repo on GitHub. Copy boilerplate from last project. Fix the file names. Set up the remote. `npm install`. Open VS Code. Open five browser tabs. Set up Netlify. Push.

Fifteen minutes gone before a single line of real code is written.

**gocode collapses that into one command.**

<br>

---

<br>

## What it looks like

```
$ gocode new

  ██████╗  ██████╗  ██████╗ ██████╗ ██████╗ ███████╗
  ...

  gocode v2.0.0  —  last: portfolio  ·  4 active

  > New Project      scaffold · push · deploy
    Resume           open last project in VS Code
    Push             commit + push to GitHub
    Deploy           push live to Netlify
    ...


  [1/5] Project name
  > tip-calculator


  [2/5] What does it do?
  > splits bills and calculates tips


  [3/5] Technology stack
  > Vanilla JS


  [4/5] Create GitHub repository?
  > Yes  ·  Public


  [5/5] Deploy to Netlify?
  > Yes


  [OK]    vanilla boilerplate ready
  [OK]    color palette injected  (hue: 214°)
  [OK]    git initialised  (branch: main)
  [OK]    repo created: https://github.com/you/tip-calculator
  [OK]    pushed to GitHub
  [OK]    deployed: https://tip-calculator.netlify.app
  [OK]    README.md generated

  ╭─────────────────────────────────────────────╮
  │  Project Ready: tip-calculator              │
  │                                             │
  │    Stack:   vanilla                         │
  │    Path:    ~/projects/tip-calculator       │
  │    GitHub:  https://github.com/you/...      │
  │    Live:    https://tip-calculator...       │
  ╰─────────────────────────────────────────────╯

  Next: gocode push  ·  gocode deploy  ·  gocode todo
```

<br>

---

<br>

## Install

**Requirements:** Go 1.22+, git

```bash
# Clone to ~/.gocode  (source lives here, binary builds here)
git clone https://github.com/acemaster-gh/gocode ~/.gocode

# Build
cd ~/.gocode && go build -ldflags "-X main.Version=2.0.0 -s -w" -o gocode .

# Add alias
echo 'alias gocode="$HOME/.gocode/gocode"' >> ~/.bashrc && source ~/.bashrc

# Verify
gocode --version
```

> Already have gocode v1? Your `~/.devsession/projects.json` is automatically detected and migrated.

<br>

---

<br>

## First run

```bash
gocode config
```

You will be asked for:

| Field | Example |
|---|---|
| GitHub username | `acemaster-gh` |
| Projects directory | `~/Desktop/js-projects` |
| Obsidian vault path | `~/Desktop/DevMaster-Vault` *(optional)* |
| Browser tabs | one URL per line |
| Browser delay | `5` seconds |

Config is saved to `~/.gocode/config.toml`. Edit it any time.

<br>

---

<br>

## Commands

<br>

### Project lifecycle

```
gocode new          Full project creation
                    scaffold → git init → GitHub repo → Netlify → VS Code → browser tabs

gocode resume       Pick a project, open in VS Code, show pending todos

gocode complete     Mark a project done, update README progress table,
                    update Obsidian status

gocode archive      Move a project to archived status

gocode delete       Tear down everything — folder, GitHub repo,
                    registry entry, Obsidian note
```

<br>

### Git

```
gocode push         Conventional commit with type picker
                    (feat / fix / docs / style / refactor / perf / test / chore)
                    then push

gocode log          Last 15 commits, pager for full history

gocode diff         Colour diff of all uncommitted changes

gocode branch       Create (feat/ fix/ refactor/ docs/ chore/),
                    switch, or delete branches

gocode stash        Save changes with a description,
                    apply or drop saved stashes
```

<br>

### Cloud

```
gocode deploy       Push to Netlify — auto-commits dirty tree first,
                    updates registry with live URL

gocode open         Open GitHub repo / live site / Netlify dashboard / local folder
                    in browser or file manager
```

<br>

### Workspace

```
gocode list         Full project table — name, stack, status, created, links

gocode status       Deep panel — branch, dirty files, todos, all URLs, dates

gocode edit         Rename (syncs folder + GitHub repo + Obsidian note)
                    Update description everywhere
                    Toggle public / private on GitHub
                    Change lifecycle status

gocode todo         Per-project task lists — add, toggle done/undone

gocode note         Quick note to Obsidian daily log or project file

gocode timer        Pomodoro — 5 / 15 / 25 / 50 min, live progress bar,
                    terminal bell on completion
```

<br>

### System

```
gocode doctor       Check every dependency — git, node, gh, netlify, VS Code —
                    show version, auth status, install hint if missing

gocode stats        Bar charts — by stack, by status, monthly activity (6 months)

gocode backup       Timestamped snapshots of registry + config
                    saved to ~/.gocode/backups/

gocode config       Setup wizard — reconfigure anything

gocode update       git pull + rebuild binary from source
```

<br>

---

<br>

## Project registry

Every project you create is tracked in `~/.gocode/data/projects.json`.

Writes are **atomic** — data is written to a temp file and renamed into place. A `Ctrl+C` mid-save never corrupts anything. On a failed `gocode new`, the partially-created folder and GitHub repo are rolled back automatically.

```jsonc
{
  "projects": [
    {
      "name": "tip-calculator",
      "path": "/home/ace/Desktop/js-projects/tip-calculator",
      "type": "vanilla",
      "status": "active",               // active · completed · archived
      "created": "2026-05-27",
      "last_opened": "2026-05-29",
      "repo_url": "https://github.com/acemaster-gh/tip-calculator",
      "netlify_url": "https://tip-calculator.netlify.app",
      "description": "splits bills and calculates tips",
      "todos": [
        { "task": "add split-bill UI", "done": false, "created": "2026-05-27" }
      ]
    }
  ]
}
```

<br>

---

<br>

## Templates

Three stacks are **compiled directly into the binary** using `//go:embed`. No internet needed to scaffold.

<br>

### Vanilla JS

```
project/
├── index.html      header · hero · card grid
├── style.css       full CSS variable system — --primary, --bg-main, --text-main ...
├── script.js       ES module entry point
└── README.md       auto-generated with live URL + progress table
```

Every new vanilla project gets a **unique color palette**. gocode generates a random HSL hue, derives a full light/dark design system from it (primary, hover, backgrounds, borders, text, success, error, warning), and injects the hex values directly into `style.css`. Pure Go — no Python, no external tools.

<br>

### Tailwind CSS

```
project/
├── src/
│   ├── index.html
│   ├── css/input.css       @import "tailwindcss" · @theme · @layer components
│   └── js/main.js
├── dist/css/style.css
└── package.json            npm run dev · npm run build
```

Tailwind v4 CSS-first config. `.btn`, `.card`, `.btn-primary`, `.btn-outline` defined in `@layer components`. Dark mode via Tailwind's dark class.

<br>

### React + Vite

Calls `npm create vite@latest` at runtime, runs `npm install` automatically. No scaffolding is embedded because npm packages cannot be compiled into a Go binary.

<br>

---

<br>

## Obsidian integration

When a vault path is set in config, gocode writes and maintains markdown files automatically. No Obsidian plugin required.

| Event | What gocode writes |
|---|---|
| `gocode new` | Creates `vault/gocode-projects/name.md` with stack, URLs, status |
| `gocode push` | Appends timestamped entry to `vault/daily/YYYY-MM-DD.md` |
| `gocode deploy` | Appends deploy entry to daily note |
| `gocode note` | Appends your text to daily note or project note |
| `gocode complete` | Updates `**Status:** active` → `**Status:** ✅ completed` |
| `gocode edit` (rename) | Renames the note file |
| `gocode delete` | Removes the note file |

<br>

---

<br>

## Plugin system

Place any executable in `~/.gocode/plugins/` and gocode calls it after lifecycle events.

**Hooks:** `post-create` · `post-push` · `post-deploy` · `post-complete`

Every hook receives a JSON payload on stdin:

```json
{
  "hook": "post-deploy",
  "project": "tip-calculator",
  "data": {
    "url": "https://tip-calculator.netlify.app",
    "path": "/home/ace/Desktop/js-projects/tip-calculator",
    "stack": "vanilla"
  }
}
```

Example — Slack notification:

```bash
#!/usr/bin/env bash
# ~/.gocode/plugins/slack-notify
PAYLOAD=$(cat)
PROJECT=$(echo "$PAYLOAD" | jq -r '.project')
URL=$(echo "$PAYLOAD"     | jq -r '.data.url')
curl -s -X POST "$SLACK_WEBHOOK" \
  -H 'Content-type: application/json' \
  -d "{\"text\": \"*$PROJECT* deployed → $URL\"}"
```

```bash
chmod +x ~/.gocode/plugins/slack-notify
```

<br>

---

<br>

## Architecture

```
~/.gocode/
│
├── gocode               ← compiled binary  (this is the whole program)
├── config.toml          ← your settings
│
├── data/
│   └── projects.json    ← registry  (atomic writes, auto-backed-up)
│
├── plugins/             ← drop executables here for lifecycle hooks
├── backups/             ← timestamped registry + config snapshots
│
└── source/
    ├── main.go                    entry point, --flag compat layer
    ├── cmd/                       one file per command (23 commands)
    │   ├── root.go                interactive menu, shared helpers
    │   ├── new.go                 project creation engine
    │   └── ...
    └── internal/
        ├── ui/        styles, banner, spinner, log functions
        ├── config/    TOML config via Viper
        ├── registry/  atomic JSON CRUD, todos, rollback
        ├── git/       all git operations via exec.Command
        ├── github/    gh CLI wrapper
        ├── netlify/   netlify CLI wrapper
        ├── obsidian/  vault note lifecycle
        ├── colors/    HSL → hex palette generation (pure Go)
        ├── template/  go:embed scaffold + README generation
        └── plugin/    hook runner
```

<br>

**Core dependencies:**

| Package | Why |
|---|---|
| `charmbracelet/huh` | All interactive forms — replaces Node.js + @clack/prompts entirely |
| `charmbracelet/lipgloss` | Terminal colour, layout, borders |
| `spf13/cobra` | Command routing |
| `spf13/viper` | TOML config |

<br>

---

<br>

## v1 → v2

gocode v1 was Bash + Node.js. v2 is a compiled Go binary. Here is what changed and why.

<br>

| | v1 | v2 |
|---|---|---|
| Runtime | bash + node + python3 + fzf | nothing — just the binary |
| Forms | Node.js V8 cold-boot on every input | native Go, instant |
| Color palette | `python3 -c "import colorsys"` | pure Go HSL math |
| Spinner | background subshell + `/dev/tty` | goroutine, no race conditions |
| Registry writes | direct overwrite | atomic temp-file + rename |
| Rollback | none | full cleanup on failure |
| Templates | files on disk at install path | compiled into binary |
| Windows | no | yes |
| Startup | ~400ms | ~8ms |

<br>

All v1 flags still work. `gocode --new`, `gocode --push`, etc. are translated to subcommands transparently so existing muscle memory and aliases are not broken.

<br>

---

<br>

## Dependencies

gocode has zero runtime dependencies. The following tools unlock optional features:

<br>

```
git          required for everything
─────────────────────────────────────────────────────────────
gh           GitHub CLI      repo create · rename · delete · visibility
netlify      Netlify CLI     deploy · redeploy
node / npm   Node.js         React+Vite scaffold · Tailwind dev server
code         VS Code         auto-open on resume and new project
```

Run `gocode doctor` to see what is installed, what version, and whether it is authenticated.

<br>

---

<br>

## Building

```bash
cd ~/.gocode

# Full build + install
make install

# Build only
make build

# Run tests
make test

# Lint
make vet

# Format
make fmt
```

<br>

---

<br>

## Roadmap

```
v2.1    test suite · golangci-lint · GitHub Actions CI
v2.2    AI commit messages via Ollama · code smell audit · natural language search
v2.3    native Vercel API · native Cloudflare Pages API
v2.4    SQLite registry backend · cross-machine sync
v2.5    WASM plugin ecosystem · plugin registry
```

<br>

---

<br>

## Contributing

```bash
# Fork → branch → commit → PR

git checkout -b feat/your-feature

# Conventional commits only
git commit -m "feat: description"
git commit -m "fix: description"
git commit -m "docs: description"

# Before pushing
make fmt vet test
```

<br>

---

<br>

## License

MIT — see [LICENSE](LICENSE)

<br>

---

<br>

<div align="center">

Built by [acemaster-gh](https://github.com/acemaster-gh)

Powered by [Charm](https://charm.sh) — the best terminal UI toolkit in existence

</div>
