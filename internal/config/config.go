package config

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultConfigFile = "youteam.toml"
	defaultHost       = "0.0.0.0"
	defaultPort       = "8080"
)

// Config contains the runtime settings used by the embedded server.
type Config struct {
	Host   string
	Port   string
	Source string
}

// Address returns the host:port value consumed by net/http.
func (c Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

// LoadConfig resolves configuration from an explicit TOML file, the default
// local TOML file, or environment variables.
func LoadConfig(path string) (Config, error) {
	if path != "" {
		return loadConfigFile(path, path)
	}

	cfg, err := loadDefaultConfigFile()
	if err == nil {
		return cfg, nil
	}

	var notFound viper.ConfigFileNotFoundError
	if !errors.As(err, &notFound) {
		return Config{}, err
	}

	return loadConfigFromEnvironment()
}

func defaultConfig(source string) Config {
	return Config{
		Host:   defaultHost,
		Port:   defaultPort,
		Source: source,
	}
}

func loadDefaultConfigFile() (Config, error) {
	v := newConfigViper()
	v.SetConfigName(strings.TrimSuffix(defaultConfigFile, ".toml"))
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	return configFromViper(v, defaultConfigFile)
}

func loadConfigFile(path string, source string) (Config, error) {
	v := newConfigViper()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("load config %q: %w", path, err)
	}

	return configFromViper(v, source)
}

func loadConfigFromEnvironment() (Config, error) {
	v := newConfigViper()
	v.AutomaticEnv()

	return configFromViper(v, "environment")
}

func newConfigViper() *viper.Viper {
	v := viper.New()
	v.SetDefault("HOST", defaultHost)
	v.SetDefault("PORT", defaultPort)
	return v
}

func configFromViper(v *viper.Viper, source string) (Config, error) {
	cfg := defaultConfig(source)
	if host := strings.TrimSpace(v.GetString("HOST")); host != "" {
		cfg.Host = host
	}

	if port := strings.TrimSpace(v.GetString("PORT")); port != "" {
		normalized, err := normalizePort(port)
		if err != nil {
			return Config{}, fmt.Errorf("invalid PORT from %s: %w", source, err)
		}
		cfg.Port = normalized
	}

	return cfg, nil
}

func normalizePort(value string) (string, error) {
	value = strings.TrimSpace(value)
	port, err := strconv.Atoi(value)
	if err != nil {
		return "", fmt.Errorf("%q is not a valid port", value)
	}
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("%q is outside the valid port range", value)
	}

	return strconv.Itoa(port), nil
}
