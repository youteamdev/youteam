package cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func newClientCommand() *cli.Command {
	return &cli.Command{
		Name:  "client",
		Usage: "Run client commands.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			return cli.ShowSubcommandHelp(cmd)
		},
		Commands: []*cli.Command{
			newClientJoinCommand(),
			newClientInfoCommand(),
		},
	}
}

func newClientJoinCommand() *cli.Command {
	return &cli.Command{
		Name:  "join",
		Usage: "Join a YouTeam server.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			_, err := fmt.Fprintln(outputWriter(cmd), "client join is not implemented yet.")
			return err
		},
	}
}

func newClientInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "info",
		Usage: "Show client information.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			_, err := fmt.Fprintln(outputWriter(cmd), "client info is not implemented yet.")
			return err
		},
	}
}
