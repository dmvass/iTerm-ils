package main

import (
	"fmt"
	"os"
)

// themePath is path to a theme directory
// could be changed at a compiletion time
var themePath string

func main() {
	if themePath == "" {
		themePath = "theme"
	}

	theme, err := NewTheme(themePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmd, err := NewCommand(theme, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
