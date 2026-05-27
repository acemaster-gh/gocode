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

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Create, switch, or delete a branch",
	RunE:  runBranch,
}

func runBranch(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	name, err := pickActiveProject("Manage branches for which project?")
	if err != nil || name == "" {
		return err
	}

	reg := getRegistry()
	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	// Show existing branches
	branches, _ := gitpkg.Branches(proj.Path)
	current := gitpkg.CurrentBranch(proj.Path)

	ui.Divider()
	fmt.Printf("  %s  %s\n",
		ui.Cyan.Render("Current branch:"),
		ui.Bold.Render(current))
	if len(branches) > 0 {
		fmt.Printf("  %s  %s\n",
			ui.Dim.Render("All branches:  "),
			ui.Dim.Render(strings.Join(branches, "  |  ")))
	}
	ui.Divider()

	// Pick action
	var action string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Branch operation").
				Options(
					huh.NewOption("⎇  New branch         create and switch", "new"),
					huh.NewOption("↔  Switch branch      checkout existing", "switch"),
					huh.NewOption("✕  Delete branch      remove local branch", "delete"),
				).
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
	case "new":
		return branchNew(proj.Path)
	case "switch":
		return branchSwitch(proj.Path, branches, current)
	case "delete":
		return branchDelete(proj.Path, branches, current)
	}
	return nil
}

func branchNew(path string) error {
	var branchType, branchName string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Branch type").
				Options(
					huh.NewOption("feat/      new feature", "feat"),
					huh.NewOption("fix/       bug fix", "fix"),
					huh.NewOption("refactor/  code restructuring", "refactor"),
					huh.NewOption("docs/      documentation", "docs"),
					huh.NewOption("chore/     maintenance", "chore"),
				).
				Value(&branchType),
			huh.NewInput().
				Title("Branch description (kebab-case)").
				Placeholder("implement-timer-feature").
				Value(&branchName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("branch description cannot be empty")
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

	fullName := branchType + "/" + strings.TrimSpace(branchName)
	if err := gitpkg.CreateBranch(path, fullName); err != nil {
		ui.Error(fmt.Sprintf("Branch creation failed: %v", err))
		return nil
	}
	ui.Done(fmt.Sprintf("Switched to new branch: %s", fullName))
	return nil
}

func branchSwitch(path string, branches []string, current string) error {
	if len(branches) <= 1 {
		ui.Info("Only one branch — create more with: gocode branch → New branch")
		return nil
	}

	opts := []huh.Option[string]{}
	for _, b := range branches {
		label := b
		if b == current {
			label = b + "  (current)"
		}
		opts = append(opts, huh.NewOption(label, b))
	}

	var target string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Switch to which branch?").
				Options(opts...).
				Value(&target),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	if err := gitpkg.SwitchBranch(path, target); err != nil {
		ui.Error(fmt.Sprintf("Switch failed: %v", err))
		return nil
	}
	ui.Done(fmt.Sprintf("Switched to: %s", target))
	return nil
}

func branchDelete(path string, branches []string, current string) error {
	// Only offer non-current branches
	opts := []huh.Option[string]{}
	for _, b := range branches {
		if b == current || b == "main" || b == "master" {
			continue
		}
		opts = append(opts, huh.NewOption(b, b))
	}

	if len(opts) == 0 {
		ui.Info("No branches available to delete.")
		return nil
	}

	var target string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Delete which branch?").
				Options(opts...).
				Value(&target),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	if err := gitpkg.DeleteBranch(path, target); err != nil {
		ui.Error(fmt.Sprintf("Delete failed: %v", err))
		return nil
	}
	ui.Done(fmt.Sprintf("Deleted branch: %s", target))
	return nil
}
