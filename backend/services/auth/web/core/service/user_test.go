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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/cache"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/config"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/domain"
	"github.com/stretchr/testify/assert"
)

type mockEncryptor struct{}

func (e mockEncryptor) Encrypt(text string) (string, error) {
	return string(text), nil
}

func (e mockEncryptor) Decrypt(ciphertext string) (string, error) {
	return string(ciphertext), nil
}

type mockAdapter struct {
}

var user = domain.UserAccess{
	ID:           "mock",
	AccessToken:  "mock",
	RefreshToken: "mock",
	TokenType:    "mock",
	Scope:        "mock",
	ExpiresAt:    1000000,
	ApiDomain:    "pipedrive",
}

func (m mockAdapter) InsertUser(ctx context.Context, user domain.UserAccess) error {
	return nil
}

func (m mockAdapter) SelectUserByID(ctx context.Context, uid string) (domain.UserAccess, error) {
	return user, nil
}

func (m mockAdapter) UpsertUser(ctx context.Context, user domain.UserAccess) (domain.UserAccess, error) {
	return domain.UserAccess{
		ID:          "mock",
		AccessToken: "mock",
	}, nil
}

func (m mockAdapter) DeleteUserByID(ctx context.Context, uid string) error {
	return nil
}

func TestUserService(t *testing.T) {
	service := NewUserService(mockAdapter{}, mockEncryptor{}, cache.NewCache(&config.CacheConfig{}), log.NewEmptyLogger())

	t.Run("save user", func(t *testing.T) {
		assert.NoError(t, service.CreateUser(context.Background(), user))
	})

	t.Run("get user", func(t *testing.T) {
		u, err := service.GetUser(context.Background(), "mock")
		assert.NoError(t, err)
		assert.Equal(t, user, u)
	})

	t.Run("get user with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		_, err := service.GetUser(ctx, "mock")
		assert.Error(t, err)
	})

	t.Run("update user token", func(t *testing.T) {
		_, err := service.UpdateUser(context.Background(), domain.UserAccess{
			ID:           "mock",
			AccessToken:  "mock",
			RefreshToken: "mock",
			TokenType:    "mock",
			Scope:        "mock",
			ExpiresAt:    100000,
			ApiDomain:    "pipedrive",
		})
		assert.NoError(t, err)
	})

	t.Run("update user token with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		_, err := service.UpdateUser(ctx, domain.UserAccess{
			ID:           "mock",
			AccessToken:  "mock",
			RefreshToken: "mock",
			TokenType:    "mock",
			Scope:        "mock",
			ExpiresAt:    100000,
			ApiDomain:    "pipedrive",
		})
		assert.Error(t, err)
	})

	t.Run("delete user", func(t *testing.T) {
		assert.NoError(t, service.DeleteUser(context.Background(), "mock"))
	})
}
