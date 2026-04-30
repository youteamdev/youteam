package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestMainRunsVersionCommand(t *testing.T) {
	originalArgs := os.Args
	originalStdout := os.Stdout
	t.Cleanup(func() {
		os.Args = originalArgs
		os.Stdout = originalStdout
	})

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	t.Cleanup(func() {
		_ = reader.Close()
		_ = writer.Close()
	})

	os.Args = []string{"youteam", "version"}
	os.Stdout = writer

	main()

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var stdout bytes.Buffer
	if _, err := io.Copy(&stdout, reader); err != nil {
		t.Fatalf("Copy() error = %v", err)
	}
	if stdout.String() != "dev\n" {
		t.Fatalf("stdout = %q, want %q", stdout.String(), "dev\n")
	}
}
