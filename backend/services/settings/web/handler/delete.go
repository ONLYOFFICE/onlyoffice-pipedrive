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
	"fmt"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/port"
	"go-micro.dev/v4/client"
)

type SettingsDeleteHandler struct {
	service port.DocSettingsService
	client  client.Client
	logger  log.Logger
}

func NewSettingsDeleteHandler(
	service port.DocSettingsService,
	client client.Client,
	logger log.Logger,
) SettingsDeleteHandler {
	return SettingsDeleteHandler{
		service: service,
		client:  client,
		logger:  logger,
	}
}

func (u SettingsDeleteHandler) DeleteSettings(ctx context.Context, cid *string, res *interface{}) error {
	_, err, _ := group.Do(fmt.Sprintf("remove-%s", *cid), func() (interface{}, error) {
		u.logger.Debugf("removing settings %s", *cid)
		if err := u.service.RemoveSettings(ctx, *cid); err != nil {
			u.logger.Debugf("could not delete settings %s: %s", *cid, err.Error())
			return nil, err
		}

		return nil, nil
	})

	return err
}
