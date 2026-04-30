package cli

import (
	"bytes"
	"context"
	"net"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

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
	stdout, stderr, err := runCLI(t, NewCommand("1.2.3-abcdef0"), "version")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if stdout != "1.2.3-abcdef0\n" {
		t.Fatalf("stdout = %q, want %q", stdout, "1.2.3-abcdef0\n")
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

func TestServerRunReturnsConfigLoadError(t *testing.T) {
	_, _, err := runCLI(t, NewCommand("test-version"), "--config", filepath.Join(t.TempDir(), "missing.toml"), "server", "run")
	if err == nil {
		t.Fatal("Run() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("Run() error = %v, want load config error", err)
	}
}

func TestServerRunStartsWithDebugLoggingUntilContextCancel(t *testing.T) {
	t.Chdir(t.TempDir())
	port := unusedTCPPort(t)
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type result struct {
		stdout string
		stderr string
		err    error
	}
	done := make(chan result, 1)
	go func() {
		stdout, stderr, err := runCLIResult(ctx, NewCommand("test-version"), "--debug", "server", "run")
		done <- result{stdout: stdout, stderr: stderr, err: err}
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case got := <-done:
		if got.err != nil {
			t.Fatalf("Run() error = %v", got.err)
		}
		if got.stdout != "" {
			t.Fatalf("stdout = %q, want empty", got.stdout)
		}
		for _, want := range []string{"debug: config source: environment", "debug: listen address: 127.0.0.1:" + port} {
			if !strings.Contains(got.stderr, want) {
				t.Fatalf("stderr = %q, want to contain %q", got.stderr, want)
			}
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run() did not stop after context cancellation")
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

func unusedTCPPort(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	defer listener.Close()

	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("SplitHostPort() error = %v", err)
	}
	return port
}

func runCLI(t *testing.T, cmd *urfavecli.Command, args ...string) (string, string, error) {
	t.Helper()

	return runCLIResult(context.Background(), cmd, args...)
}

func runCLIResult(ctx context.Context, cmd *urfavecli.Command, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Writer = &stdout
	cmd.ErrWriter = &stderr

	runArgs := append([]string{cmd.Name}, args...)
	err := cmd.Run(ctx, runArgs)
	return stdout.String(), stderr.String(), err
}
