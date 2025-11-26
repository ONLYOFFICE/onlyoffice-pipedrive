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
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	plog "github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/port"
	"github.com/mitchellh/mapstructure"
	"go-micro.dev/v4/cache"
	"golang.org/x/oauth2"
)

type settingsService struct {
	adapter     port.DocSettingsServiceAdapter
	encryptor   crypto.Encryptor
	cache       cache.Cache
	credentials *oauth2.Config
	logger      plog.Logger
}

func NewSettingsService(
	adapter port.DocSettingsServiceAdapter,
	encryptor crypto.Encryptor,
	cache cache.Cache,
	credentials *oauth2.Config,
	logger plog.Logger,
) port.DocSettingsService {
	return settingsService{
		adapter:     adapter,
		encryptor:   encryptor,
		cache:       cache,
		credentials: credentials,
		logger:      logger,
	}
}

func (s settingsService) CreateSettings(ctx context.Context, settings domain.DocSettings) error {
	s.logger.Debugf("validating company %s settings to perform a persist action", settings.CompanyID)
	if err := settings.Validate(); err != nil {
		return err
	}

	esecret, err := s.encryptor.Encrypt(settings.DocSecret, []byte(s.credentials.ClientSecret))
	if err != nil {
		return err
	}

	s.logger.Debugf("settings %s are valid. Persisting to database", settings.CompanyID)
	if err := s.adapter.InsertSettings(ctx, domain.DocSettings{
		CompanyID:       settings.CompanyID,
		DocAddress:      settings.DocAddress,
		DocSecret:       esecret,
		DocHeader:       settings.DocHeader,
		DemoEnabled:     settings.DemoEnabled,
		DemoStarted:     settings.DemoStarted,
		AutofillEnabled: settings.AutofillEnabled,
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

	var settings domain.DocSettings
	var err error
	if res, _, err := s.cache.Get(ctx, id); err == nil && res != nil {
		s.logger.Debugf("found settings %s in the cache", id)
		if err := mapstructure.Decode(res, &settings); err != nil {
			s.logger.Errorf("could not decode from cache: %s", err.Error())
		}
	}

	if settings.CompanyID == "" {
		settings, err = s.adapter.SelectSettings(ctx, id)
		if err != nil {
			return settings, err
		}
		s.cache.Put(ctx, id, settings, 1*time.Minute)
	}

	s.logger.Debugf("found settings: %v", settings)
	dsecret, err := s.encryptor.Decrypt(settings.DocSecret, []byte(s.credentials.ClientSecret))
	if err != nil {
		return settings, err
	}

	return domain.DocSettings{
		CompanyID:       cid,
		DocAddress:      settings.DocAddress,
		DocSecret:       dsecret,
		DocHeader:       settings.DocHeader,
		DemoEnabled:     settings.DemoEnabled,
		DemoStarted:     settings.DemoStarted,
		AutofillEnabled: settings.AutofillEnabled,
	}, nil
}

func (s settingsService) UpdateSettings(ctx context.Context, settings domain.DocSettings) (domain.DocSettings, error) {
	s.logger.Debugf("validating settings %s to perform an update action", settings.CompanyID)
	if err := settings.Validate(); err != nil {
		return settings, err
	}

	persistedSettings, err := s.adapter.SelectSettings(ctx, settings.CompanyID)
	if err != nil {
		if settings.DemoEnabled {
			settings.DemoStarted = time.Now()
		}
	} else {
		if settings.DemoEnabled {
			if persistedSettings.DemoEnabled && !persistedSettings.DemoStarted.IsZero() {
				settings.DemoStarted = persistedSettings.DemoStarted
			} else {
				settings.DemoStarted = time.Now()
			}
		}
	}

	esecret, err := s.encryptor.Encrypt(settings.DocSecret, []byte(s.credentials.ClientSecret))
	if err != nil {
		return settings, err
	}

	s.logger.Debugf("settings %s are valid to perform an update action", settings.CompanyID)
	if _, err := s.adapter.UpsertSettings(ctx, domain.DocSettings{
		CompanyID:       settings.CompanyID,
		DocAddress:      settings.DocAddress,
		DocSecret:       esecret,
		DocHeader:       settings.DocHeader,
		DemoEnabled:     settings.DemoEnabled,
		DemoStarted:     settings.DemoStarted,
		AutofillEnabled: settings.AutofillEnabled,
	}); err != nil {
		return settings, err
	}

	s.cache.Delete(ctx, settings.CompanyID)

	s.logger.Debugf("successfully persisted %s settings", settings.CompanyID)
	return settings, nil
}

func (s settingsService) RemoveSettings(ctx context.Context, cid string) error {
	id := strings.TrimSpace(cid)
	s.logger.Debugf("validating cid %s to perform a delete action", id)

	if id == "" {
		return &InvalidServiceParameterError{
			Name:   "CID",
			Reason: "Should not be blank",
		}
	}

	s.cache.Delete(ctx, cid)

	s.logger.Debugf("uid %s is valid to perform a delete action", id)
	return s.adapter.DeleteSettings(ctx, id)
}
