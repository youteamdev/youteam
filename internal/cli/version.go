package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func newVersionCommand(version string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print the CLI version.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			_, err := fmt.Fprintln(outputWriter(cmd), version)
			return err
		},
	}
}
