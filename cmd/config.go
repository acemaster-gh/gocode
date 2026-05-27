package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/acemaster-gh/gocode/internal/config"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"setup", "configure"},
	Short:   "Run the setup wizard",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigWizard()
	},
}

func runConfigWizard() error {
	ui.Banner(Version)

	existing, _ := config.Load()
	defaults := &config.Config{
		ProjectsDir:        defaultProjectsDir(),
		SleepBeforeBrowser: 5,
		BrowserTabs: []string{
			"https://claude.ai",
			"https://chatgpt.com",
			"https://developer.mozilla.org",
			"http://localhost:5500",
		},
	}
	if existing != nil {
		defaults = existing
	}
	if defaults.ObsidianVault == "" {
		defaults.ObsidianVault = obsidian.AutoDetectVault()
	}

	fmt.Printf("  %s\n\n", ui.Dim.Render("Configure gocode — saved to ~/.gocode/config.toml"))

	ghUser      := defaults.GHUser
	projectsDir := defaults.ProjectsDir
	vaultPath   := defaults.ObsidianVault
	tabs        := strings.Join(defaults.BrowserTabs, "\n")
	sleepStr    := fmt.Sprintf("%d", defaults.SleepBeforeBrowser)

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("GitHub username").
				Placeholder("acemaster-gh").
				Value(&ghUser).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("GitHub username is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Projects directory  (will be created if missing)").
				Placeholder(defaults.ProjectsDir).
				Value(&projectsDir),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Obsidian vault path  (leave empty to skip)").
				Placeholder("/home/you/Desktop/MyVault").
				Value(&vaultPath),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Browser tabs on project open  (one URL per line)").
				Placeholder("https://claude.ai").
				Value(&tabs),
			huh.NewInput().
				Title("Seconds to wait before opening browser").
				Placeholder("5").
				Value(&sleepStr),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(ui.Dim.Render("Setup cancelled."))
			return nil
		}
		return err
	}

	sleep := 5
	fmt.Sscanf(strings.TrimSpace(sleepStr), "%d", &sleep)

	tabList := []string{}
	for _, line := range strings.Split(tabs, "\n") {
		if t := strings.TrimSpace(line); t != "" && strings.HasPrefix(t, "http") {
			tabList = append(tabList, t)
		}
	}
	if len(tabList) == 0 {
		tabList = defaults.BrowserTabs
	}

	cfg := &config.Config{
		GHUser:             strings.TrimSpace(ghUser),
		ProjectsDir:        expandHome(strings.TrimSpace(projectsDir)),
		ObsidianVault:      expandHome(strings.TrimSpace(vaultPath)),
		SleepBeforeBrowser: sleep,
		BrowserTabs:        tabList,
	}

	os.MkdirAll(cfg.ProjectsDir, 0755)
	os.MkdirAll(config.DataDir(), 0755)
	os.MkdirAll(config.PluginsDir(), 0755)

	sp := ui.StartSpinner("Saving config...")
	saveErr := config.Save(cfg)
	sp.Stop()
	if saveErr != nil {
		return fmt.Errorf("saving config: %w", saveErr)
	}

	ui.Divider()
	ui.OK(fmt.Sprintf("GitHub user:   %s", cfg.GHUser))
	ui.OK(fmt.Sprintf("Projects dir:  %s", cfg.ProjectsDir))
	if cfg.ObsidianVault != "" {
		ui.OK(fmt.Sprintf("Obsidian:      %s", cfg.ObsidianVault))
	}
	ui.OK(fmt.Sprintf("Config saved:  %s", config.ConfigPath()))
	ui.Divider()
	ui.Done("gocode is configured.  Run: gocode new")
	return nil
}

func defaultProjectsDir() string {
	home, _ := os.UserHomeDir()
	for _, c := range []string{
		home + "/Desktop/js-projects",
		home + "/Documents/projects",
		home + "/projects",
	} {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return home + "/projects"
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return path
}
