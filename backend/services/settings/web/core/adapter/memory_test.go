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

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestMemoryAdapter(t *testing.T) {
	adapter := NewMemoryDocserverAdapter()

	t.Run("save settings", func(t *testing.T) {
		assert.NoError(t, adapter.InsertSettings(context.Background(), settings))
	})

	t.Run("save the same settings object", func(t *testing.T) {
		assert.NoError(t, adapter.InsertSettings(context.Background(), settings))
	})

	t.Run("get settings by cid", func(t *testing.T) {
		s, err := adapter.SelectSettings(context.Background(), "mock")
		assert.NoError(t, err)
		assert.Equal(t, settings, s)
	})

	t.Run("update settings by cid", func(t *testing.T) {
		s, err := adapter.UpsertSettings(context.Background(), domain.DocSettings{
			CompanyID:  "mock",
			DocAddress: "mock",
			DocSecret:  "mock",
		})
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})

	t.Run("delete settings by cid", func(t *testing.T) {
		assert.NoError(t, adapter.DeleteSettings(context.Background(), "mock"))
	})

	t.Run("get invalid settings", func(t *testing.T) {
		_, err := adapter.SelectSettings(context.Background(), "mock")
		assert.Error(t, err)
	})
}
