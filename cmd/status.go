package cmd

import (
	"fmt"
	"os"
	"strings"

	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"info"},
	Short:   "Show detailed status of a project",
	RunE:    runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	name, err := pickAllProject("Status for which project?")
	if err != nil || name == "" {
		return err
	}

	reg := getRegistry()
	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	// ── Status badge ──────────────────────────────────────────────────────────
	statusColor := map[string]lipgloss.Color{
		"active":    "#73c991",
		"completed": "#4ec9b0",
		"archived":  "#808080",
	}
	col := statusColor[proj.Status]
	if col == "" {
		col = "#ffd700"
	}
	badge := lipgloss.NewStyle().
		Foreground(lipgloss.Color(col)).
		Bold(true).
		Render("● " + strings.ToUpper(proj.Status))

	// ── Git info ──────────────────────────────────────────────────────────────
	branch := "—"
	dirty := 0
	if _, err := os.Stat(proj.Path + "/.git"); err == nil {
		branch = gitpkg.CurrentBranch(proj.Path)
		dirty = gitpkg.UncommittedCount(proj.Path)
	}

	dirtyStr := ui.Green.Render("clean")
	if dirty > 0 {
		dirtyStr = ui.Yellow.Render(fmt.Sprintf("%d change(s)", dirty))
	}

	// ── Todos ─────────────────────────────────────────────────────────────────
	totalTodos, doneTodos := 0, 0
	for _, t := range proj.Todos {
		totalTodos++
		if t.Done {
			doneTodos++
		}
	}

	// ── Render panel ─────────────────────────────────────────────────────────
	col2w := 22
	row := func(label, value string) string {
		l := ui.Dim.Render(fmt.Sprintf("  %-*s", col2w, label))
		return l + value
	}

	lines := []string{
		ui.Bold.Render("  " + name) + "  " + badge,
		"",
		row("Description:", proj.Description),
		row("Stack:", proj.Type),
		row("Created:", proj.Created),
		row("Last opened:", proj.LastOpened),
		row("Path:", proj.Path),
	}

	if proj.RepoURL != "" && proj.RepoURL != "local-only" {
		lines = append(lines, row("GitHub:", ui.Cyan.Render(proj.RepoURL)))
	} else {
		lines = append(lines, row("GitHub:", ui.Dim.Render("local only")))
	}

	if proj.NetlifyURL != "" && proj.NetlifyURL != "not-deployed" {
		lines = append(lines, row("Live:", ui.Cyan.Render(proj.NetlifyURL)))
	}

	lines = append(lines,
		"",
		row("Branch:", branch),
		row("Working tree:", dirtyStr),
	)

	if totalTodos > 0 {
		lines = append(lines, row("Todos:",
			fmt.Sprintf("%d/%d complete", doneTodos, totalTodos)))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("75")).
		Padding(0, 1).
		Render(strings.Join(lines, "\n"))

	fmt.Println()
	fmt.Println(box)
	fmt.Println()

	return nil
}
