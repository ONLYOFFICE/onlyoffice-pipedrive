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

package service

import (
	"context"
	"strings"

	plog "github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/port"
)

type aiAccessService struct {
	adapter port.AICodeAccessServiceAdapter
	logger  plog.Logger
}

func NewAIAccessService(
	adapter port.AICodeAccessServiceAdapter,
	logger plog.Logger,
) port.AICodeAccessService {
	return aiAccessService{
		adapter: adapter,
		logger:  logger,
	}
}

func (s aiAccessService) UpsertCodeAccess(ctx context.Context, codeAccess domain.AICodeAccess) error {
	s.logger.Debugf("validating AI code access %s to perform an upsert action", codeAccess.Code)
	if err := codeAccess.Validate(); err != nil {
		return err
	}

	s.logger.Debugf("AI code access %s is valid. Persisting to database", codeAccess.Code)
	if err := s.adapter.UpsertCodeAccess(ctx, codeAccess); err != nil {
		return err
	}

	return nil
}

func (s aiAccessService) GetCodeAccess(ctx context.Context, code string) (domain.AICodeAccess, error) {
	s.logger.Debugf("trying to select AI code access with code: %s", code)
	code = strings.TrimSpace(code)

	if code == "" {
		return domain.AICodeAccess{}, &InvalidServiceParameterError{
			Name:   "Code",
			Reason: "Should not be blank",
		}
	}

	codeAccess, err := s.adapter.SelectCodeAccess(ctx, code)
	if err != nil {
		return codeAccess, err
	}

	s.logger.Debugf("found AI code access: %v", codeAccess)

	return codeAccess, nil
}

func (s aiAccessService) RemoveCodeAccess(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)
	s.logger.Debugf("validating code %s to perform a delete action", code)

	if code == "" {
		return &InvalidServiceParameterError{
			Name:   "Code",
			Reason: "Should not be blank",
		}
	}

	s.logger.Debugf("code %s is valid to perform a delete action", code)
	return s.adapter.RemoveCodeAccess(ctx, code)
}
