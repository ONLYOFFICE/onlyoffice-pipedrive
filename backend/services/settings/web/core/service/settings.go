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
	"errors"
	"strings"

	plog "github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/port"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/crypto"
)

var _ErrOperationTimeout = errors.New("operation timeout")

type settingsService struct {
	adapter   port.DocSettingsServiceAdapter
	encryptor crypto.Encryptor
	logger    plog.Logger
}

func NewSettingsService(
	adapter port.DocSettingsServiceAdapter,
	encryptor crypto.Encryptor,
	logger plog.Logger,
) port.DocSettingsService {
	return settingsService{
		adapter:   adapter,
		encryptor: encryptor,
		logger:    logger,
	}
}

func (s settingsService) CreateSettings(ctx context.Context, settings domain.DocSettings) error {
	s.logger.Debugf("validating company %s settings to perform a persist action", settings.CompanyID)
	if err := settings.Validate(); err != nil {
		return err
	}

	esecret, err := s.encryptor.Encrypt(settings.DocSecret)
	if err != nil {
		return err
	}

	s.logger.Debugf("settings %s are valid. Persisting to database", settings.CompanyID)
	if err := s.adapter.InsertSettings(ctx, domain.DocSettings{
		CompanyID:  settings.CompanyID,
		DocAddress: settings.DocAddress,
		DocSecret:  esecret,
		DocHeader:  settings.DocHeader,
	}); err != nil {
		return err
	}

	return nil
}

func (s settingsService) GetSettings(ctx context.Context, cid string) (domain.DocSettings, error) {
	s.logger.Debugf("trying to select settings for company with id: %s", cid)
	id := strings.TrimSpace(cid)

	if id == "" {
		return domain.DocSettings{}, &InvalidServiceParameterError{
			Name:   "CID",
			Reason: "Should not be blank",
		}
	}

	settings, err := s.adapter.SelectSettings(ctx, id)
	if err != nil {
		return settings, err
	}

	s.logger.Debugf("found settings: %v", settings)
	dsecret, err := s.encryptor.Decrypt(settings.DocSecret)
	if err != nil {
		return settings, err
	}

	return domain.DocSettings{
		CompanyID:  cid,
		DocAddress: settings.DocAddress,
		DocSecret:  dsecret,
		DocHeader:  settings.DocHeader,
	}, nil
}

func (s settingsService) UpdateSettings(ctx context.Context, settings domain.DocSettings) (domain.DocSettings, error) {
	s.logger.Debugf("validating settings %s to perform an update action", settings.CompanyID)
	if err := settings.Validate(); err != nil {
		return settings, err
	}

	esecret, err := s.encryptor.Encrypt(settings.DocSecret)
	if err != nil {
		return settings, err
	}

	s.logger.Debugf("settings %s are valid to perform an update action", settings.CompanyID)
	if _, err := s.adapter.UpsertSettings(ctx, domain.DocSettings{
		CompanyID:  settings.CompanyID,
		DocAddress: settings.DocAddress,
		DocSecret:  esecret,
		DocHeader:  settings.DocHeader,
	}); err != nil {
		return settings, err
	}

	s.logger.Debugf("successfully persisted %s settings", settings.CompanyID)
	return settings, nil
}

func (s settingsService) DeleteSettings(ctx context.Context, cid string) error {
	id := strings.TrimSpace(cid)
	s.logger.Debugf("validating cid %s to perform a delete action", id)

	if id == "" {
		return &InvalidServiceParameterError{
			Name:   "CID",
			Reason: "Should not be blank",
		}
	}

	s.logger.Debugf("uid %s is valid to perform a delete action", id)
	return s.adapter.DeleteSettings(ctx, id)
}
