package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HookType = string

const (
	HookPostCreate   HookType = "post-create"
	HookPrePush      HookType = "pre-push"
	HookPostPush     HookType = "post-push"
	HookPostDeploy   HookType = "post-deploy"
	HookPostComplete HookType = "post-complete"
)

type Payload struct {
	Hook    HookType          `json:"hook"`
	Project string            `json:"project"`
	Data    map[string]string `json:"data"`
}

// RunHooks finds all executables in pluginsDir and calls each with JSON payload via stdin.
// Errors from individual plugins are silently ignored.
func RunHooks(pluginsDir string, hook HookType, p Payload) {
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return // plugins dir doesn't exist — fine
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(pluginsDir, e.Name())
		info, err := e.Info()
		if err != nil {
			continue
		}
		// Check executable bit
		if info.Mode()&0111 == 0 {
			continue
		}

		cmd := exec.Command(path)
		cmd.Stdin = strings.NewReader(string(payload))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run() // ignore errors from plugins
	}
}

// PrintPayload pretty-prints a payload for debugging.
func PrintPayload(p Payload) {
	data, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(data))
}
