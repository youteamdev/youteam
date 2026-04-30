package cli

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	urfavecli "github.com/urfave/cli/v3"
)

func TestNewCommandRootMetadata(t *testing.T) {
	cmd := NewCommand("test-version")

	if cmd.Name != "youteam" {
		t.Fatalf("Name = %q, want %q", cmd.Name, "youteam")
	}
	if cmd.Usage != "Run YouTeam client and server commands." {
		t.Fatalf("Usage = %q, want %q", cmd.Usage, "Run YouTeam client and server commands.")
	}
	if cmd.UsageText != "youteam [global options] command [command options]" {
		t.Fatalf("UsageText = %q, want %q", cmd.UsageText, "youteam [global options] command [command options]")
	}
	if !cmd.HideVersion {
		t.Fatal("HideVersion = false, want true")
	}

	assertCommandNames(t, cmd.Commands, []string{"version", "client", "server"})
	assertRootFlags(t, cmd)
}

func TestVersionCommandWritesConfiguredVersion(t *testing.T) {
	stdout, stderr, err := runCLI(t, NewCommand("1.2.3"), "version")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if stdout != "1.2.3\n" {
		t.Fatalf("stdout = %q, want %q", stdout, "1.2.3\n")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestClientCommandsPrintPlaceholderMessages(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "join",
			args: []string{"client", "join"},
			want: "client join is not implemented yet.\n",
		},
		{
			name: "info",
			args: []string{"client", "info"},
			want: "client info is not implemented yet.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCLI(t, NewCommand("test-version"), tt.args...)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if stdout != tt.want {
				t.Fatalf("stdout = %q, want %q", stdout, tt.want)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}
		})
	}
}

func TestCommandGroupsShowSubcommandHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "client",
			args:     []string{"client"},
			contains: []string{"Run client commands.", "join", "info"},
		},
		{
			name:     "server",
			args:     []string{"server"},
			contains: []string{"Run server commands.", "run"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCLI(t, NewCommand("test-version"), tt.args...)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			for _, want := range tt.contains {
				if !strings.Contains(stdout, want) {
					t.Fatalf("stdout = %q, want to contain %q", stdout, want)
				}
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}
		})
	}
}

func TestDebugLoggerWritesOnlyWhenEnabled(t *testing.T) {
	var disabled bytes.Buffer
	newDebugLogger(false, &disabled).Debugf("listen address: %s", "127.0.0.1:8080")
	if disabled.String() != "" {
		t.Fatalf("disabled logger wrote %q, want empty", disabled.String())
	}

	var enabled bytes.Buffer
	newDebugLogger(true, &enabled).Debugf("listen address: %s", "127.0.0.1:8080")
	if enabled.String() != "debug: listen address: 127.0.0.1:8080\n" {
		t.Fatalf("enabled logger wrote %q", enabled.String())
	}
}

func assertRootFlags(t *testing.T, cmd *urfavecli.Command) {
	t.Helper()

	if len(cmd.Flags) != 2 {
		t.Fatalf("len(Flags) = %d, want 2", len(cmd.Flags))
	}

	configFlag, ok := cmd.Flags[0].(*urfavecli.StringFlag)
	if !ok {
		t.Fatalf("Flags[0] = %T, want *cli.StringFlag", cmd.Flags[0])
	}
	if configFlag.Name != "config" {
		t.Fatalf("config flag name = %q, want %q", configFlag.Name, "config")
	}
	if !configFlag.TakesFile {
		t.Fatal("config flag TakesFile = false, want true")
	}

	debugFlag, ok := cmd.Flags[1].(*urfavecli.BoolFlag)
	if !ok {
		t.Fatalf("Flags[1] = %T, want *cli.BoolFlag", cmd.Flags[1])
	}
	if debugFlag.Name != "debug" {
		t.Fatalf("debug flag name = %q, want %q", debugFlag.Name, "debug")
	}
}

func assertCommandNames(t *testing.T, commands []*urfavecli.Command, want []string) {
	t.Helper()

	got := make([]string, 0, len(commands))
	for _, command := range commands {
		got = append(got, command.Name)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("command names = %v, want %v", got, want)
	}
}

func runCLI(t *testing.T, cmd *urfavecli.Command, args ...string) (string, string, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Writer = &stdout
	cmd.ErrWriter = &stderr

	runArgs := append([]string{cmd.Name}, args...)
	err := cmd.Run(context.Background(), runArgs)
	return stdout.String(), stderr.String(), err
}
