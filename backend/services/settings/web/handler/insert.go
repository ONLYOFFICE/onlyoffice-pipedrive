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

package handler

import (
	"context"
	"fmt"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/port"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
)

type SettingsInsertHandler struct {
	service port.DocSettingsService
	logger  log.Logger
}

func NewSettingsInsertHandler(
	service port.DocSettingsService,
	logger log.Logger,
) SettingsInsertHandler {
	return SettingsInsertHandler{
		service: service,
		logger:  logger,
	}
}

func (i SettingsInsertHandler) InsertSettings(ctx context.Context, req request.DocSettings, res *interface{}) error {
	_, err, _ := group.Do(fmt.Sprintf("insert-%d", req.CompanyID), func() (interface{}, error) {
		settings, err := i.service.UpdateSettings(ctx, domain.DocSettings{
			CompanyID:       fmt.Sprint(req.CompanyID),
			DocAddress:      req.DocAddress,
			DocHeader:       req.DocHeader,
			DocSecret:       req.DocSecret,
			DemoEnabled:     req.DemoEnabled,
			AutofillEnabled: req.AutofillEnabled,
		})

		if err != nil {
			i.logger.Errorf("could not update settings: %s", err.Error())
			return nil, err
		}

		return settings, nil
	})

	return err
}
