package cmd

import (
	"fmt"
	"os"
	"os/exec"

	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show recent git commits",
	RunE:  runLog,
}

func runLog(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	name, err := pickActiveProject("Show log for which project?")
	if err != nil || name == "" {
		return err
	}

	reg := getRegistry()
	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	fmt.Println()
	fmt.Printf("  %s\n", ui.Bold.Render(fmt.Sprintf("Last 15 commits — %s", name)))
	ui.Divider()

	out, err := gitpkg.Log(proj.Path, 15)
	if err != nil {
		ui.Error("Could not read git log.")
		return nil
	}
	if out == "" {
		ui.Info("No commits yet.")
		return nil
	}

	fmt.Println(out)
	ui.Divider()

	// Offer full log in pager if there are more commits
	allOut, _ := gitpkg.Log(proj.Path, 9999)
	lineCount := countLines(allOut)
	if lineCount > 15 {
		fmt.Printf("\n  %s\n",
			ui.Dim.Render(fmt.Sprintf("(%d total commits — press q to exit pager)", lineCount)))
		if ui.Confirm("Open full log in pager?") {
			c := exec.Command("git", "log", "--oneline", "--color=always")
			c.Dir = proj.Path
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			c.Run()
		}
	}

	return nil
}

func countLines(s string) int {
	n := 0
	for _, c := range s {
		if c == '\n' {
			n++
		}
	}
	return n + 1
}
