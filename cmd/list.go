package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all projects",
	RunE:    runList,
}

func runList(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	reg := getRegistry()
	all, err := reg.All()
	if err != nil {
		return err
	}

	if len(all) == 0 {
		ui.Info("No projects yet. Run: gocode new")
		return nil
	}

	// ── Header ────────────────────────────────────────────────────────────────
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("75"))
	header := fmt.Sprintf("  %s  %-26s %-12s %-14s %s",
		headerStyle.Render("  "),
		headerStyle.Render("NAME"),
		headerStyle.Render("STACK"),
		headerStyle.Render("CREATED"),
		headerStyle.Render("GITHUB / LIVE"),
	)
	fmt.Println(header)
	ui.Divider()

	// ── Rows ──────────────────────────────────────────────────────────────────
	for _, p := range all {
		icon, col := statusIndicator(p.Status)

		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(col))
		indicator := nameStyle.Render(icon)

		// truncate long names
		name := p.Name
		if len(name) > 24 {
			name = name[:21] + "..."
		}

		// Build link summary
		links := []string{}
		if p.RepoURL != "" && p.RepoURL != "local-only" {
			links = append(links, ui.Cyan.Render("⎇ github"))
		}
		if p.NetlifyURL != "" && p.NetlifyURL != "not-deployed" {
			links = append(links, ui.Green.Render("🌐 live"))
		}
		linkStr := strings.Join(links, "  ")
		if linkStr == "" {
			linkStr = ui.Dim.Render("local")
		}

		fmt.Fprintf(os.Stderr, "  %s  %-26s %-12s %-14s %s\n",
			indicator,
			nameStyle.Render(name),
			ui.Dim.Render(p.Type),
			ui.Dim.Render(p.Created),
			linkStr,
		)
	}

	ui.Divider()

	// ── Footer counts ──────────────────────────────────────────────────────────
	active, completed, archived := reg.CountByStatus()
	fmt.Fprintf(os.Stderr, "  %s  %s  %s\n\n",
		ui.Green.Render(fmt.Sprintf("● %d active", active)),
		ui.Cyan.Render(fmt.Sprintf("✓ %d completed", completed)),
		ui.Dim.Render(fmt.Sprintf("○ %d archived", archived)),
	)

	return nil
}

func statusIndicator(status string) (icon string, color string) {
	switch status {
	case "active":
		return "●", "#73c991"
	case "completed":
		return "✓", "#4ec9b0"
	case "archived":
		return "○", "#808080"
	default:
		return "?", "#ffd700"
	}
}
