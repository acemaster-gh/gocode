package main

import (
	"os"
	"strings"

	"github.com/acemaster-gh/gocode/cmd"
)

// Version is set via -ldflags at build time.
var Version = "2.0.0"

func main() {
	// Translate legacy --flag style to cobra subcommand style.
	// e.g. "gocode --new" → "gocode new"
	// Multi-word flags: --open-github → "open-github" (handled in open.go aliases)
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if strings.HasPrefix(arg, "--") && arg != "--help" && arg != "--version" {
			os.Args[1] = strings.TrimPrefix(arg, "--")
		}
	}
	cmd.Execute(Version)
}
