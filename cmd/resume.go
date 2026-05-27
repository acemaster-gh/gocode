package cmd

import (
	"fmt"
	"os/exec"

	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a project in VS Code",
	RunE:  runResume,
}

func runResume(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}

	ui.Banner(Version)
	reg := getRegistry()

	name, err := pickActiveProject("Resume which project?")
	if err != nil || name == "" {
		return err
	}

	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	dirty := gitpkg.UncommittedCount(proj.Path)
	if dirty > 0 {
		ui.Warn(fmt.Sprintf("%d uncommitted change(s) in %s", dirty, name))
	}

	branch := gitpkg.CurrentBranch(proj.Path)
	ui.Step(fmt.Sprintf("Resuming %s  (branch: %s)", name, branch))

	reg.UpdateLastOpened(name)

	if err := exec.Command("code", proj.Path).Start(); err != nil {
		ui.Warn("Could not open VS Code — is it installed?")
	}

	vault := obsidian.New(cfg.ObsidianVault)
	go vault.DailyLog("Resumed: " + name)

	ui.OK(fmt.Sprintf("Opened: %s", proj.Path))

	if len(proj.Todos) > 0 {
		pending := 0
		for _, t := range proj.Todos {
			if !t.Done {
				pending++
			}
		}
		if pending > 0 {
			ui.Info(fmt.Sprintf("%d pending todo(s)  —  run: gocode todo", pending))
		}
	}

	return nil
}
