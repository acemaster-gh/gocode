package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/acemaster-gh/gocode/internal/colors"
	"github.com/acemaster-gh/gocode/internal/config"
	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	ghpkg "github.com/acemaster-gh/gocode/internal/github"
	"github.com/acemaster-gh/gocode/internal/netlify"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/plugin"
	"github.com/acemaster-gh/gocode/internal/registry"
	tmpl "github.com/acemaster-gh/gocode/internal/template"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// Rollback state — tracks what was created so we can clean up on failure.
type rollbackState struct {
	projectPath string
	repoURL     string
	ghUser      string
}

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create"},
	Short:   "Create a new project",
	Long:    "Scaffold a new project, create GitHub repo, deploy to Netlify",
	RunE:    runNew,
}

func runNew(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}

	ui.Banner(Version)

	// ── Step 1-5: Collect project info via huh form ───────────────────────────
	var (
		name         string
		description  string
		stack        string
		createGit    bool
		visibility   string
		deployNow    bool
	)

	err = huh.NewForm(
		// Group 1: Name + description
		huh.NewGroup(
			huh.NewInput().
				Title("[1/5] Project name").
				Placeholder("tip-calculator").
				Value(&name).
				Validate(func(s string) error {
					s = strings.TrimSpace(s)
					if s == "" {
						return errors.New("name cannot be empty")
					}
					if strings.ContainsAny(s, " /\\:*?\"<>|") {
						return errors.New("no spaces or special characters")
					}
					return nil
				}),
			huh.NewInput().
				Title("[2/5] What does it do?").
				Placeholder("a short description").
				Value(&description),
		),
		// Group 2: Stack
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("[3/5] Technology stack").
				Options(
					huh.NewOption("Vanilla JS      HTML + CSS + JS, CSS variable system", "vanilla"),
					huh.NewOption("Tailwind CSS    utility-first, v4 CLI pipeline", "tailwind"),
					huh.NewOption("React + Vite    SPA with hot-reload", "react"),
				).
				Value(&stack),
		),
		// Group 3: GitHub
		huh.NewGroup(
			huh.NewConfirm().
				Title("[4/5] Create GitHub repository?").
				Value(&createGit),
			huh.NewSelect[string]().
				Title("Repository visibility").
				Options(
					huh.NewOption("🌐 Public", "public"),
					huh.NewOption("🔒 Private", "private"),
				).
				Value(&visibility),
		),
		// Group 4: Deploy
		huh.NewGroup(
			huh.NewConfirm().
				Title("[5/5] Deploy to Netlify?").
				Value(&deployNow),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(ui.Dim.Render("cancelled."))
			return nil
		}
		return err
	}

	// Clean name
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "-")
	description = strings.TrimSpace(description)

	// Check for duplicate
	reg := getRegistry()
	if reg.Exists(name) {
		ui.Error(fmt.Sprintf("Project '%s' already exists in registry.", name))
		return nil
	}

	projectPath := filepath.Join(cfg.ProjectsDir, name)
	if _, err := os.Stat(projectPath); err == nil {
		ui.Error(fmt.Sprintf("Directory already exists: %s", projectPath))
		return nil
	}

	rb := &rollbackState{}

	// ── Create directory ──────────────────────────────────────────────────────
	ui.Step(fmt.Sprintf("Creating project: %s", name))
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	rb.projectPath = projectPath

	// ── Scaffold boilerplate ──────────────────────────────────────────────────
	sp := ui.StartSpinner(fmt.Sprintf("Scaffolding %s...", stack))
	scaffoldErr := tmpl.Scaffold(stack, projectPath, name)
	sp.Stop()
	if scaffoldErr != nil {
		doRollback(rb)
		return fmt.Errorf("scaffold: %w", scaffoldErr)
	}
	ui.OK(fmt.Sprintf("%s boilerplate ready.", stack))

	// ── Color palette (vanilla only) ──────────────────────────────────────────
	if stack == "vanilla" {
		cssPath := filepath.Join(projectPath, "style.css")
		if _, err := os.Stat(cssPath); err == nil {
			palette := colors.Generate()
			if err := palette.InjectIntoCSS(cssPath); err == nil {
				ui.OK(fmt.Sprintf("Color palette injected (hue: %d°)", palette.Hue))
			}
		}
	}

	// ── Git init ──────────────────────────────────────────────────────────────
	ui.Step("Initialising git...")
	if err := gitpkg.InitAndCommit(projectPath); err != nil {
		ui.Warn(fmt.Sprintf("git init failed: %v", err))
	} else {
		ui.OK("Git initialised (branch: main)")
	}

	// ── GitHub search (background) ────────────────────────────────────────────
	go func() {
		results, err := ghpkg.SearchSimilar(name)
		if err != nil || len(results) == 0 {
			return
		}
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, ui.Bold.Render("  Similar projects on GitHub:"))
		ui.Divider()
		for _, r := range results {
			fmt.Fprintf(os.Stderr, "  %s  %s\n",
				ui.Cyan.Render(fmt.Sprintf("⭐ %d", r.StargazersCount)),
				ui.Bold.Render(r.FullName))
			fmt.Fprintf(os.Stderr, "       %s\n", ui.Dim.Render(r.URL))
		}
		ui.Divider()
	}()

	// ── GitHub repo creation ──────────────────────────────────────────────────
	repoURL := "local-only"
	if createGit {
		sp = ui.StartSpinner(fmt.Sprintf("Creating GitHub repo (%s)...", visibility))
		url, err := ghpkg.CreateRepo(name, visibility)
		sp.Stop()
		if err != nil {
			ui.Warn(fmt.Sprintf("GitHub: %v", err))
		} else {
			repoURL = url
			rb.repoURL = repoURL
			rb.ghUser = cfg.GHUser
			ui.OK(fmt.Sprintf("Repo created: %s", repoURL))

			// Push
			sp = ui.StartSpinner("Pushing to GitHub...")
			if err := gitpkg.AddRemoteAndPush(projectPath, repoURL); err != nil {
				sp.Stop()
				ui.Warn("Push failed — run: gocode push")
			} else {
				sp.Stop()
				ui.OK("Pushed to GitHub.")
			}
		}
	}

	// ── Netlify deploy ────────────────────────────────────────────────────────
	netlifyURL := "not-deployed"
	if deployNow {
		sp = ui.StartSpinner("Deploying to Netlify...")
		netlifyURL = netlify.Full(name, projectPath, repoURL)
		sp.Stop()
		if netlifyURL != "not-deployed" {
			ui.OK(fmt.Sprintf("Deployed: %s", netlifyURL))
		} else {
			ui.Warn("Netlify deploy failed — run: gocode deploy")
		}
	}

	// ── README ────────────────────────────────────────────────────────────────
	if err := tmpl.GenerateREADME(name, stack, projectPath,
		repoURL, netlifyURL, description); err == nil {
		ui.OK("README.md generated.")
	}

	// Push README update
	gitpkg.PushFinal(projectPath, repoURL)

	// ── Registry ──────────────────────────────────────────────────────────────
	if err := reg.Add(registry.Project{
		Name:        name,
		Path:        projectPath,
		Type:        stack,
		RepoURL:     repoURL,
		NetlifyURL:  netlifyURL,
		Description: description,
	}); err != nil {
		ui.Warn(fmt.Sprintf("Registry: %v", err))
	}

	// ── Obsidian (background) ─────────────────────────────────────────────────
	vault := obsidian.New(cfg.ObsidianVault)
	go func() {
		vault.ProjectNote(name, stack, repoURL, netlifyURL, description)
		vault.DailyLog("Created: " + name)
	}()

	// ── Plugin hooks (background) ─────────────────────────────────────────────
	go plugin.RunHooks(config.PluginsDir(), plugin.HookPostCreate, plugin.Payload{
		Hook:    plugin.HookPostCreate,
		Project: name,
		Data: map[string]string{
			"path":    projectPath,
			"stack":   stack,
			"repo":    repoURL,
			"netlify": netlifyURL,
		},
	})

	// ── Open VS Code ──────────────────────────────────────────────────────────
	go func() {
		exec.Command("code", projectPath).Start()
	}()

	// ── Open browser tabs ─────────────────────────────────────────────────────
	go func() {
		time.Sleep(time.Duration(cfg.SleepBeforeBrowser) * time.Second)
		tabs := cfg.BrowserTabs
		if netlifyURL != "not-deployed" {
			tabs = append(tabs, netlifyURL)
		}
		for _, tab := range tabs {
			openURL(tab)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// ── Summary box ───────────────────────────────────────────────────────────
	ui.SummaryBox(name, stack, projectPath, repoURL, netlifyURL)

	fmt.Printf("  %s\n\n",
		ui.Dim.Render("Next: gocode push  ·  gocode deploy  ·  gocode todo"))

	return nil
}

// ── Rollback ──────────────────────────────────────────────────────────────────

func doRollback(rb *rollbackState) {
	if rb.projectPath != "" {
		ui.Warn("Rolling back: removing project directory...")
		os.RemoveAll(rb.projectPath)
	}
	if rb.repoURL != "" && rb.repoURL != "local-only" && rb.ghUser != "" {
		// Extract repo name from URL
		parts := strings.Split(rb.repoURL, "/")
		if len(parts) > 0 {
			repoName := parts[len(parts)-1]
			ui.Warn("Rolling back: deleting GitHub repo...")
			ghpkg.DeleteRepo(rb.ghUser, repoName)
		}
	}
}
