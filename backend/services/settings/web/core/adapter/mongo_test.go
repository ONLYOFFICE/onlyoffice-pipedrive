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

package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/stretchr/testify/assert"
)

var settings = domain.DocSettings{
	CompanyID:  "mock",
	DocAddress: "mock",
	DocSecret:  "mock",
}

func TestMongoAdapter(t *testing.T) {
	adapter := NewMongoDocserverAdapter("mongodb://localhost:27017")

	t.Run("save settings with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		assert.Error(t, adapter.InsertSettings(ctx, settings))
	})

	t.Run("save settings", func(t *testing.T) {
		assert.NoError(t, adapter.InsertSettings(context.Background(), settings))
	})

	t.Run("save the same settings object", func(t *testing.T) {
		assert.NoError(t, adapter.InsertSettings(context.Background(), settings))
	})

	t.Run("get settings by id with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		_, err := adapter.SelectSettings(ctx, "mock")
		assert.Error(t, err)
	})

	t.Run("get settings by id", func(t *testing.T) {
		s, err := adapter.SelectSettings(context.Background(), "mock")
		assert.NoError(t, err)
		assert.Equal(t, settings, s)
	})

	t.Run("delete settings by id with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		assert.Error(t, adapter.DeleteSettings(ctx, "mock"))
	})

	t.Run("delete settings by id", func(t *testing.T) {
		assert.NoError(t, adapter.DeleteSettings(context.Background(), "mock"))
	})

	t.Run("get invalid settings", func(t *testing.T) {
		_, err := adapter.SelectSettings(context.Background(), "mock")
		assert.Error(t, err)
	})

	t.Run("invald settings update", func(t *testing.T) {
		_, err := adapter.UpsertSettings(context.Background(), domain.DocSettings{
			CompanyID:  "mock",
			DocAddress: "mock",
		})
		assert.Error(t, err)
	})

	t.Run("update settings with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
		defer cancel()
		_, err := adapter.UpsertSettings(ctx, domain.DocSettings{
			CompanyID:  "mock",
			DocAddress: "mock",
			DocSecret:  "mock",
		})
		assert.Error(t, err)
	})

	t.Run("update settings", func(t *testing.T) {
		_, err := adapter.UpsertSettings(context.Background(), domain.DocSettings{
			CompanyID:  "mock",
			DocAddress: "mock",
			DocSecret:  "mock",
		})
		assert.NoError(t, err)
	})

	adapter.DeleteSettings(context.Background(), "mock")
}
