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

package handler

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/adapter"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/service"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/stretchr/testify/assert"
)

type mockEncryptor struct{}

func (e mockEncryptor) Encrypt(text string) (string, error) {
	return string(text), nil
}

func (e mockEncryptor) Decrypt(ciphertext string) (string, error) {
	return string(ciphertext), nil
}

func TestSelectCaching(t *testing.T) {
	adapter := adapter.NewMemoryUserAdapter()
	service := service.NewUserService(adapter, mockEncryptor{}, log.NewEmptyLogger())
	pclient := pclient.NewPipedriveAuthClient("clientID", "clientSecret")

	sel := NewUserSelectHandler(service, nil, pclient, log.NewEmptyLogger())

	service.CreateUser(context.Background(), domain.UserAccess{
		ID:           "mock",
		AccessToken:  "mock",
		RefreshToken: "mock",
		TokenType:    "mock",
		Scope:        "mock",
		ExpiresAt:    time.Now().Add(24 * time.Hour).UnixMilli(),
	})

	t.Run("get user", func(t *testing.T) {
		var res domain.UserAccess
		id := "mock"
		assert.NoError(t, sel.GetUser(context.Background(), &id, &res))
		assert.NotEmpty(t, res)
	})
}
