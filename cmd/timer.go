package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Pomodoro focus timer",
	RunE:  runTimer,
}

func runTimer(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	// Config form
	var duration int
	var label string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Focus duration").
				Options(
					huh.NewOption("⏱  25 min  — standard Pomodoro", 25),
					huh.NewOption("⏱  50 min  — deep work", 50),
					huh.NewOption("⏱  15 min  — short sprint", 15),
					huh.NewOption("⏱  5 min   — quick task", 5),
				).
				Value(&duration),
			huh.NewInput().
				Title("What are you working on?").
				Placeholder("implementing the timer feature").
				Value(&label),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	if label == "" {
		label = "deep work session"
	}

	total := time.Duration(duration) * time.Minute
	end := time.Now().Add(total)

	// Handle Ctrl+C gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	fmt.Println()
	fmt.Printf("  %s\n", ui.Bold.Render(fmt.Sprintf("Focus: %s", label)))
	fmt.Printf("  %s\n\n", ui.Dim.Render(fmt.Sprintf("%d minutes — started at %s",
		duration, time.Now().Format("15:04"))))

	// Progress bar width
	barWidth := 40
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		select {
		case <-sigCh:
			fmt.Fprintf(os.Stderr, "\r\033[K")
			ui.Warn("Timer cancelled.")
			return nil

		case <-tick.C:
			remaining := time.Until(end)
			if remaining <= 0 {
				fmt.Fprintf(os.Stderr, "\r\033[K")
				showTimerDone(label, duration)
				return nil
			}

			elapsed := total - remaining
			progress := float64(elapsed) / float64(total)
			filled := int(progress * float64(barWidth))

			bar := ""
			for i := 0; i < barWidth; i++ {
				if i < filled {
					bar += ui.Green.Render("█")
				} else {
					bar += ui.Dim.Render("░")
				}
			}

			mins := int(remaining.Minutes())
			secs := int(remaining.Seconds()) % 60

			fmt.Fprintf(os.Stderr, "\r  %s  %s  %s",
				bar,
				ui.Cyan.Bold(true).Render(fmt.Sprintf("%02d:%02d", mins, secs)),
				ui.Dim.Render("remaining"),
			)
		}
	}
}

func showTimerDone(label string, duration int) {
	fmt.Println()
	fmt.Println()

	box := ui.BoxStyle.Render(
		fmt.Sprintf("%s\n\n%s\n%s",
			ui.Green.Bold(true).Render("  ✓ Session complete!"),
			ui.Bold.Render("  "+label),
			ui.Dim.Render(fmt.Sprintf("  %d minutes of focused work ⚡", duration)),
		),
	)
	fmt.Println(box)
	fmt.Println()

	// Terminal bell
	fmt.Print("\a")

	// Suggest break
	ui.Info("Take a 5-minute break. You earned it.")
}
