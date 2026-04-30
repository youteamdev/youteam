package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"
)

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
