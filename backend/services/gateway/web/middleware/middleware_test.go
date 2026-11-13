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

package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var testCredentials = &oauth2.Config{ClientID: "client-id", ClientSecret: "client-secret"}

func signContext(t *testing.T, secret string, uid, cid int) string {
	t.Helper()
	token, err := crypto.NewJwtManager(&config.CryptoConfig{}).Sign(secret, request.PipedriveTokenContext{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
		UID:              uid,
		CID:              cid,
	})
	assert.NoError(t, err)
	return token
}

func TestContextMiddleware(t *testing.T) {
	m := BuildHandleContextMiddleware(crypto.NewJwtManager(&config.CryptoConfig{}), testCredentials, log.NewEmptyLogger())

	var got request.PipedriveTokenContext
	handler := m.Protect(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		got, _ = r.Context().Value("X-Pipedrive-App-Context").(request.PipedriveTokenContext)
		rw.WriteHeader(http.StatusOK)
	}))

	cases := []struct {
		name   string
		header string
		status int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"wrong secret", signContext(t, "another-secret", 1, 2), http.StatusForbidden},
		{"garbage token", "not-a-jwt", http.StatusForbidden},
		{"valid token", signContext(t, testCredentials.ClientSecret, 11, 22), http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
			if tc.header != "" {
				req.Header.Set("X-Pipedrive-App-Context", tc.header)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tc.status, rec.Code)
		})
	}

	assert.Equal(t, 11, got.UID)
	assert.Equal(t, 22, got.CID)
}

func TestAuthMiddleware(t *testing.T) {
	m := BuildHandleAuthMiddleware(testCredentials, log.NewEmptyLogger())
	handler := m.Protect(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	valid := base64.StdEncoding.EncodeToString([]byte("client-id:client-secret"))
	invalid := base64.StdEncoding.EncodeToString([]byte("client-id:wrong"))

	cases := []struct {
		name   string
		header string
		status int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"wrong credentials", "Basic " + invalid, http.StatusForbidden},
		{"valid credentials", "Basic " + valid, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/oauth/auth", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tc.status, rec.Code)
		})
	}
}
