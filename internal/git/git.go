package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// run executes a git command in dir and returns combined output.
func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// runSilent runs a git command, inheriting stdio (for interactive use).
func runSilent(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InitAndCommit runs git init + first commit in dir.
func InitAndCommit(dir string) error {
	if _, err := run(dir, "init", "-q"); err != nil {
		return fmt.Errorf("git init: %w", err)
	}
	_, _ = run(dir, "checkout", "-b", "main")
	_, _ = run(dir, "branch", "-M", "main")
	if _, err := run(dir, "add", "."); err != nil {
		return err
	}
	_, err := run(dir, "commit", "-m", "feat: initial commit — created with gocode v2.0.0", "-q")
	return err
}

// AddRemoteAndPush adds the remote and pushes main.
func AddRemoteAndPush(dir, url string) error {
	run(dir, "remote", "add", "origin", url+".git")
	run(dir, "remote", "set-url", "origin", url+".git")
	run(dir, "branch", "-M", "main")
	_, err := run(dir, "push", "-u", "origin", "main", "-q")
	return err
}

// CommitAndPush stages all, commits with msg, pushes.
func CommitAndPush(dir, msg string) error {
	if _, err := run(dir, "add", "."); err != nil {
		return err
	}
	if _, err := run(dir, "commit", "-m", msg, "-q"); err != nil {
		return err
	}
	_, err := run(dir, "push", "-q")
	return err
}

// PushFinal pushes README update after create (best-effort).
func PushFinal(dir, repoURL string) {
	if repoURL == "local-only" || repoURL == "" {
		return
	}
	run(dir, "add", ".")
	run(dir, "commit", "-m", "docs: update README", "-q")
	run(dir, "push", "-q")
}

// IsClean returns true if there are no uncommitted changes.
func IsClean(dir string) bool {
	out, err := run(dir, "status", "--short")
	return err == nil && strings.TrimSpace(out) == ""
}

// UncommittedCount returns the number of changed files.
func UncommittedCount(dir string) int {
	out, err := run(dir, "status", "--short")
	if err != nil || strings.TrimSpace(out) == "" {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(out), "\n"))
}

// CurrentBranch returns the current branch name.
func CurrentBranch(dir string) string {
	out, err := run(dir, "branch", "--show-current")
	if err != nil {
		return "?"
	}
	return out
}

// Log returns the last n commit lines.
func Log(dir string, n int) (string, error) {
	out, err := run(dir, "log", "--oneline", fmt.Sprintf("-%d", n), "--color=always")
	return out, err
}

// Diff returns uncommitted changes.
func Diff(dir string) (string, error) {
	out, err := run(dir, "diff", "--color=always")
	return out, err
}

// Branches returns a list of local branches.
func Branches(dir string) ([]string, error) {
	out, err := run(dir, "branch", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	var branches []string
	for _, b := range strings.Split(out, "\n") {
		b = strings.TrimSpace(b)
		if b != "" {
			branches = append(branches, b)
		}
	}
	return branches, nil
}

// CreateBranch creates and switches to a new branch.
func CreateBranch(dir, name string) error {
	_, err := run(dir, "checkout", "-b", name)
	return err
}

// SwitchBranch switches to an existing branch.
func SwitchBranch(dir, name string) error {
	_, err := run(dir, "checkout", name)
	return err
}

// DeleteBranch deletes a local branch.
func DeleteBranch(dir, name string) error {
	_, err := run(dir, "branch", "-d", name)
	return err
}

// StashList returns stash entries.
func StashList(dir string) ([]string, error) {
	out, err := run(dir, "stash", "list")
	if err != nil || out == "" {
		return nil, err
	}
	return strings.Split(out, "\n"), nil
}

// StashPush saves current changes to stash with a message.
func StashPush(dir, msg string) error {
	_, err := run(dir, "stash", "push", "-m", msg)
	return err
}

// StashApply applies stash at index.
func StashApply(dir string, idx int) error {
	_, err := run(dir, "stash", "apply", fmt.Sprintf("stash@{%d}", idx))
	return err
}

// StashDrop drops stash at index.
func StashDrop(dir string, idx int) error {
	_, err := run(dir, "stash", "drop", fmt.Sprintf("stash@{%d}", idx))
	return err
}

// Status returns a short status string.
func Status(dir string) string {
	out, err := run(dir, "status", "--short")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}
