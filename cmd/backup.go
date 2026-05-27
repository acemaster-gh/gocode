package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/acemaster-gh/gocode/internal/config"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Back up registry and config",
	RunE:  runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	backupDir := filepath.Join(config.GocodeDir(), "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("creating backup dir: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	backed := 0
	total := 0

	files := []struct {
		src  string
		name string
	}{
		{filepath.Join(config.DataDir(), "projects.json"), "projects"},
		{config.ConfigPath(), "config"},
	}

	fmt.Println()
	fmt.Printf("  %s\n", ui.Bold.Render("Backup"))
	ui.Divider()

	for _, f := range files {
		total++
		dest := filepath.Join(backupDir, fmt.Sprintf("%s_%s.bak", f.name, timestamp))

		if _, err := os.Stat(f.src); os.IsNotExist(err) {
			fmt.Printf("  %-16s %s\n",
				f.name, ui.Dim.Render("not found, skipping"))
			continue
		}

		sp := ui.StartSpinner(fmt.Sprintf("Backing up %s...", f.name))
		err := copyFile(f.src, dest)
		sp.Stop()

		if err != nil {
			fmt.Printf("  %-16s %s — %v\n",
				f.name, ui.Red.Render("[FAIL]"), err)
		} else {
			backed++
			size := fileSize(dest)
			fmt.Printf("  %-16s %s  %s\n",
				f.name,
				ui.Green.Render("[OK]"),
				ui.Dim.Render(fmt.Sprintf("%s  (%s)", filepath.Base(dest), size)))
		}
	}

	ui.Divider()
	fmt.Printf("  %s\n\n", ui.Dim.Render(fmt.Sprintf("Backed up %d/%d files to: %s", backed, total, backupDir)))

	if backed == total {
		ui.Done("Backup complete.")
	} else {
		ui.Warn("Some files could not be backed up.")
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func fileSize(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "?"
	}
	b := info.Size()
	switch {
	case b < 1024:
		return fmt.Sprintf("%d B", b)
	case b < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	default:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	}
}
