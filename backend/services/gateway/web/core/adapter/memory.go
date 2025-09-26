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
	"context"
	"encoding/json"
	"errors"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/port"
)

type memoryAICodeAccessAdapter struct {
	kvs map[string][]byte
}

func NewMemoryAICodeAccessAdapter() port.AICodeAccessServiceAdapter {
	return &memoryAICodeAccessAdapter{
		kvs: make(map[string][]byte),
	}
}

func (m *memoryAICodeAccessAdapter) save(code domain.AICodeAccess) error {
	buffer, err := json.Marshal(code)

	if err != nil {
		return err
	}

	m.kvs[code.Code] = buffer

	return nil
}

func (m *memoryAICodeAccessAdapter) SelectCodeAccess(ctx context.Context, code string) (domain.AICodeAccess, error) {
	buffer, ok := m.kvs[code]
	var codeAcc domain.AICodeAccess

	if !ok {
		return codeAcc, errors.New("code with this code doesn't exist")
	}

	if err := json.Unmarshal(buffer, &codeAcc); err != nil {
		return codeAcc, err
	}

	return codeAcc, nil
}

func (m *memoryAICodeAccessAdapter) UpsertCodeAccess(ctx context.Context, code domain.AICodeAccess) error {
	if err := m.save(code); err != nil {
		return err
	}

	return nil
}

func (m *memoryAICodeAccessAdapter) RemoveCodeAccess(ctx context.Context, code string) error {
	if _, ok := m.kvs[code]; !ok {
		return errors.New("code with this code doesn't exist")
	}

	delete(m.kvs, code)

	return nil
}
