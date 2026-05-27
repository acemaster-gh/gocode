package cmd

import (
	"errors"
	"fmt"
	"strings"

	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/plugin"
	"github.com/acemaster-gh/gocode/internal/config"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Commit and push to GitHub",
	RunE:  runPush,
}

func runPush(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}

	ui.Banner(Version)

	reg := getRegistry()

	// Pick project
	name, err := pickActiveProject("Push which project?")
	if err != nil || name == "" {
		return err
	}

	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	// Check for changes
	if gitpkg.IsClean(proj.Path) {
		ui.Warn("Nothing to commit — working tree is clean.")
		return nil
	}

	// Show what changed
	changes := gitpkg.Status(proj.Path)
	if changes != "" {
		ui.Divider()
		fmt.Println(ui.Cyan.Render("  Changed files:"))
		for _, line := range strings.Split(changes, "\n") {
			fmt.Println(ui.Dim.Render("  " + line))
		}
		ui.Divider()
	}

	// Commit form
	var commitType, commitMsg string

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Commit type").
				Options(
					huh.NewOption("feat      ✨  new feature or enhancement", "feat"),
					huh.NewOption("fix       🐛  bug fix", "fix"),
					huh.NewOption("docs      📝  documentation changes", "docs"),
					huh.NewOption("style     💄  formatting, no logic change", "style"),
					huh.NewOption("refactor  ♻   code restructuring", "refactor"),
					huh.NewOption("perf      ⚡  performance improvement", "perf"),
					huh.NewOption("test      🧪  adding or updating tests", "test"),
					huh.NewOption("chore     🔧  build, config, tooling", "chore"),
				).
				Value(&commitType),
			huh.NewInput().
				Title("Commit message").
				Placeholder("what did you change?").
				Value(&commitMsg).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("message cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(ui.Dim.Render("cancelled."))
			return nil
		}
		return err
	}

	fullMsg := fmt.Sprintf("%s: %s", commitType, strings.TrimSpace(commitMsg))

	sp := ui.StartSpinner("Pushing to GitHub...")
	pushErr := gitpkg.CommitAndPush(proj.Path, fullMsg)
	sp.Stop()

	if pushErr != nil {
		ui.Error(fmt.Sprintf("Push failed: %v", pushErr))
		ui.Info("Try: cd " + proj.Path + " && git push --set-upstream origin main")
		return nil
	}

	reg.UpdateLastOpened(name)

	// Obsidian + plugin hooks in background
	vault := obsidian.New(cfg.ObsidianVault)
	go vault.DailyLog(fmt.Sprintf("Pushed [%s]: %s", name, fullMsg))
	go plugin.RunHooks(config.PluginsDir(), plugin.HookPostPush, plugin.Payload{
		Hook:    plugin.HookPostPush,
		Project: name,
		Data:    map[string]string{"message": fullMsg},
	})

	ui.OK(fmt.Sprintf("Pushed: %s", ui.Bold.Render(fullMsg)))
	return nil
}
