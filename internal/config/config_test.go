package config

import (
	"os"
	"path/filepath"
	"testing"
)

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
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", "18080")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Host != "127.0.0.1" {
		t.Fatalf("Host = %q, want 127.0.0.1", cfg.Host)
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

	if err := os.WriteFile(filepath.Join(dir, defaultConfigFile), []byte("HOST = '127.0.0.1'\nPORT = 18081\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Host != "127.0.0.1" {
		t.Fatalf("Host = %q, want 127.0.0.1", cfg.Host)
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

	if err := os.WriteFile(path, []byte("HOST = '127.0.0.1'\nPORT = '18082'\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Host != "127.0.0.1" {
		t.Fatalf("Host = %q, want 127.0.0.1", cfg.Host)
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
