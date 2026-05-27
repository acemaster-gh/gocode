package netlify

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func available() bool {
	_, err := exec.LookPath("netlify")
	return err == nil
}

// AuthStatus checks netlify authentication.
func AuthStatus() bool {
	if !available() {
		return false
	}
	out, err := exec.Command("netlify", "status").CombinedOutput()
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(out))
	return !strings.Contains(lower, "not logged") &&
		!strings.Contains(lower, "not authenticated") &&
		!strings.Contains(lower, "please log in") &&
		!strings.Contains(lower, "login required")
}

type siteJSON struct {
	ID     string `json:"id"`
	SSLURL string `json:"ssl_url"`
	URL    string `json:"url"`
}

// CreateSite creates a Netlify site.
// Returns (siteID, siteURL, error). Falls back to auto-named site if name taken.
func CreateSite(name string) (string, string, error) {
	if !available() {
		return "", "", fmt.Errorf("netlify CLI not found")
	}

	// Try named site first
	out, err := exec.Command("netlify", "sites:create",
		"--name", name, "--json").Output()
	if err != nil {
		// Name likely taken — try auto
		out, err = exec.Command("netlify", "sites:create", "--json").Output()
		if err != nil {
			return "", "", fmt.Errorf("netlify sites:create failed")
		}
	}

	var site siteJSON
	if err := json.Unmarshal(out, &site); err != nil || site.ID == "" {
		return "", "", fmt.Errorf("could not parse netlify site response")
	}

	url := site.SSLURL
	if url == "" {
		url = site.URL
	}
	return site.ID, url, nil
}

// Deploy runs netlify deploy --prod to the given site.
// Returns the live URL.
func Deploy(dir, siteID string) (string, error) {
	cmd := exec.Command("netlify", "deploy", "--prod", "--dir", ".", "--site", siteID)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("netlify deploy: %s", strings.TrimSpace(string(out)))
	}

	// Extract live URL from output
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, ".netlify.app") {
			// Find the URL
			for _, word := range strings.Fields(line) {
				if strings.HasPrefix(word, "https://") && strings.Contains(word, ".netlify.app") {
					return strings.TrimRight(word, ","), nil
				}
			}
		}
	}
	return "", nil
}

// LinkSite links a local project directory to a Netlify site.
func LinkSite(dir, siteID string) error {
	cmd := exec.Command("netlify", "link", "--id", siteID)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Redeploy runs netlify deploy --prod in the given dir (uses existing link).
func Redeploy(dir string) (string, error) {
	cmd := exec.Command("netlify", "deploy", "--prod", "--dir", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("redeploy: %s", strings.TrimSpace(string(out)))
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, ".netlify.app") {
			for _, word := range strings.Fields(line) {
				if strings.HasPrefix(word, "https://") && strings.Contains(word, ".netlify.app") {
					return strings.TrimRight(word, ","), nil
				}
			}
		}
	}
	return "", nil
}

// Full is the high-level create+deploy+link workflow used in --new.
// Returns live URL or "not-deployed".
func Full(projectName, dir, repoURL string) string {
	if !available() {
		return "not-deployed"
	}
	if !AuthStatus() {
		return "not-deployed"
	}

	siteID, siteURL, err := CreateSite(projectName)
	if err != nil || siteID == "" {
		return "not-deployed"
	}
	_ = siteURL

	liveURL, err := Deploy(dir, siteID)
	if err != nil || liveURL == "" {
		return "not-deployed"
	}

	// Link site if we have a git repo
	if repoURL != "local-only" && repoURL != "" {
		_ = LinkSite(dir, siteID) // best-effort
	}

	return liveURL
}
