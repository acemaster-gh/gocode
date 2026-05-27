package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/acemaster-gh/gocode/internal/registry"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show project statistics",
	RunE:  runStats,
}

func runStats(cmd *cobra.Command, args []string) error {
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

	active, completed, archived := reg.CountByStatus()
	total := len(all)

	stacks := map[string]int{}
	for _, p := range all {
		stacks[p.Type]++
	}

	withGitHub, withNetlify := 0, 0
	for _, p := range all {
		if p.RepoURL != "local-only" && p.RepoURL != "" {
			withGitHub++
		}
		if p.NetlifyURL != "not-deployed" && p.NetlifyURL != "" {
			withNetlify++
		}
	}

	// ── Render ────────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Printf("  %s\n\n", ui.Bold.Render("Project Statistics"))

	fmt.Printf("  %s\n", ui.Bold.Render("Totals"))
	ui.Divider()
	statRow("Total projects", fmt.Sprintf("%d", total))
	statRow("Active",         fmt.Sprintf("%d", active))
	statRow("Completed",      fmt.Sprintf("%d", completed))
	statRow("Archived",       fmt.Sprintf("%d", archived))
	statRow("On GitHub",      fmt.Sprintf("%d / %d", withGitHub, total))
	statRow("Deployed live",  fmt.Sprintf("%d / %d", withNetlify, total))

	if total > 0 {
		pct := float64(completed) / float64(total)
		fmt.Printf("\n  %s %s  %s\n",
			ui.Dim.Render("Completion:"),
			miniBar(pct, 30),
			ui.Green.Render(fmt.Sprintf("%.0f%%", pct*100)),
		)
	}

	// Stack breakdown
	fmt.Printf("\n  %s\n", ui.Bold.Render("By Stack"))
	ui.Divider()
	for _, stack := range []string{"vanilla", "tailwind", "react"} {
		count := stacks[stack]
		if count == 0 {
			continue
		}
		bar := miniBar(float64(count)/float64(total), 20)
		fmt.Printf("  %-12s %s  %d\n", ui.Cyan.Render(stack), bar, count)
	}

	// Monthly activity
	months := monthlyBreakdown(all, 6)
	maxCount := 0
	for _, mv := range months {
		if mv.count > maxCount {
			maxCount = mv.count
		}
	}

	fmt.Printf("\n  %s\n", ui.Bold.Render("Monthly Activity  (last 6 months)"))
	ui.Divider()
	for _, mv := range months {
		barLen := 0
		if maxCount > 0 {
			barLen = int(float64(mv.count) / float64(maxCount) * 20)
		}
		bar := strings.Repeat(ui.Cyan.Render("▪"), barLen)
		if barLen == 0 {
			bar = ui.Dim.Render("·")
		}
		fmt.Printf("  %-8s %s  %d\n", ui.Dim.Render(mv.label), bar, mv.count)
	}

	fmt.Println()
	return nil
}

func statRow(label, value string) {
	fmt.Printf("  %-20s %s\n",
		ui.Dim.Render(label),
		lipgloss.NewStyle().Bold(true).Render(value))
}

func miniBar(pct float64, width int) string {
	filled := int(pct * float64(width))
	if filled > width {
		filled = width
	}
	return ui.Green.Render(strings.Repeat("█", filled)) +
		ui.Dim.Render(strings.Repeat("░", width-filled))
}

type monthVal struct {
	label string
	count int
}

func monthlyBreakdown(projects []registry.Project, numMonths int) []monthVal {
	result := make([]monthVal, numMonths)
	now := time.Now()
	for i := numMonths - 1; i >= 0; i-- {
		t := now.AddDate(0, -(numMonths-1-i), 0)
		ym := t.Format("2006-01")
		mv := monthVal{label: t.Format("Jan 06")}
		for _, p := range projects {
			if len(p.Created) >= 7 && p.Created[:7] == ym {
				mv.count++
			}
		}
		result[i] = mv
	}
	return result
}
