package cmd

import (
	"os/exec"
	"runtime"
)

// openURL opens a URL in the default browser, cross-platform.
func openURL(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	}
}

// openPath opens a file or directory in the system file manager.
func openPath(path string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", path).Start()
	case "darwin":
		exec.Command("open", path).Start()
	case "windows":
		exec.Command("explorer", path).Start()
	}
}
