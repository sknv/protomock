package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

func Parse(filePath string, cfg any) error {
	if err := readFile(filePath, cfg); err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if err := readEnv(cfg); err != nil {
		return fmt.Errorf("read env: %w", err)
	}

	return nil
}

func readFile(filePath string, cfg any) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(cfg); err != nil {
		return fmt.Errorf("decode yaml file: %w", err)
	}

	return nil
}

func readEnv(cfg any) error {
	if err := envconfig.Process("", cfg); err != nil {
		return fmt.Errorf("process from env: %w", err)
	}

	return nil
}
