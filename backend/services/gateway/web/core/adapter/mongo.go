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
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/port"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

type aiCodeAccessCollection struct {
	mgm.DefaultModel `bson:",inline"`
	Code             string `json:"code" bson:"code"`
	CodeHash         string `json:"code_hash" bson:"code_hash"`
	UserID           string `json:"user_id" bson:"user_id"`
	FileID           string `json:"file_id" bson:"file_id"`
	DealID           string `json:"deal_id" bson:"deal_id"`
}

type mongoAICodeAccessAdapter struct {
	encryptor   crypto.Encryptor
	credentials *oauth2.Config
}

func hashCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

func NewMongoAICodeAccessAdapter(url string, encryptor crypto.Encryptor, credentials *oauth2.Config) port.AICodeAccessServiceAdapter {
	if err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: 3 * time.Second}, "pipedrive",
		options.Client().ApplyURI(url),
	); err != nil {
		log.Fatalf("mongo initialization error: %s", err.Error())
	}

	return &mongoAICodeAccessAdapter{
		encryptor:   encryptor,
		credentials: credentials,
	}
}

func (m *mongoAICodeAccessAdapter) save(ctx context.Context, code domain.AICodeAccess) error {
	encryptedCode, err := m.encryptor.Encrypt(code.Code, []byte(m.credentials.ClientSecret))
	if err != nil {
		return err
	}

	codeHash := hashCode(code.Code)
	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		c := &aiCodeAccessCollection{}
		collection := mgm.Coll(&aiCodeAccessCollection{})

		if err := collection.FirstWithCtx(ctx, bson.M{"user_id": code.UserID, "file_id": code.FileID}, c); err != nil {
			if cerr := collection.CreateWithCtx(ctx, &aiCodeAccessCollection{
				Code:     encryptedCode,
				CodeHash: codeHash,
				UserID:   code.UserID,
				FileID:   code.FileID,
				DealID:   code.DealID,
			}); cerr != nil {
				return cerr
			}

			return session.CommitTransaction(sc)
		}

		c.Code = encryptedCode
		c.CodeHash = codeHash
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
	if err := collection.FirstWithCtx(ctx, bson.M{"code_hash": hashCode(code)}, c); err != nil {
		return domain.AICodeAccess{}, err
	}

	decryptedCode, err := m.encryptor.Decrypt(c.Code, []byte(m.credentials.ClientSecret))
	if err != nil {
		return domain.AICodeAccess{}, err
	}

	return domain.AICodeAccess{
		Code:   decryptedCode,
		UserID: c.UserID,
		FileID: c.FileID,
		DealID: c.DealID,
	}, nil
}

func (m *mongoAICodeAccessAdapter) SelectCodeAccessByUserAndFile(ctx context.Context, userID, fileID string) (domain.AICodeAccess, error) {
	userID = strings.TrimSpace(userID)
	fileID = strings.TrimSpace(fileID)

	if userID == "" {
		return domain.AICodeAccess{}, ErrInvalidCode
	}

	if fileID == "" {
		return domain.AICodeAccess{}, ErrInvalidCode
	}

	c := &aiCodeAccessCollection{}
	collection := mgm.Coll(&aiCodeAccessCollection{})
	if err := collection.FirstWithCtx(ctx, bson.M{"user_id": userID, "file_id": fileID}, c); err != nil {
		return domain.AICodeAccess{}, err
	}

	decryptedCode, err := m.encryptor.Decrypt(c.Code, []byte(m.credentials.ClientSecret))
	if err != nil {
		return domain.AICodeAccess{}, err
	}

	return domain.AICodeAccess{
		Code:   decryptedCode,
		UserID: c.UserID,
		FileID: c.FileID,
		DealID: c.DealID,
	}, nil
}

func (m *mongoAICodeAccessAdapter) RemoveCodeAccess(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)

	if code == "" {
		return ErrInvalidCode
	}

	_, err := mgm.Coll(&aiCodeAccessCollection{}).DeleteMany(ctx, bson.M{"code_hash": hashCode(code)})
	return err
}
