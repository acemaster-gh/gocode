package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Add a quick note to Obsidian",
	RunE:  runNote,
}

func runNote(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}

	if cfg.ObsidianVault == "" {
		ui.Error("No Obsidian vault configured. Run: gocode config")
		return nil
	}

	ui.Banner(Version)

	// Pick project or daily
	var target string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Add note where?").
				Options(
					huh.NewOption("📅  Daily note    append to today's journal", "daily"),
					huh.NewOption("📁  Project note  add to a specific project", "project"),
				).
				Value(&target),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	var noteText string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Your note").
				Placeholder("quick thought or task...").
				Value(&noteText).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("note cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	noteText = strings.TrimSpace(noteText)
	vault := obsidian.New(cfg.ObsidianVault)

	switch target {
	case "daily":
		if err := vault.DailyLog(noteText); err != nil {
			ui.Error(fmt.Sprintf("Could not write daily note: %v", err))
			return nil
		}
		ui.Done("Added to today's daily note.")

	case "project":
		name, err := pickAllProject("Add note to which project?")
		if err != nil || name == "" {
			return err
		}
		if err := vault.QuickNote(name, noteText); err != nil {
			ui.Error(fmt.Sprintf("Could not write project note: %v", err))
			return nil
		}
		ui.Done(fmt.Sprintf("Note added to: %s", name))
	}

	return nil
}
