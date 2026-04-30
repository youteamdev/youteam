package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"
	"youteam.dev/youteam/allinone/internal/config"
	"youteam.dev/youteam/allinone/internal/server"
)

// NewCommand builds the root youteam CLI command.
func NewCommand(version string) *cli.Command {
	return &cli.Command{
		Name:        "youteam",
		Usage:       "Run YouTeam client and server commands.",
		UsageText:   "youteam [global options] command [command options]",
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Usage:     "Load configuration from `path`.",
				TakesFile: true,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging.",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Print the CLI version.",
				Action: func(_ context.Context, cmd *cli.Command) error {
					_, err := fmt.Fprintln(outputWriter(cmd), version)
					return err
				},
			},
			{
				Name:  "client",
				Usage: "Run client commands.",
				Action: func(_ context.Context, cmd *cli.Command) error {
					return cli.ShowSubcommandHelp(cmd)
				},
				Commands: []*cli.Command{
					{
						Name:  "join",
						Usage: "Join a YouTeam server.",
						Action: func(_ context.Context, cmd *cli.Command) error {
							_, err := fmt.Fprintln(outputWriter(cmd), "client join is not implemented yet.")
							return err
						},
					},
					{
						Name:  "info",
						Usage: "Show client information.",
						Action: func(_ context.Context, cmd *cli.Command) error {
							_, err := fmt.Fprintln(outputWriter(cmd), "client info is not implemented yet.")
							return err
						},
					},
				},
			},
			{
				Name:  "server",
				Usage: "Run server commands.",
				Action: func(_ context.Context, cmd *cli.Command) error {
					return cli.ShowSubcommandHelp(cmd)
				},
				Commands: []*cli.Command{
					{
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
					},
				},
			},
		},
	}
}

type debugLogger struct {
	enabled bool
	writer  io.Writer
}

func newDebugLogger(enabled bool, writer io.Writer) debugLogger {
	if writer == nil {
		writer = os.Stderr
	}

	return debugLogger{enabled: enabled, writer: writer}
}

func (l debugLogger) Debugf(format string, args ...any) {
	if !l.enabled {
		return
	}

	_, _ = fmt.Fprintf(l.writer, "debug: "+format+"\n", args...)
}

func outputWriter(cmd *cli.Command) io.Writer {
	if cmd != nil && cmd.Root().Writer != nil {
		return cmd.Root().Writer
	}

	return os.Stdout
}

func errorWriter(cmd *cli.Command) io.Writer {
	if cmd != nil && cmd.Root().ErrWriter != nil {
		return cmd.Root().ErrWriter
	}

	return os.Stderr
}
