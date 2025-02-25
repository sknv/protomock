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
	Port     int    `yaml:"port" envconfig:"HTTP_SERVER_PORT"`
	MocksDir string `yaml:"mocksdir" envconfig:"HTTP_SERVER_MOCKSDIR"`
}

type GRPCServerConfig struct {
	Enabled  bool   `yaml:"enabled" envconfig:"GRPC_SERVER_ENABLED"`
	Port     int    `yaml:"port" envconfig:"GRPC_SERVER_PORT"`
	MocksDir string `yaml:"mocksdir" envconfig:"GRPC_SERVER_MOCKSDIR"`
}

type Config struct {
	Log        LogConfig        `yaml:"log"`
	HTTPServer HTTPServerConfig `yaml:"httpserver"`
	GRPCServer GRPCServerConfig `yaml:"grpcserver"`
}

func Parse(filePath string) (*Config, error) {
	var cfg Config

	if err := config.Parse(filePath, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
