package cmd

import (
	"fmt"

	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show uncommitted changes",
	RunE:  runDiff,
}

func runDiff(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	name, err := pickActiveProject("Show diff for which project?")
	if err != nil || name == "" {
		return err
	}

	reg := getRegistry()
	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	fmt.Println()
	fmt.Printf("  %s\n", ui.Bold.Render(fmt.Sprintf("Diff — %s", name)))
	ui.Divider()

	out, err := gitpkg.Diff(proj.Path)
	if err != nil {
		ui.Error("Could not read diff.")
		return nil
	}
	if out == "" {
		ui.Green.Render("  Working tree clean — nothing to show.")
		ui.Info("Working tree clean — nothing to show.")
		return nil
	}

	fmt.Println(out)
	ui.Divider()
	return nil
}
