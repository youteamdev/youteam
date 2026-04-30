package main

import (
	"context"
	"fmt"
	"os"

	"youteam.dev/youteam/allinone/internal/cli"
)

var version = "dev"

func main() {
	cmd := cli.NewCommand(version)
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
