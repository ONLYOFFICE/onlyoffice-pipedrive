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
	"log"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web/core/port"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type docSettingsCollection struct {
	mgm.DefaultModel `bson:",inline"`
	CompanyID        string    `json:"company_id" bson:"company_id"`
	DocAddress       string    `json:"doc_address" bson:"doc_address"`
	DocSecret        string    `json:"doc_secret" bson:"doc_secret"`
	DocHeader        string    `json:"doc_header" bson:"doc_header"`
	DemoEnabled      bool      `json:"demo_enabled" bson:"demo_enabled"`
	DemoStarted      time.Time `json:"demo_started" bson:"demo_started"`
}

type mongoUserAdapter struct {
}

func NewMongoDocserverAdapter(url string) port.DocSettingsServiceAdapter {
	if err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: 3 * time.Second}, "pipedrive",
		options.Client().ApplyURI(url),
	); err != nil {
		log.Fatalf("mongo initialization error: %s", err.Error())
	}

	return &mongoUserAdapter{}
}

func (m *mongoUserAdapter) save(ctx context.Context, settings domain.DocSettings) error {
	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		u := &docSettingsCollection{}
		collection := mgm.Coll(&docSettingsCollection{})

		if err := collection.FirstWithCtx(ctx, bson.M{"company_id": settings.CompanyID}, u); err != nil {
			if cerr := collection.CreateWithCtx(ctx, &docSettingsCollection{
				CompanyID:   settings.CompanyID,
				DocAddress:  settings.DocAddress,
				DocSecret:   settings.DocSecret,
				DocHeader:   settings.DocHeader,
				DemoEnabled: settings.DemoEnabled,
				DemoStarted: settings.DemoStarted,
			}); cerr != nil {
				return cerr
			}

			return session.CommitTransaction(sc)
		}

		u.CompanyID = settings.CompanyID
		u.DocAddress = settings.DocAddress
		u.DocSecret = settings.DocSecret
		u.DocHeader = settings.DocHeader
		u.DemoEnabled = settings.DemoEnabled
		if u.DemoStarted.IsZero() {
			u.DemoStarted = settings.DemoStarted
		}

		u.UpdatedAt = time.Now()

		if err := collection.UpdateWithCtx(ctx, u); err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (m *mongoUserAdapter) InsertSettings(ctx context.Context, settings domain.DocSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}

	return m.save(ctx, settings)
}

func (m *mongoUserAdapter) SelectSettings(ctx context.Context, cid string) (domain.DocSettings, error) {
	cid = strings.TrimSpace(cid)

	if cid == "" {
		return domain.DocSettings{}, ErrInvalidCompanyID
	}

	settings := &docSettingsCollection{}
	collection := mgm.Coll(settings)

	if err := collection.FirstWithCtx(ctx, bson.M{"company_id": cid}, settings); err != nil {
		return domain.DocSettings{}, err
	}

	return domain.DocSettings{
		CompanyID:   settings.CompanyID,
		DocAddress:  settings.DocAddress,
		DocSecret:   settings.DocSecret,
		DocHeader:   settings.DocHeader,
		DemoEnabled: settings.DemoEnabled,
		DemoStarted: settings.DemoStarted,
	}, nil
}

func (m *mongoUserAdapter) UpsertSettings(ctx context.Context, settings domain.DocSettings) (domain.DocSettings, error) {
	if err := settings.Validate(); err != nil {
		return settings, err
	}

	return settings, m.save(ctx, settings)
}

func (m *mongoUserAdapter) DeleteSettings(ctx context.Context, cid string) error {
	cid = strings.TrimSpace(cid)

	if cid == "" {
		return ErrInvalidCompanyID
	}

	_, err := mgm.Coll(&docSettingsCollection{}).DeleteMany(ctx, bson.M{"company_id": bson.M{operator.Eq: cid}})
	return err
}
