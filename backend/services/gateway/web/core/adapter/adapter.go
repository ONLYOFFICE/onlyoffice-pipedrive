/**
 *
 * (c) Copyright Ascensio System SIA 2025
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

package adapter

import (
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/port"
	"golang.org/x/oauth2"
)

func BuildNewAICodeAccessAdapter(config *config.StorageConfig, encryptor crypto.Encryptor, credentials *oauth2.Config) port.AICodeAccessServiceAdapter {
	adapter := NewMemoryAICodeAccessAdapter()
	if config.Storage.URL != "" {
		adapter = NewMongoAICodeAccessAdapter(config.Storage.URL, encryptor, credentials)
	}

	return adapter
}
