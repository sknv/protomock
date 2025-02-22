package config

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/sknv/protomock/pkg/config"
)

const _defaultConfigFilePath = "./configs/protomock.yaml"

func FilePathFlag() *string {
	return flag.String("c", _defaultConfigFilePath, "configuration file path")
}

// ----------------------------------------------------------------------------

type LogConfig struct {
	Level slog.Level `yaml:"level" envconfig:"LOG_LEVEL"`
}

type HTTPServerConfig struct {
	Enabled  bool   `yaml:"enabled" envconfig:"HTTP_SERVER_ENABLED"`
	Address  string `yaml:"address" envconfig:"HTTP_SERVER_ADDRESS"`
	MocksDir string `yaml:"mocksdir" envconfig:"HTTP_SERVER_MOCKSDIR"`
}

type Config struct {
	Log        LogConfig        `yaml:"log"`
	HTTPServer HTTPServerConfig `yaml:"httpserver"`
}

func Parse(filePath string) (*Config, error) {
	var cfg Config

	if err := config.Parse(filePath, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
