package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func available() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// AuthStatus returns true if gh is authenticated.
func AuthStatus() bool {
	if !available() {
		return false
	}
	cmd := exec.Command("gh", "auth", "status")
	return cmd.Run() == nil
}

// CreateRepo creates a GitHub repo with the given visibility.
// Returns the repo URL or "local-only" on failure.
func CreateRepo(name, visibility string) (string, error) {
	if !available() {
		return "local-only", fmt.Errorf("gh not found")
	}
	if !AuthStatus() {
		return "local-only", fmt.Errorf("not authenticated — run: gh auth login")
	}

	vis := "--" + visibility
	cmd := exec.Command("gh", "repo", "create", name, vis)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "local-only", fmt.Errorf("gh repo create: %s", strings.TrimSpace(string(out)))
	}

	// gh outputs the URL on success
	url := strings.TrimSpace(string(out))
	if !strings.HasPrefix(url, "https://") {
		// Try to parse from output
		for _, line := range strings.Split(url, "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "https://github.com") {
				url = strings.TrimSpace(line)
				break
			}
		}
	}
	return url, nil
}

// RenameRepo renames a GitHub repo. Returns new URL.
func RenameRepo(user, oldName, newName string) (string, error) {
	cmd := exec.Command("gh", "repo", "rename", newName,
		"-R", user+"/"+oldName, "--yes")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("rename failed: %s", strings.TrimSpace(string(out)))
	}
	return "https://github.com/" + user + "/" + newName, nil
}

// SetVisibility sets a repo's visibility.
func SetVisibility(user, name, visibility string) error {
	cmd := exec.Command("gh", "repo", "edit",
		user+"/"+name, "--visibility", visibility)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("set visibility: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// SetDescription updates the repo description.
func SetDescription(user, name, desc string) error {
	cmd := exec.Command("gh", "repo", "edit",
		user+"/"+name, "--description", desc)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("set description: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// DeleteRepo deletes a repository.
func DeleteRepo(user, name string) error {
	cmd := exec.Command("gh", "repo", "delete", user+"/"+name, "--yes")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("delete repo: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// SearchSimilar searches GitHub for similar repos by topic.
type SearchResult struct {
	FullName        string `json:"fullName"`
	StargazersCount int    `json:"stargazersCount"`
	URL             string `json:"url"`
}

func SearchSimilar(projectName string) ([]SearchResult, error) {
	if !available() {
		return nil, nil
	}

	// Clean up project name for search
	topic := strings.NewReplacer("-", " ", "_", " ").Replace(projectName)
	for _, stop := range []string{"app", "project", "tool", "my", "the", "a"} {
		topic = strings.ReplaceAll(topic, " "+stop+" ", " ")
		topic = strings.TrimPrefix(topic, stop+" ")
		topic = strings.TrimSuffix(topic, " "+stop)
	}
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil, nil
	}

	cmd := exec.Command("gh", "search", "repos", topic,
		"--language", "javascript",
		"--sort", "stars",
		"--limit", "5",
		"--json", "fullName,stargazersCount,url")
	cmd.Args = append(cmd.Args) // no timeout here; caller should use goroutine

	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	var results []SearchResult
	if err := json.Unmarshal(out, &results); err != nil {
		return nil, nil
	}
	return results, nil
}
