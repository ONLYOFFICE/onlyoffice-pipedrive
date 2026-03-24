/**
 *
 * (c) Copyright Ascensio System SIA 2026
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

package handler

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/cache"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/adapter"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/service"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

type mockEncryptor struct{}

func (e mockEncryptor) Encrypt(text string, key []byte) (string, error) {
	return string(text), nil
}

func (e mockEncryptor) Decrypt(ciphertext string, key []byte) (string, error) {
	return string(ciphertext), nil
}

func TestSelectCaching(t *testing.T) {
	adapter := adapter.NewMemoryDocserverAdapter()
	service := service.NewSettingsService(
		adapter, mockEncryptor{}, cache.NewCache(&config.CacheConfig{}),
		&oauth2.Config{
			ClientID:     "mock",
			ClientSecret: "mock",
		}, log.NewEmptyLogger(),
	)

	sel := NewSettingsSelectHandler(service, nil, log.NewEmptyLogger())

	service.CreateSettings(context.Background(), domain.DocSettings{
		CompanyID:   "mock",
		DocAddress:  "mock",
		DocSecret:   "mock",
		DocHeader:   "mock",
		DemoEnabled: true,
		DemoStarted: time.Now(),
	})

	t.Run("get settings", func(t *testing.T) {
		var res response.DocSettingsResponse
		id := "mock"
		assert.NoError(t, sel.GetSettings(context.Background(), &id, &res))
		assert.NotEmpty(t, res)
	})
}
