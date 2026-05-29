package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/acemaster-gh/gocode/internal/config"
	"github.com/acemaster-gh/gocode/internal/registry"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var Version string

var rootCmd = &cobra.Command{
	Use:   "gocode",
	Short: "Developer workflow system",
	Long:  "gocode v2 — terminal-based developer workflow automation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMenu()
	},
}

func Execute(version string) {
	Version = version
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		newCmd, pushCmd, resumeCmd, deployCmd,
		statusCmd, listCmd, editCmd, deleteCmd,
		logCmd, diffCmd, branchCmd, stashCmd,
		todoCmd, timerCmd, noteCmd,
		completeCmd, archiveCmd,
		doctorCmd, backupCmd, statsCmd,
		configCmd, updateCmd, openCmd,
	)
}

// ── Main Menu ─────────────────────────────────────────────────────────────────

func runMenu() error {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			ui.Warn("No config found — running setup wizard...")
			return runConfigWizard()
		}
		return err
	}
	_ = cfg

	ui.Banner(Version)

	reg := getRegistry()
	active, _ := reg.Active()
	lastProject := reg.LastActive()

	hint := fmt.Sprintf("%d active project(s)", len(active))
	if lastProject != "" {
		hint = fmt.Sprintf("last: %s  ·  %d active", lastProject, len(active))
	}

	var action string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("gocode v%s  —  %s", Version, ui.Dim.Render(hint))).
				Options(
					huh.NewOption("  New Project      scaffold · push · deploy", "new"),
					huh.NewOption("  Resume           open last project in VS Code", "resume"),
					huh.NewOption("  Push             commit + push to GitHub", "push"),
					huh.NewOption("  Deploy           push live to Netlify", "deploy"),
					huh.NewOption("  All Projects     table view of everything", "list"),
					huh.NewOption("  Status           detailed info on a project", "status"),
					huh.NewOption("  Edit             name · description · visibility", "edit"),
					huh.NewOption("  Delete           remove project everywhere", "delete"),
					huh.NewOption("  Quick Note       append to Obsidian", "note"),
					huh.NewOption("  Settings         config · doctor · backup · stats", "settings"),
				).Value(&action),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	return dispatchMenu(action)
}

func dispatchMenu(action string) error {
	switch action {
	case "new":
		return runNew(newCmd, nil)
	case "resume":
		return runResume(resumeCmd, nil)
	case "push":
		return runPush(pushCmd, nil)
	case "deploy":
		return runDeploy(deployCmd, nil)
	case "list":
		return runList(listCmd, nil)
	case "status":
		return runStatus(statusCmd, nil)
	case "edit":
		return runEdit(editCmd, nil)
	case "delete":
		return runDelete(deleteCmd, nil)
	case "note":
		return runNote(noteCmd, nil)
	case "settings":
		return runSettingsMenu()
	}
	return nil
}

func runSettingsMenu() error {
	var action string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Settings").
				Options(
					huh.NewOption("  Setup Wizard    reconfigure gocode", "config"),
					huh.NewOption("  Doctor          check all dependencies", "doctor"),
					huh.NewOption("  Backup          back up registry + config", "backup"),
					huh.NewOption("  Stats           project statistics", "stats"),
					huh.NewOption("  Update gocode   pull + rebuild", "update"),
				).Value(&action),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	switch action {
	case "config":
		return runConfigWizard()
	case "doctor":
		return runDoctor(doctorCmd, nil)
	case "backup":
		return runBackup(backupCmd, nil)
	case "stats":
		return runStats(statsCmd, nil)
	case "update":
		return runUpdate(updateCmd, nil)
	}
	return nil
}

// ── Shared helpers ────────────────────────────────────────────────────────────

func getRegistry() *registry.Registry {
	return registry.New(config.DataDir() + "/projects.json")
}

func mustLoadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			ui.Error("gocode not configured. Run: gocode config")
		}
		return nil, err
	}
	return cfg, nil
}

func pickActiveProject(prompt string) (string, error) {
	reg := getRegistry()
	active, err := reg.Active()
	if err != nil || len(active) == 0 {
		return "", fmt.Errorf("no active projects — run: gocode new")
	}
	if len(active) == 1 {
		return active[0].Name, nil
	}
	opts := make([]huh.Option[string], len(active))
	for i, p := range active {
		opts[i] = huh.NewOption(p.Name, p.Name)
	}
	var chosen string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title(prompt).Options(opts...).Value(&chosen),
		),
	).WithTheme(huh.ThemeCharm()).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		}
		return "", err
	}
	return chosen, nil
}

func pickAllProject(prompt string) (string, error) {
	reg := getRegistry()
	all, err := reg.All()
	if err != nil || len(all) == 0 {
		return "", fmt.Errorf("no projects found — run: gocode new")
	}
	opts := make([]huh.Option[string], len(all))
	for i, p := range all {
		label := fmt.Sprintf("%-28s  [%s]", p.Name, p.Status)
		opts[i] = huh.NewOption(label, p.Name)
	}
	var chosen string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title(prompt).Options(opts...).Value(&chosen),
		),
	).WithTheme(huh.ThemeCharm()).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		}
		return "", err
	}
	return chosen, nil
}

func openInBrowser(url string) { openURL(url) }
