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

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/port"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type aiCodeAccessCollection struct {
	mgm.DefaultModel `bson:",inline"`
	Code             string `json:"code" bson:"code"`
	UserID           string `json:"user_id" bson:"user_id"`
	DealID           string `json:"deal_id" bson:"deal_id"`
}

type mongoAICodeAccessAdapter struct {
}

func NewMongoAICodeAccessAdapter(url string) port.AICodeAccessServiceAdapter {
	if err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: 3 * time.Second}, "pipedrive",
		options.Client().ApplyURI(url),
	); err != nil {
		log.Fatalf("mongo initialization error: %s", err.Error())
	}

	return &mongoAICodeAccessAdapter{}
}

func (m *mongoAICodeAccessAdapter) save(ctx context.Context, code domain.AICodeAccess) error {
	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		c := &aiCodeAccessCollection{}
		collection := mgm.Coll(&aiCodeAccessCollection{})

		if err := collection.FirstWithCtx(ctx, bson.M{"code": code.Code}, c); err != nil {
			if cerr := collection.CreateWithCtx(ctx, &aiCodeAccessCollection{
				Code:   code.Code,
				UserID: code.UserID,
				DealID: code.DealID,
			}); cerr != nil {
				return cerr
			}

			return session.CommitTransaction(sc)
		}

		c.UserID = code.UserID
		c.DealID = code.DealID
		c.UpdatedAt = time.Now()

		if err := collection.UpdateWithCtx(ctx, c); err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

func (m *mongoAICodeAccessAdapter) UpsertCodeAccess(ctx context.Context, code domain.AICodeAccess) error {
	if err := code.Validate(); err != nil {
		return err
	}

	return m.save(ctx, code)
}

func (m *mongoAICodeAccessAdapter) SelectCodeAccess(ctx context.Context, code string) (domain.AICodeAccess, error) {
	code = strings.TrimSpace(code)

	if code == "" {
		return domain.AICodeAccess{}, ErrInvalidCode
	}

	c := &aiCodeAccessCollection{}
	collection := mgm.Coll(&aiCodeAccessCollection{})
	return domain.AICodeAccess{
		Code:   c.Code,
		UserID: c.UserID,
		DealID: c.DealID,
	}, collection.FirstWithCtx(ctx, bson.M{"code": code}, c)
}

func (m *mongoAICodeAccessAdapter) RemoveCodeAccess(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)

	if code == "" {
		return ErrInvalidCode
	}

	_, err := mgm.Coll(&aiCodeAccessCollection{}).DeleteMany(ctx, bson.M{"code": code})
	return err
}
