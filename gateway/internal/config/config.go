package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Server struct {
		BindAddr string `yaml:"bind_addr"`
	} `yaml:"server"`

	AuthServiceURL     string `yaml:"auth_service_url"`
	CurrencyServiceURL string `yaml:"currency_service_url"`
}

func LoadConfig() (*Config, error) {
	path := filepath.Join("gateway", "internal", "config", "config.example.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse YAML: %w", err)
	}

	if v := os.Getenv("BIND_ADDR"); v != "" {
		cfg.Server.BindAddr = v
	}
	if v := os.Getenv("AUTH_SERVICE_URL"); v != "" {
		cfg.AuthServiceURL = v
	}
	if v := os.Getenv("CURRENCY_SERVICE_URL"); v != "" {
		cfg.CurrencyServiceURL = v
	}

	if cfg.Server.BindAddr == "" {
		cfg.Server.BindAddr = ":8080"
	}

	return &cfg, nil
}
