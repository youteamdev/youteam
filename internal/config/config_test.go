package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const loopbackHost = "127.0.0.1"

func TestConfigAddressUsesHostPort(t *testing.T) {
	cfg := Config{Host: loopbackHost, Port: "8080"}

	if got := cfg.Address(); got != loopbackHost+":8080" {
		t.Fatalf("Address() = %q, want %q", got, loopbackHost+":8080")
	}
}

func TestLoadConfigDefaultsToEnvironmentWithFallbacks(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("HOST", "")
	t.Setenv("PORT", "")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Host != defaultHost {
		t.Fatalf("Host = %q, want %q", cfg.Host, defaultHost)
	}
	if cfg.Port != defaultPort {
		t.Fatalf("Port = %q, want %q", cfg.Port, defaultPort)
	}
	if cfg.Source != "environment" {
		t.Fatalf("Source = %q, want environment", cfg.Source)
	}
}

func TestLoadConfigUsesEnvironmentWhenNoConfigFileExists(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("HOST", loopbackHost)
	t.Setenv("PORT", "18080")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Host != loopbackHost {
		t.Fatalf("Host = %q, want %q", cfg.Host, loopbackHost)
	}
	if cfg.Port != "18080" {
		t.Fatalf("Port = %q, want 18080", cfg.Port)
	}
}

func TestLoadConfigUsesDefaultFileBeforeEnvironment(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOST", "192.0.2.1")
	t.Setenv("PORT", "19090")

	if err := os.WriteFile(filepath.Join(dir, defaultConfigFile), []byte("HOST = '"+loopbackHost+"'\nPORT = 18081\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Host != loopbackHost {
		t.Fatalf("Host = %q, want %q", cfg.Host, loopbackHost)
	}
	if cfg.Port != "18081" {
		t.Fatalf("Port = %q, want 18081", cfg.Port)
	}
	if cfg.Source != defaultConfigFile {
		t.Fatalf("Source = %q, want %q", cfg.Source, defaultConfigFile)
	}
}

func TestLoadConfigUsesExplicitConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.toml")
	t.Setenv("HOST", "192.0.2.2")
	t.Setenv("PORT", "19091")

	if err := os.WriteFile(path, []byte("HOST = '"+loopbackHost+"'\nPORT = '18082'\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Host != loopbackHost {
		t.Fatalf("Host = %q, want %q", cfg.Host, loopbackHost)
	}
	if cfg.Port != "18082" {
		t.Fatalf("Port = %q, want 18082", cfg.Port)
	}
	if cfg.Source != path {
		t.Fatalf("Source = %q, want %q", cfg.Source, path)
	}
}

func TestLoadConfigReturnsErrorForMissingExplicitConfig(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.toml"))
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want error")
	}
}

func TestLoadConfigReturnsErrorForInvalidDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile(filepath.Join(dir, defaultConfigFile), []byte("PORT = [\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := LoadConfig("")
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want error")
	}
}

func TestLoadConfigReturnsErrorForInvalidEnvironmentPort(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "0")

	_, err := LoadConfig("")
	if err == nil {
		t.Fatal("LoadConfig() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "invalid PORT from environment") {
		t.Fatalf("LoadConfig() error = %v, want invalid environment port error", err)
	}
}

func TestNormalizePortTrimsValidPort(t *testing.T) {
	got, err := normalizePort(" 65535 ")
	if err != nil {
		t.Fatalf("normalizePort() error = %v", err)
	}
	if got != "65535" {
		t.Fatalf("normalizePort() = %q, want %q", got, "65535")
	}
}
