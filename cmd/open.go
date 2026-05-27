package cmd

import (
	"errors"
	"fmt"

	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"open-github", "open-live", "open-netlify", "open-folder"},
	Short:   "Open project URLs or folders",
	RunE:    runOpen,
}

func runOpen(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	name, err := pickAllProject("Open resources for which project?")
	if err != nil || name == "" {
		return err
	}

	reg := getRegistry()
	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	// Build available options
	opts := []huh.Option[string]{}

	opts = append(opts, huh.NewOption("📁  Open folder    open in file manager", "folder"))

	if proj.RepoURL != "local-only" && proj.RepoURL != "" {
		opts = append(opts,
			huh.NewOption("⎇  GitHub repo    open in browser", "github"))
	}

	if proj.NetlifyURL != "not-deployed" && proj.NetlifyURL != "" {
		opts = append(opts,
			huh.NewOption("🌐  Live site      open deployed URL in browser", "live"),
			huh.NewOption("📊  Netlify panel  open netlify.com dashboard", "netlify"))
	}

	if len(opts) == 0 {
		ui.Info("No external links configured for this project.")
		return nil
	}

	var target string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Open: %s", name)).
				Options(opts...).
				Value(&target),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return err
	}

	switch target {
	case "folder":
		openPath(proj.Path)
		ui.OK(fmt.Sprintf("Opened folder: %s", proj.Path))

	case "github":
		openURL(proj.RepoURL)
		ui.OK(fmt.Sprintf("Opening: %s", proj.RepoURL))

	case "live":
		openURL(proj.NetlifyURL)
		ui.OK(fmt.Sprintf("Opening: %s", proj.NetlifyURL))

	case "netlify":
		openURL("https://app.netlify.com")
		ui.OK("Opening Netlify dashboard.")
	}

	reg.UpdateLastOpened(name)
	return nil
}
