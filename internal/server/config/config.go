package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gox/pkg/fs"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Auth  map[string]string `yaml:"auth"`
	Socks Socks             `yaml:"socks"`
	Https Https             `yaml:"https"`
}

type Socks struct {
	Listen string `yaml:"listen"`
	Enable *bool  `yaml:"enable"`
}

type Https struct {
	Listen string `yaml:"listen"`
	Enable *bool  `yaml:"enable"`
}

func New(filename string, username, password string) (*Config, error) {
	config := defaultConfig()

	data, err := os.ReadFile(filename)
	if err != nil {
		slog.Info("failed to read config, use default config", "path", filename, "error", err)
	} else if err := yaml.Unmarshal(data, &config); err != nil {
		config = defaultConfig()
		slog.Info("failed to unmarshal config, use default config", "path", filename, "error", err)
	}

	normalize(&config)

	if username != "" && password != "" {
		config.Auth = map[string]string{
			username: password,
		}
	}

	slog.Info(
		"config ready",
		"path", filename,
		"auth_enabled", len(config.Auth) != 0,
		"socks_listen", config.Socks.Listen,
		"https_listen", config.Https.Listen,
	)

	return &config, nil
}

func Default(path string) error {
	return save(defaultConfig(), path)
}

func save(config any, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}
	if err := fs.WriteFile(path, data); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}
	return nil
}

func defaultConfig() Config {
	return Config{
		Auth: map[string]string{},
		Socks: Socks{
			Enable: new(true),
			Listen: "0.0.0.0:31080",
		},
		Https: Https{
			Enable: new(true),
			Listen: "0.0.0.0:38443",
		},
	}
}

func normalize(config *Config) {
	if config.Auth == nil {
		config.Auth = map[string]string{}
	}
	if strings.TrimSpace(config.Socks.Listen) == "" {
		config.Socks.Listen = "0.0.0.0:31080"
	}
	if strings.TrimSpace(config.Https.Listen) == "" {
		config.Https.Listen = "0.0.0.0:38443"
	}
	if config.Socks.Enable == nil {
		config.Socks.Enable = new(true)
	}
	if config.Https.Enable == nil {
		config.Https.Enable = new(true)
	}
}

//go:fix inline
func boolPtr(value bool) *bool {
	return new(value)
}
