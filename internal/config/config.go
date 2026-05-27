package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	GHUser             string   `mapstructure:"gh_user"`
	ProjectsDir        string   `mapstructure:"projects_dir"`
	ObsidianVault      string   `mapstructure:"obsidian_vault"`
	BrowserTabs        []string `mapstructure:"browser_tabs"`
	SleepBeforeBrowser int      `mapstructure:"sleep_before_browser"`
}

var configDir string

func init() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".gocode")
}

func ConfigPath() string { return filepath.Join(configDir, "config.toml") }
func GocodeDir() string  { return configDir }
func DataDir() string    { return filepath.Join(configDir, "data") }
func PluginsDir() string { return filepath.Join(configDir, "plugins") }

var ErrNotFound = errors.New("config not found — run: gocode config")

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(configDir)
	v.SetDefault("sleep_before_browser", 5)
	v.SetDefault("browser_tabs", []string{
		"https://claude.ai",
		"https://chatgpt.com",
		"https://developer.mozilla.org",
		"http://localhost:5500",
	})

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) || os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	tabs := ""
	for i, tab := range cfg.BrowserTabs {
		if i > 0 {
			tabs += ", "
		}
		tabs += fmt.Sprintf("%q", tab)
	}
	content := fmt.Sprintf(
		"# gocode v2.0.0 config\ngh_user             = %q\nprojects_dir        = %q\nobsidian_vault      = %q\nsleep_before_browser = %d\nbrowser_tabs        = [%s]\n",
		cfg.GHUser, cfg.ProjectsDir, cfg.ObsidianVault, cfg.SleepBeforeBrowser, tabs,
	)
	return os.WriteFile(ConfigPath(), []byte(content), 0644)
}
