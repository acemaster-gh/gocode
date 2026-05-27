package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ── Styles ────────────────────────────────────────────────────────────────────

var (
	Cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	Yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("221"))
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("183"))
	Dim     = lipgloss.NewStyle().Faint(true)
	Bold    = lipgloss.NewStyle().Bold(true)

	TagInfo  = Cyan.Copy().Bold(true).Render("[INFO] ")
	TagOK    = Green.Copy().Bold(true).Render("[OK]   ")
	TagWarn  = Yellow.Copy().Bold(true).Render("[WARN] ")
	TagError = Red.Copy().Bold(true).Render("[ERR]  ")
	TagStep  = Cyan.Copy().Bold(true).Render("[....] ")
	TagDone  = Green.Copy().Bold(true).Render("[DONE] ")

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("114")).
			Padding(0, 1)

	SummaryTitle = Green.Copy().Bold(true)
	SummaryDim   = Dim.Copy()
)

// ── Banner ────────────────────────────────────────────────────────────────────

func Banner(version string) {
	cols := termWidth()
	fmt.Println()
	if cols >= 60 {
		c := Cyan.Bold(true)
		fmt.Println(c.Render("  ██████╗  ██████╗  ██████╗ ██████╗ ██████╗ ███████╗"))
		fmt.Println(c.Render(" ██╔════╝ ██╔═══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝"))
		fmt.Println(c.Render(" ██║  ███╗██║   ██║██║     ██║   ██║██║  ██║█████╗  "))
		fmt.Println(c.Render(" ██║   ██║██║   ██║██║     ██║   ██║██║  ██║██╔══╝  "))
		fmt.Println(c.Render(" ╚██████╔╝╚██████╔╝╚██████╗╚██████╔╝██████╔╝███████╗"))
		fmt.Println(c.Render("  ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝"))
		fmt.Println(Dim.Render("  developer workflow system · v" + version))
	} else {
		fmt.Println(Cyan.Bold(true).Render("  gocode v" + version))
	}
	fmt.Println()
}

func termWidth() int {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 80
	}
	var w int
	fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &w)
	if w < 20 {
		return 80
	}
	return w
}

// ── Log functions (write to stderr) ──────────────────────────────────────────

func Info(msg string)  { fmt.Fprintln(os.Stderr, TagInfo+msg) }
func OK(msg string)    { fmt.Fprintln(os.Stderr, TagOK+msg) }
func Warn(msg string)  { fmt.Fprintln(os.Stderr, TagWarn+msg) }
func Error(msg string) { fmt.Fprintln(os.Stderr, TagError+msg) }
func Step(msg string)  { fmt.Fprintln(os.Stderr, TagStep+Bold.Render(msg)) }
func Done(msg string)  { fmt.Fprintln(os.Stderr, TagDone+Bold.Render(msg)) }

func Divider() {
	fmt.Fprintln(os.Stderr, Dim.Render("────────────────────────────────────────────────────"))
}

// ── Summary box ───────────────────────────────────────────────────────────────

func SummaryBox(name, stack, path, repo, netlify string) {
	lines := []string{
		SummaryTitle.Render("✓ Project Ready: " + name),
		"",
		SummaryDim.Render("  Stack:   ") + stack,
		SummaryDim.Render("  Path:    ") + path,
	}
	if repo != "local-only" && repo != "" {
		lines = append(lines, SummaryDim.Render("  GitHub:  ")+repo)
	}
	if netlify != "not-deployed" && netlify != "" {
		lines = append(lines, SummaryDim.Render("  Live:    ")+netlify)
	}
	fmt.Println()
	fmt.Println(BoxStyle.Render(strings.Join(lines, "\n")))
	fmt.Println()
}

// ── Spinner ────────────────────────────────────────────────────────────────────

type Spinner struct {
	msg  string
	stop chan struct{}
	done chan struct{}
	mu   sync.Mutex
}

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func StartSpinner(msg string) *Spinner {
	s := &Spinner{
		msg:  msg,
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Fprintf(os.Stderr, "\r\033[K")
				close(s.done)
				return
			default:
				s.mu.Lock()
				m := s.msg
				s.mu.Unlock()
				fmt.Fprintf(os.Stderr, "\r%s %s", Cyan.Render(frames[i]), m)
				i = (i + 1) % len(frames)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return s
}

func (s *Spinner) UpdateMsg(msg string) {
	s.mu.Lock()
	s.msg = msg
	s.mu.Unlock()
}

func (s *Spinner) Stop() {
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
	<-s.done
}

// ── Confirm ───────────────────────────────────────────────────────────────────

func Confirm(prompt string) bool {
	fmt.Fprintf(os.Stderr, "%s %s [y/N]: ",
		Yellow.Render("[?]"), prompt)
	var ans string
	fmt.Scanln(&ans)
	return strings.ToLower(strings.TrimSpace(ans)) == "y"
}
