package config

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v2"
)

type OAuthCredentialsConfig struct {
	Credentials struct {
		ClientID     string `yaml:"client_id" env:"CLIENT_ID,overwrite"`
		ClientSecret string `yaml:"client_secret" env:"CLIENT_SECRET,overwrite"`
		RedirectURI  string `yaml:"redirect_uri" env:"REDIRECT_URI,overwrite"`
	} `yaml:"oauth"`
}

func (zc *OAuthCredentialsConfig) Validate() error {
	zc.Credentials.ClientID = strings.TrimSpace(zc.Credentials.ClientID)
	zc.Credentials.ClientSecret = strings.TrimSpace(zc.Credentials.ClientSecret)

	if zc.Credentials.ClientID == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "ClientID",
			Reason:    "Should not be empty",
		}
	}

	if zc.Credentials.ClientSecret == "" {
		return &InvalidConfigurationParameterError{
			Parameter: "ClientSecret",
			Reason:    "Should not be empty",
		}
	}

	return nil
}

func BuildNewCredentialsConfig(path string) func() (*OAuthCredentialsConfig, error) {
	return func() (*OAuthCredentialsConfig, error) {
		var config OAuthCredentialsConfig
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