package cmd

import (
	"errors"
	"fmt"
	"strings"

	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var stashCmd = &cobra.Command{
	Use:   "stash",
	Short: "Save, apply, or drop git stashes",
	RunE:  runStash,
}

func runStash(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	name, err := pickActiveProject("Manage stashes for which project?")
	if err != nil || name == "" {
		return err
	}

	reg := getRegistry()
	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	stashes, _ := gitpkg.StashList(proj.Path)

	// Show stash list
	if len(stashes) > 0 {
		ui.Divider()
		fmt.Println(ui.Cyan.Render("  Saved stashes:"))
		for i, s := range stashes {
			fmt.Printf("  %s  %s\n",
				ui.Yellow.Render(fmt.Sprintf("[%d]", i)),
				ui.Dim.Render(s))
		}
		ui.Divider()
	}

	var action string
	opts := []huh.Option[string]{
		huh.NewOption("💾  Save stash       push current changes to stack", "save"),
	}
	if len(stashes) > 0 {
		opts = append(opts,
			huh.NewOption("📤  Apply stash      restore saved changes", "apply"),
			huh.NewOption("✕   Drop stash       remove without applying", "drop"),
		)
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Stash operation").
				Options(opts...).
				Value(&action),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	switch action {
	case "save":
		return stashSave(proj.Path)
	case "apply":
		return stashApply(proj.Path, stashes)
	case "drop":
		return stashDrop(proj.Path, stashes)
	}
	return nil
}

func stashSave(path string) error {
	var msg string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Stash description").
				Placeholder("WIP: implementing feature X").
				Value(&msg).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("description cannot be empty")
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

	if err := gitpkg.StashPush(path, strings.TrimSpace(msg)); err != nil {
		ui.Error(fmt.Sprintf("Stash failed: %v", err))
		return nil
	}
	ui.Done("Changes stashed.")
	return nil
}

func stashApply(path string, stashes []string) error {
	opts := make([]huh.Option[int], len(stashes))
	for i, s := range stashes {
		opts[i] = huh.NewOption(s, i)
	}

	var idx int
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Apply which stash?").
				Options(opts...).
				Value(&idx),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	if err := gitpkg.StashApply(path, idx); err != nil {
		ui.Error(fmt.Sprintf("Apply failed: %v", err))
		return nil
	}
	ui.Done(fmt.Sprintf("Applied stash@{%d}", idx))
	return nil
}

func stashDrop(path string, stashes []string) error {
	opts := make([]huh.Option[int], len(stashes))
	for i, s := range stashes {
		opts[i] = huh.NewOption(s, i)
	}

	var idx int
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Drop which stash?").
				Options(opts...).
				Value(&idx),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	if err := gitpkg.StashDrop(path, idx); err != nil {
		ui.Error(fmt.Sprintf("Drop failed: %v", err))
		return nil
	}
	ui.Done(fmt.Sprintf("Dropped stash@{%d}", idx))
	return nil
}
