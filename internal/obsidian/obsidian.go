package obsidian

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Vault struct{ path string }

func New(path string) *Vault      { return &Vault{path: path} }
func (v *Vault) enabled() bool    { return v.path != "" }

func (v *Vault) ProjectNote(name, stack, repoURL, netlifyURL, desc string) error {
	if !v.enabled() { return nil }
	dir := filepath.Join(v.path, "gocode-projects")
	if err := os.MkdirAll(dir, 0755); err != nil { return err }

	repoLine, liveLine := "", ""
	if repoURL != "" && repoURL != "local-only" {
		repoLine = fmt.Sprintf("**Repo:** %s\n", repoURL)
	}
	if netlifyURL != "" && netlifyURL != "not-deployed" {
		liveLine = fmt.Sprintf("**Live:** %s\n", netlifyURL)
	}
	if desc == "" { desc = "_No description._" }

	content := fmt.Sprintf("# %s\n\n**Type:** %s\n**Created:** %s\n**Status:** active\n%s%s**Description:** %s\n\n## Notes\n\n_Add notes here._\n\n## Progress\n\n- [ ] Initial setup complete\n\n---\n_Created by gocode v2.0.0_\n",
		name, stack, time.Now().Format("2006-01-02"),
		repoLine, liveLine, desc)

	return os.WriteFile(filepath.Join(dir, name+".md"), []byte(content), 0644)
}

func (v *Vault) DailyLog(entry string) error {
	if !v.enabled() { return nil }
	today := time.Now().Format("2006-01-02")
	dir := filepath.Join(v.path, "daily")
	if err := os.MkdirAll(dir, 0755); err != nil { return err }
	line := fmt.Sprintf("- %s — gocode: %s\n", time.Now().Format("15:04"), entry)
	f, err := os.OpenFile(filepath.Join(dir, today+".md"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { return err }
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

func (v *Vault) QuickNote(projectName, note string) error {
	if !v.enabled() { return nil }
	notePath := filepath.Join(v.path, "gocode-projects", projectName+".md")
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(notePath), 0755)
		os.WriteFile(notePath, []byte(fmt.Sprintf("# %s\n\n## Notes\n\n", projectName)), 0644)
	}
	f, err := os.OpenFile(notePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil { return err }
	defer f.Close()
	_, err = fmt.Fprintf(f, "\n- **%s** — %s", time.Now().Format("2006-01-02 15:04"), note)
	return err
}

func (v *Vault) MarkComplete(name string) error {
	if !v.enabled() { return nil }
	p := filepath.Join(v.path, "gocode-projects", name+".md")
	data, err := os.ReadFile(p)
	if err != nil { return nil }
	updated := strings.Replace(string(data), "**Status:** active", "**Status:** ✅ completed", 1)
	return os.WriteFile(p, []byte(updated), 0644)
}

func (v *Vault) RenameNote(oldName, newName string) error {
	if !v.enabled() { return nil }
	old := filepath.Join(v.path, "gocode-projects", oldName+".md")
	if _, err := os.Stat(old); os.IsNotExist(err) { return nil }
	return os.Rename(old, filepath.Join(v.path, "gocode-projects", newName+".md"))
}

func (v *Vault) DeleteNote(name string) error {
	if !v.enabled() { return nil }
	p := filepath.Join(v.path, "gocode-projects", name+".md")
	if _, err := os.Stat(p); os.IsNotExist(err) { return nil }
	return os.Remove(p)
}

func (v *Vault) UpdateDescription(name, desc string) error {
	if !v.enabled() { return nil }
	p := filepath.Join(v.path, "gocode-projects", name+".md")
	data, err := os.ReadFile(p)
	if err != nil { return nil }
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "**Description:**") {
			lines[i] = fmt.Sprintf("**Description:** %s", desc)
			break
		}
	}
	return os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0644)
}

func AutoDetectVault() string {
	home, _ := os.UserHomeDir()
	result := ""
	searchVault(home, 0, 4, &result)
	return result
}

func searchVault(dir string, depth, maxDepth int, result *string) {
	if depth > maxDepth || *result != "" { return }
	entries, err := os.ReadDir(dir)
	if err != nil { return }
	for _, e := range entries {
		if e.Name() == ".obsidian" && e.IsDir() { *result = dir; return }
		if e.IsDir() && len(e.Name()) > 0 && e.Name()[0] != '.' {
			searchVault(filepath.Join(dir, e.Name()), depth+1, maxDepth, result)
		}
		if *result != "" { return }
	}
}
