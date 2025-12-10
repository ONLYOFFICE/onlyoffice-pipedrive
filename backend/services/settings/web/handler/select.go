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

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/port"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"go-micro.dev/v4/client"
)

type SettingsSelectHandler struct {
	service port.DocSettingsService
	client  client.Client
	logger  log.Logger
}

func NewSettingsSelectHandler(
	service port.DocSettingsService,
	client client.Client,
	logger log.Logger,
) SettingsSelectHandler {
	return SettingsSelectHandler{
		service: service,
		client:  client,
		logger:  logger,
	}
}

func (u SettingsSelectHandler) GetSettings(ctx context.Context, cid *string, res *response.DocSettingsResponse) error {
	settings, err, _ := group.Do(*cid, func() (interface{}, error) {
		settings, err := u.service.GetSettings(ctx, *cid)
		if err != nil {
			u.logger.Warnf("could not get company %s settings. Reason: %s", *cid, err.Error())
			return settings, nil
		}

		return settings, nil
	})

	if set, ok := settings.(domain.DocSettings); ok {
		pluginsEnabled := set.GetPluginsEnabled()
		autofillEnabled := set.GetAutofillEnabled()
		*res = response.DocSettingsResponse{
			DocAddress:      set.DocAddress,
			DocSecret:       set.DocSecret,
			DocHeader:       set.DocHeader,
			DemoEnabled:     set.DemoEnabled,
			DemoStarted:     set.DemoStarted,
			PluginsEnabled:  &pluginsEnabled,
			AutofillEnabled: &autofillEnabled,
		}
		return nil
	}

	return err
}
