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
	"encoding/json"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/port"
)

type memoryDocserverAdapter struct {
	kvs map[string][]byte
}

func NewMemoryDocserverAdapter() port.DocSettingsServiceAdapter {
	return &memoryDocserverAdapter{
		kvs: make(map[string][]byte),
	}
}

func (m *memoryDocserverAdapter) save(settings domain.DocSettings) error {
	buffer, err := json.Marshal(settings)

	if err != nil {
		return err
	}

	m.kvs[settings.CompanyID] = buffer

	return nil
}

func (m *memoryDocserverAdapter) InsertSettings(ctx context.Context, settings domain.DocSettings) error {
	return m.save(settings)
}

func (m *memoryDocserverAdapter) SelectSettings(ctx context.Context, cid string) (domain.DocSettings, error) {
	buffer, ok := m.kvs[cid]
	var settings domain.DocSettings

	if !ok {
		return settings, ErrNoCompanySettings
	}

	if err := json.Unmarshal(buffer, &settings); err != nil {
		return settings, err
	}

	return settings, nil
}

func (m *memoryDocserverAdapter) UpsertSettings(ctx context.Context, settings domain.DocSettings) (domain.DocSettings, error) {
	if err := m.save(settings); err != nil {
		return domain.DocSettings{}, err
	}

	return settings, nil
}

func (m *memoryDocserverAdapter) DeleteSettings(ctx context.Context, cid string) error {
	if _, ok := m.kvs[cid]; !ok {
		return ErrNoCompanySettings
	}

	delete(m.kvs, cid)

	return nil
}
