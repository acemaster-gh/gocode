package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/acemaster-gh/gocode/internal/config"
	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	ghpkg "github.com/acemaster-gh/gocode/internal/github"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/registry"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"rename"},
	Short:   "Edit project name, description, or visibility",
	RunE:    runEdit,
}

func runEdit(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}
	ui.Banner(Version)

	reg := getRegistry()
	name, err := pickAllProject("Edit which project?")
	if err != nil || name == "" {
		return err
	}
	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	var field string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Editing: %s", name)).
				Options(
					huh.NewOption("✎  Name          rename folder + GitHub repo", "name"),
					huh.NewOption("✎  Description   update description everywhere", "description"),
					huh.NewOption("🔒  Visibility    toggle public / private", "visibility"),
					huh.NewOption("📌  Status        change project lifecycle state", "status"),
				).Value(&field),
		),
	).WithTheme(huh.ThemeCharm()).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	vault := obsidian.New(cfg.ObsidianVault)
	switch field {
	case "name":
		return editName(cfg, reg, vault, proj)
	case "description":
		return editDescription(cfg, reg, vault, proj)
	case "visibility":
		return editVisibility(cfg, reg, proj)
	case "status":
		return editStatus(reg, proj)
	}
	return nil
}

func editName(cfg *config.Config, reg *registry.Registry, vault *obsidian.Vault, proj registry.Project) error {
	var newName string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("New project name").
				Placeholder(proj.Name).
				Value(&newName).
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
		),
	).WithTheme(huh.ThemeCharm()).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}
	newName = strings.TrimSpace(newName)
	if newName == proj.Name {
		ui.Info("Name unchanged.")
		return nil
	}

	newPath := filepath.Join(filepath.Dir(proj.Path), newName)
	if err := os.Rename(proj.Path, newPath); err != nil {
		return fmt.Errorf("folder rename: %w", err)
	}
	ui.OK("Folder renamed.")

	newRepo := proj.RepoURL
	if proj.RepoURL != "local-only" && proj.RepoURL != "" {
		sp := ui.StartSpinner("Renaming GitHub repo...")
		url, rErr := ghpkg.RenameRepo(cfg.GHUser, proj.Name, newName)
		sp.Stop()
		if rErr != nil {
			ui.Warn(fmt.Sprintf("GitHub rename failed: %v", rErr))
		} else {
			newRepo = url
			gitpkg.AddRemoteAndPush(newPath, newRepo)
			ui.OK("GitHub repo renamed.")
		}
	}

	vault.RenameNote(proj.Name, newName)
	reg.Rename(proj.Name, newName, newPath)
	if newRepo != proj.RepoURL {
		reg.UpdateField(newName, "repo_url", newRepo)
	}
	ui.Done(fmt.Sprintf("Renamed: %s → %s", proj.Name, newName))
	return nil
}

func editDescription(cfg *config.Config, reg *registry.Registry, vault *obsidian.Vault, proj registry.Project) error {
	var newDesc string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("New description").
				Placeholder(proj.Description).
				Value(&newDesc).
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
	newDesc = strings.TrimSpace(newDesc)

	if proj.RepoURL != "local-only" && proj.RepoURL != "" {
		sp := ui.StartSpinner("Updating GitHub description...")
		ghpkg.SetDescription(cfg.GHUser, proj.Name, newDesc)
		sp.Stop()
	}

	readmePath := filepath.Join(proj.Path, "README.md")
	if data, err := os.ReadFile(readmePath); err == nil {
		updated := strings.Replace(string(data), proj.Description, newDesc, 1)
		os.WriteFile(readmePath, []byte(updated), 0644)
	}

	vault.UpdateDescription(proj.Name, newDesc)
	reg.UpdateField(proj.Name, "description", newDesc)
	ui.Done("Description updated everywhere.")
	return nil
}

func editVisibility(cfg *config.Config, reg *registry.Registry, proj registry.Project) error {
	if proj.RepoURL == "local-only" || proj.RepoURL == "" {
		ui.Error("This project has no GitHub repo.")
		return nil
	}
	var newVis string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Visibility for: %s", proj.Name)).
				Options(
					huh.NewOption("🌐 Public", "public"),
					huh.NewOption("🔒 Private", "private"),
				).Value(&newVis),
		),
	).WithTheme(huh.ThemeCharm()).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}
	sp := ui.StartSpinner(fmt.Sprintf("Setting repo to %s...", newVis))
	ghErr := ghpkg.SetVisibility(cfg.GHUser, proj.Name, newVis)
	sp.Stop()
	if ghErr != nil {
		ui.Error(fmt.Sprintf("Failed: %v", ghErr))
		return nil
	}
	ui.Done(fmt.Sprintf("Visibility set to: %s", newVis))
	return nil
}

func editStatus(reg *registry.Registry, proj registry.Project) error {
	var newStatus string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Status for: %s  (now: %s)", proj.Name, proj.Status)).
				Options(
					huh.NewOption("● Active", "active"),
					huh.NewOption("✓ Completed", "completed"),
					huh.NewOption("○ Archived", "archived"),
				).Value(&newStatus),
		),
	).WithTheme(huh.ThemeCharm()).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}
	reg.SetStatus(proj.Name, newStatus)
	ui.Done(fmt.Sprintf("Status set to: %s", newStatus))
	return nil
}
