package cli

import (
	"context"

	"github.com/urfave/cli/v3"
	"youteam.dev/youteam/allinone/internal/config"
	"youteam.dev/youteam/allinone/internal/server"
)

func newServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Run server commands.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			return cli.ShowSubcommandHelp(cmd)
		},
		Commands: []*cli.Command{
			newServerRunCommand(),
		},
	}
}

func newServerRunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Start the embedded HTTP server.",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			logger := newDebugLogger(cmd.Bool("debug"), errorWriter(cmd))

			cfg, err := config.LoadConfig(cmd.String("config"))
			if err != nil {
				return err
			}

			logger.Debugf("config source: %s", cfg.Source)
			logger.Debugf("listen address: %s", cfg.Address())

			return server.RunHTTPServer(ctx, cfg)
		},
	}
}
