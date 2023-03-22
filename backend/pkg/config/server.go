package config

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Namespace   string `yaml:"namespace" env:"SERVER_NAMESPACE,overwrite"`
	Name        string `yaml:"name" env:"SERVER_NAME,overwrite"`
	Version     int    `yaml:"version" env:"SERVER_VERSION,overwrite"`
	Address     string `yaml:"address" env:"SERVER_ADDRESS,overwrite"`
	ReplAddress string `yaml:"repl_address" env:"REPL_ADDRESS,overwrite"`
	Debug       bool   `yaml:"debug" env:"SERVER_DEBUG,overwrite"`
}

func (hs *ServerConfig) Validate() error {
	hs.Namespace = strings.TrimSpace(hs.Namespace)
	hs.Name = strings.TrimSpace(hs.Name)
	hs.Address = strings.TrimSpace(hs.Address)
	hs.ReplAddress = strings.TrimSpace(hs.ReplAddress)

	if hs.Namespace == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Namespace",
			Reason:    "Should not be empty",
		}
	}

	if hs.Name == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Name",
			Reason:    "Should not be empty",
		}
	}

	if hs.Address == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Address",
			Reason:    "Should not be empty",
		}
	}

	if hs.ReplAddress == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "Repl Address",
			Reason:    "Should not be empty",
		}
	}

	return nil
}

func BuildNewServerConfig(path string) func() (*ServerConfig, error) {
	return func() (*ServerConfig, error) {
		var config ServerConfig
		if path != "" {
			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			decoder := yaml.NewDecoder(file)

			if err := decoder.Decode(&config); err != nil {
				return nil, err
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		if err := envconfig.Process(ctx, &config); err != nil {
			return nil, err
		}

		return &config, config.Validate()
	}
}