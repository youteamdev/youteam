package cli

import (
	"github.com/urfave/cli/v3"
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
			newVersionCommand(version),
			newClientCommand(),
			newServerCommand(),
		},
	}
}
