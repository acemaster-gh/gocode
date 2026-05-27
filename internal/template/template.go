package template

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed templates
var embeddedFS embed.FS

// Scaffold copies the template for the given stack into destPath.
func Scaffold(stack, destPath, projectName string) error {
	switch stack {
	case "vanilla":
		return scaffoldEmbedded("templates/vanilla", destPath, projectName)
	case "tailwind":
		return scaffoldEmbedded("templates/tailwind", destPath, projectName)
	case "react":
		return scaffoldReact(destPath, projectName)
	default:
		return fmt.Errorf("unknown stack: %s", stack)
	}
}

func scaffoldEmbedded(fsPath, destPath, projectName string) error {
	sub, err := fs.Sub(embeddedFS, fsPath)
	if err != nil {
		return fmt.Errorf("template sub: %w", err)
	}
	return fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dest := filepath.Join(destPath, path)
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		data, err := fs.ReadFile(sub, path)
		if err != nil {
			return err
		}
		content := strings.ReplaceAll(string(data), "PROJECT_NAME", projectName)
		return os.WriteFile(dest, []byte(content), 0644)
	})
}

func scaffoldReact(destPath, projectName string) error {
	tmp, err := os.MkdirTemp("", "gocode-react-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	c := exec.Command("npm", "create", "vite@latest", projectName,
		"--", "--template", "react", "--yes")
	c.Dir = tmp
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("npm create vite: %w", err)
	}

	src := filepath.Join(tmp, projectName)
	if err := copyDir(src, destPath); err != nil {
		return err
	}

	install := exec.Command("npm", "install", "--silent")
	install.Dir = destPath
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	return install.Run()
}

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, d.Type().Perm())
	})
}

// GenerateREADME writes a project README.md.
func GenerateREADME(name, stack, path, repoURL, netlifyURL, desc string) error {
	today := time.Now().Format("2006-01-02")
	descText := desc
	if descText == "" {
		descText = "_Add a short description of what this project does._"
	}

	var liveMD, repoMD string
	if netlifyURL != "" && netlifyURL != "not-deployed" {
		liveMD = fmt.Sprintf("## Live Demo\n\n**%s**\n", netlifyURL)
	} else {
		liveMD = "## Live Demo\n\n_Not deployed yet._\n"
	}
	if repoURL != "" && repoURL != "local-only" {
		repoMD = fmt.Sprintf("## Repo\n\n**%s**\n", repoURL)
	}

	content := fmt.Sprintf("# %s\n\n**Started:** %s  |  **Stack:** %s\n\n---\n\n## Description\n\n%s\n\n%s\n%s\n## Built With\n\n%s\n\n## Preview\n\n```\n%s\n```\n\n## Progress\n\n| Date | Update |\n|------|--------|\n| %s | Project created |\n\n---\n\n_Built with [gocode](https://github.com/acemaster-gh/gocode) v2.0.0 ⚡_\n",
		name, today, stack,
		descText, liveMD, repoMD,
		stackBullets(stack), asciiPreview(stack), today)

	return os.WriteFile(filepath.Join(path, "README.md"), []byte(content), 0644)
}

// UpdateREADMEOnComplete appends a completion row.
func UpdateREADMEOnComplete(projectPath string) {
	p := filepath.Join(projectPath, "README.md")
	data, err := os.ReadFile(p)
	if err != nil {
		return
	}
	today := time.Now().Format("2006-01-02")
	row := fmt.Sprintf("| %s | ✅ Project completed |", today)
	updated := strings.Replace(string(data), "| Project created |",
		"| Project created |\n"+row, 1)
	os.WriteFile(p, []byte(updated), 0644)
}

func stackBullets(stack string) string {
	switch stack {
	case "tailwind":
		return "- HTML\n- Tailwind CSS v4\n- JavaScript"
	case "react":
		return "- React\n- Vite\n- JavaScript"
	default:
		return "- HTML\n- CSS (CSS variables)\n- JavaScript (Vanilla)"
	}
}

func asciiPreview(stack string) string {
	if stack == "react" {
		return "App\n ├── Header\n ├── Main\n └── Footer"
	}
	return "┌────────────────────────────────┐\n│  Header / Nav                  │\n│  ┌────────┐  ┌────────┐        │\n│  │ Button │  │ Button │        │\n│  └────────┘  └────────┘        │\n│  Main content area             │\n└────────────────────────────────┘"
}
