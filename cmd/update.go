package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/acemaster-gh/gocode/internal/config"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update gocode to the latest version",
	RunE:  runUpdate,
}

const repoURL = "https://github.com/acemaster-gh/gocode"

func runUpdate(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	fmt.Printf("  %s\n", ui.Bold.Render("Update gocode"))
	ui.Divider()
	ui.Info(fmt.Sprintf("Current version: %s", Version))

	// Find the source directory (go source needs to be present)
	homeDir, _ := os.UserHomeDir()
	srcDir := filepath.Join(homeDir, "gocode-go")

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		ui.Error(fmt.Sprintf("Source not found at: %s", srcDir))
		ui.Info(fmt.Sprintf("Clone the repo first: git clone %s ~/gocode-go", repoURL))
		return nil
	}

	// Git pull
	sp := ui.StartSpinner("Pulling latest changes...")
	pullCmd := exec.Command("git", "pull", "--rebase", "origin", "main")
	pullCmd.Dir = srcDir
	out, pullErr := pullCmd.CombinedOutput()
	sp.Stop()

	if pullErr != nil {
		ui.Error(fmt.Sprintf("git pull failed: %s", string(out)))
		return nil
	}
	ui.OK("Pulled latest changes.")

	// Rebuild
	binaryPath := filepath.Join(config.GocodeDir(), "gocode")

	sp = ui.StartSpinner("Building new binary...")
	buildCmd := exec.Command("go", "build",
		"-ldflags", fmt.Sprintf("-X main.Version=%s", getLatestTag(srcDir)),
		"-o", binaryPath,
		".")
	buildCmd.Dir = srcDir
	buildOut, buildErr := buildCmd.CombinedOutput()
	sp.Stop()

	if buildErr != nil {
		ui.Error(fmt.Sprintf("Build failed: %s", string(buildOut)))
		return nil
	}

	// Make executable
	os.Chmod(binaryPath, 0755)

	ui.OK("Binary rebuilt.")
	ui.Done(fmt.Sprintf("gocode updated → %s", binaryPath))

	// Print PATH hint if needed
	if runtime.GOOS != "windows" {
		ui.Info("If gocode is in your PATH, restart your shell or run: exec $SHELL")
	}

	return nil
}

func getLatestTag(srcDir string) string {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").
		Output()
	if err != nil {
		return "2.0.0"
	}
	v := string(out)
	if len(v) > 0 && v[0] == 'v' {
		v = v[1:]
	}
	return v
}
