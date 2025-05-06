package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	GRPCPort           int    `yaml:"grpc_port"`
	DatabaseURL        string `yaml:"database_url"`
	ExternalServiceURL string `yaml:"external_service_url"`
	MigrationsPath     string `yaml:"migrations_path"`
	LogsPath           string `yaml:"logs_path"`
}

// LoadConfig читает файл config.example.yaml из internal/config
func LoadConfig() (*Config, error) {
	path := filepath.Join("internal", "config", "config.example.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse YAML: %w", err)
	}
	return &cfg, nil
}
