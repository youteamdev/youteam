// Package main wires the youteam command-line application.
package main

import (
	"context"
	"fmt"
	"os"

	"youteam.dev/youteam/allinone/internal/cli"
)

var (
	version = defaultVersion
	commit  = defaultCommit
)

func main() {
	cmd := cli.NewCommand(displayVersion(version, commit))
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
