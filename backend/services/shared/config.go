/**
 *
 * (c) Copyright Ascensio System SIA 2023
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package shared

import (
	"context"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v2"
)

type OnlyofficeConfig struct {
	Onlyoffice struct {
		Builder  OnlyofficeBuilderConfig  `yaml:"builder"`
		Callback OnlyofficeCallbackConfig `yaml:"callback"`
	} `yaml:"onlyoffice"`
}

func (oc *OnlyofficeConfig) Validate() error {
	if err := oc.Onlyoffice.Builder.Validate(); err != nil {
		return err
	}

	return oc.Onlyoffice.Callback.Validate()
}

func BuildNewOnlyofficeConfig(path string) func() (*OnlyofficeConfig, error) {
	return func() (*OnlyofficeConfig, error) {
		var config OnlyofficeConfig
		config.Onlyoffice.Callback.MaxSize = 20000000
		config.Onlyoffice.Callback.UploadTimeout = 120
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

type OnlyofficeBuilderConfig struct {
	GatewayURL       string `yaml:"gateway_url" env:"ONLYOFFICE_GATEWAY_URL,overwrite"`
	CallbackURL      string `yaml:"callback_url" env:"ONLYOFFICE_CALLBACK_URL,overwrite"`
	AllowedDownloads int    `yaml:"allowed_downloads" env:"ONLYOFFICE_ALLOWED_DOWNLOADS,overwrite"`
}

func (oc *OnlyofficeBuilderConfig) Validate() error {
	return nil
}

type OnlyofficeCallbackConfig struct {
	MaxSize       int64 `yaml:"max_size" env:"ONLYOFFICE_CALLBACK_MAX_SIZE,overwrite"`
	UploadTimeout int   `yaml:"upload_timeout" env:"ONLYOFFICE_CALLBACK_UPLOAD_TIMEOUT,overwrite"`
}

func (c *OnlyofficeCallbackConfig) Validate() error {
	return nil
}
