package cmd

import (
	"errors"
	"fmt"
	"os"

	ghpkg "github.com/acemaster-gh/gocode/internal/github"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a project from disk, GitHub, and registry",
	RunE:    runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}

	ui.Banner(Version)
	reg := getRegistry()

	name, err := pickAllProject("Delete which project?")
	if err != nil || name == "" {
		return err
	}

	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	// Show what will be affected
	ui.Divider()
	ui.Warn(fmt.Sprintf("Deleting: %s", name))
	ui.Info(fmt.Sprintf("Path:   %s", proj.Path))
	if proj.RepoURL != "local-only" && proj.RepoURL != "" {
		ui.Info(fmt.Sprintf("GitHub: %s", proj.RepoURL))
	}
	ui.Divider()

	var confirmed, deleteGithub bool

	formFields := []huh.Field{
		huh.NewConfirm().
			Title(fmt.Sprintf("Permanently delete '%s'?", name)).
			Value(&confirmed),
	}
	if proj.RepoURL != "local-only" && proj.RepoURL != "" {
		formFields = append(formFields,
			huh.NewConfirm().
				Title("Also delete the GitHub repository?").
				Value(&deleteGithub),
		)
	}

	err = huh.NewForm(
		huh.NewGroup(formFields...),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(ui.Dim.Render("cancelled."))
			return nil
		}
		return err
	}
	if !confirmed {
		fmt.Println(ui.Dim.Render("cancelled."))
		return nil
	}

	// Delete folder
	if _, err := os.Stat(proj.Path); err == nil {
		sp := ui.StartSpinner("Removing folder...")
		rmErr := os.RemoveAll(proj.Path)
		sp.Stop()
		if rmErr != nil {
			ui.Warn(fmt.Sprintf("Folder delete failed: %v", rmErr))
		} else {
			ui.OK("Folder deleted.")
		}
	}

	// Delete GitHub repo
	if deleteGithub && proj.RepoURL != "local-only" && proj.RepoURL != "" {
		sp := ui.StartSpinner("Deleting GitHub repo...")
		ghErr := ghpkg.DeleteRepo(cfg.GHUser, name)
		sp.Stop()
		if ghErr != nil {
			ui.Warn(fmt.Sprintf("GitHub delete failed: %v", ghErr))
		} else {
			ui.OK("GitHub repo deleted.")
		}
	}

	// Remove from registry
	if err := reg.Delete(name); err != nil {
		ui.Warn(fmt.Sprintf("Registry: %v", err))
	} else {
		ui.OK("Removed from registry.")
	}

	// Remove Obsidian note
	vault := obsidian.New(cfg.ObsidianVault)
	vault.DeleteNote(name)

	ui.Done(fmt.Sprintf("Deleted: %s", name))
	return nil
}
