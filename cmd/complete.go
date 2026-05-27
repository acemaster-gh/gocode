package cmd

import (
	"errors"
	"fmt"

	"github.com/acemaster-gh/gocode/internal/config"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/plugin"
	tmpl "github.com/acemaster-gh/gocode/internal/template"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a project as complete",
	RunE:  runComplete,
}

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a project",
	RunE:  runArchive,
}

func runComplete(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}

	ui.Banner(Version)
	reg := getRegistry()

	name, err := pickActiveProject("Mark which project as complete?")
	if err != nil || name == "" {
		return err
	}

	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	var confirm bool
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Mark '%s' as complete?", name)).
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil || !confirm {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	// Update registry
	reg.SetStatus(name, "completed")

	// Update README
	tmpl.UpdateREADMEOnComplete(proj.Path)
	ui.OK("README progress table updated.")

	// Update Obsidian
	vault := obsidian.New(cfg.ObsidianVault)
	vault.MarkComplete(name)
	vault.DailyLog("Completed: " + name)

	// Plugin hook
	go plugin.RunHooks(config.PluginsDir(), plugin.HookPostComplete, plugin.Payload{
		Hook:    plugin.HookPostComplete,
		Project: name,
	})

	ui.Done(fmt.Sprintf("Project completed: %s 🎉", name))
	return nil
}

func runArchive(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)
	reg := getRegistry()

	name, err := pickAllProject("Archive which project?")
	if err != nil || name == "" {
		return err
	}

	var confirm bool
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Archive '%s'?", name)).
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil || !confirm {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	reg.SetStatus(name, "archived")
	ui.Done(fmt.Sprintf("Archived: %s", name))
	return nil
}
