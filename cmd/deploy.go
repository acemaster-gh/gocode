package cmd

import (
	"errors"
	"fmt"

	"github.com/acemaster-gh/gocode/internal/config"
	gitpkg "github.com/acemaster-gh/gocode/internal/git"
	"github.com/acemaster-gh/gocode/internal/netlify"
	"github.com/acemaster-gh/gocode/internal/obsidian"
	"github.com/acemaster-gh/gocode/internal/plugin"
	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a project to Netlify",
	RunE:  runDeploy,
}

func runDeploy(cmd *cobra.Command, args []string) error {
	cfg, err := mustLoadConfig()
	if err != nil {
		return err
	}

	ui.Banner(Version)

	if !netlify.AuthStatus() {
		ui.Error("Not logged in to Netlify. Run: netlify login")
		return nil
	}

	reg := getRegistry()

	name, err := pickActiveProject("Deploy which project?")
	if err != nil || name == "" {
		return err
	}

	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	// Auto-push uncommitted changes first
	if !gitpkg.IsClean(proj.Path) {
		var autoPush bool
		huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("You have uncommitted changes. Push first?").
					Value(&autoPush),
			),
		).WithTheme(huh.ThemeCharm()).Run()

		if autoPush {
			sp := ui.StartSpinner("Pushing changes...")
			gitpkg.CommitAndPush(proj.Path, "chore: pre-deploy update")
			sp.Stop()
		}
	}

	// Deploy
	sp := ui.StartSpinner("Deploying to Netlify...")

	var liveURL string
	var deployErr error

	if proj.NetlifyURL != "" && proj.NetlifyURL != "not-deployed" {
		// Re-deploy existing site
		liveURL, deployErr = netlify.Redeploy(proj.Path)
	} else {
		// Full new deploy
		liveURL = netlify.Full(name, proj.Path, proj.RepoURL)
		if liveURL == "not-deployed" {
			deployErr = fmt.Errorf("deploy failed")
		}
	}

	sp.Stop()

	if deployErr != nil {
		ui.Error(fmt.Sprintf("Deploy failed: %v", deployErr))
		return nil
	}

	if liveURL == "" || liveURL == "not-deployed" {
		ui.Warn("Could not determine live URL. Check: netlify status")
		return nil
	}

	// Update registry
	reg.UpdateField(name, "netlify_url", liveURL)

	// Show result
	ui.Divider()
	ui.OK(fmt.Sprintf("Deployed: %s", ui.Cyan.Copy().Bold(true).Render(liveURL)))
	ui.Divider()

	// Open in browser
	var openBrowser bool
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Open in browser?").
				Value(&openBrowser),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err == nil && openBrowser {
		openURL(liveURL)
	}

	// Ignore ErrUserAborted — treat it as "no"
	if err != nil && !errors.Is(err, huh.ErrUserAborted) {
		return err
	}

	// Background tasks
	vault := obsidian.New(cfg.ObsidianVault)
	go vault.DailyLog(fmt.Sprintf("Deployed [%s]: %s", name, liveURL))
	go plugin.RunHooks(config.PluginsDir(), plugin.HookPostDeploy, plugin.Payload{
		Hook:    plugin.HookPostDeploy,
		Project: name,
		Data:    map[string]string{"url": liveURL},
	})

	return nil
}
