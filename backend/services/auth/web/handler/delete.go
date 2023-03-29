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

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/core/port"
	"go-micro.dev/v4/client"
)

type UserDeleteHandler struct {
	service port.UserAccessService
	client  client.Client
	logger  log.Logger
}

func NewUserDeleteHandler(
	service port.UserAccessService,
	client client.Client,
	logger log.Logger,
) UserDeleteHandler {
	return UserDeleteHandler{
		service: service,
		client:  client,
		logger:  logger,
	}
}

func (u UserDeleteHandler) DeleteUser(ctx context.Context, uid *string, res *interface{}) error {
	u.logger.Debugf("removing user %s", *uid)
	return u.service.DeleteUser(ctx, *uid)
}
