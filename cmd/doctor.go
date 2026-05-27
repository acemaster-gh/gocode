package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check all required dependencies",
	RunE:  runDoctor,
}

type depCheck struct {
	name    string
	bin     string
	versionArg string
	required bool
	hint    string
}

var deps = []depCheck{
	{"Git",          "git",     "--version",  true,  "https://git-scm.com/downloads"},
	{"Node.js",      "node",    "--version",  false, "https://nodejs.org"},
	{"npm",          "npm",     "--version",  false, "installed with Node.js"},
	{"VS Code",      "code",    "--version",  false, "https://code.visualstudio.com"},
	{"GitHub CLI",   "gh",      "--version",  false, "https://cli.github.com"},
	{"Netlify CLI",  "netlify", "--version",  false, "npm install -g netlify-cli"},
	{"fzf",          "fzf",     "--version",  false, "brew install fzf  or  apt install fzf"},
}

func runDoctor(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	fmt.Printf("  %s\n\n", ui.Bold.Render("Environment Health Check"))

	colW := 16
	allGood := true

	for _, dep := range deps {
		path, err := exec.LookPath(dep.bin)
		if err != nil {
			mark := ui.Dim.Render("[SKIP]")
			if dep.required {
				mark = ui.Red.Bold(true).Render("[MISS]")
				allGood = false
			}
			label := fmt.Sprintf("  %-*s %s", colW, dep.name, mark)
			fmt.Println(label)
			if dep.required {
				fmt.Printf("    %s → %s\n",
					ui.Dim.Render("install:"),
					ui.Cyan.Render(dep.hint))
			}
			continue
		}

		// Get version
		version := ""
		out, verErr := exec.Command(path, dep.versionArg).CombinedOutput()
		if verErr == nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			if len(lines) > 0 {
				v := lines[0]
				// Extract just the version number
				for _, part := range strings.Fields(v) {
					if strings.HasPrefix(part, "v") || isVersionLike(part) {
						version = part
						break
					}
				}
				if version == "" && len(v) < 30 {
					version = strings.TrimSpace(v)
				}
			}
		}

		mark := ui.Green.Bold(true).Render("[OK]  ")
		label := fmt.Sprintf("  %-*s %s %s",
			colW, dep.name, mark, ui.Dim.Render(version))
		fmt.Println(label)
	}

	fmt.Println()
	ui.Divider()

	// GitHub auth status
	checkGitHubAuth()

	// Netlify auth status
	checkNetlifyAuth()

	fmt.Println()
	if allGood {
		ui.Done("All required dependencies are present.")
	} else {
		ui.Warn("Some required dependencies are missing. See install hints above.")
	}

	return nil
}

func checkGitHubAuth() {
	out, err := exec.Command("gh", "auth", "status").CombinedOutput()
	if err != nil {
		fmt.Printf("  %-16s %s  %s\n",
			"GitHub Auth",
			ui.Yellow.Render("[WARN] "),
			ui.Dim.Render("run: gh auth login"))
		return
	}
	lower := strings.ToLower(string(out))
	if strings.Contains(lower, "logged in") {
		fmt.Printf("  %-16s %s\n",
			"GitHub Auth",
			ui.Green.Bold(true).Render("[OK]   authenticated"))
	} else {
		fmt.Printf("  %-16s %s  %s\n",
			"GitHub Auth",
			ui.Yellow.Render("[WARN] "),
			ui.Dim.Render("run: gh auth login"))
	}
}

func checkNetlifyAuth() {
	out, err := exec.Command("netlify", "status").CombinedOutput()
	if err != nil {
		fmt.Printf("  %-16s %s  %s\n",
			"Netlify Auth",
			ui.Yellow.Render("[WARN] "),
			ui.Dim.Render("run: netlify login"))
		return
	}
	lower := strings.ToLower(string(out))
	if strings.Contains(lower, "email") || strings.Contains(lower, "logged in") {
		fmt.Printf("  %-16s %s\n",
			"Netlify Auth",
			ui.Green.Bold(true).Render("[OK]   authenticated"))
	} else {
		fmt.Printf("  %-16s %s  %s\n",
			"Netlify Auth",
			ui.Yellow.Render("[WARN] "),
			ui.Dim.Render("run: netlify login"))
	}
}

func isVersionLike(s string) bool {
	if len(s) < 3 {
		return false
	}
	dots := strings.Count(s, ".")
	return dots >= 1 && dots <= 3
}
