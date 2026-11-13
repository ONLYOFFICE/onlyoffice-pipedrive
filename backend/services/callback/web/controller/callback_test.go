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

package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/stretchr/testify/assert"
)

func TestCreateFileVersion(t *testing.T) {
	c := CallbackController{}
	cases := map[string]string{
		"report.docx":            "report (1).docx",
		"report (1).docx":        "report (2).docx",
		"report (9).docx":        "report (10).docx",
		"my (draft).xlsx":        "my (draft) (1).xlsx",
		"archive.tar.gz":         "archive.tar (1).gz",
		"README":                 "README (1)",
		"Q1 plan (12).pptx":      "Q1 plan (13).pptx",
		"no-space(3).docx":       "no-space(3) (1).docx",
		"unicode résumé.docx":    "unicode résumé (1).docx",
		"unicode résumé (4).pdf": "unicode résumé (5).pdf",
	}

	for in, want := range cases {
		assert.Equal(t, want, c.createFileVersion(in), in)
	}
}

func TestDemoModeValidity(t *testing.T) {
	c := CallbackController{}
	assert.False(t, c.isDemoModeValid(response.DocSettingsResponse{DemoEnabled: false}))
	assert.True(t, c.isDemoModeValid(response.DocSettingsResponse{DemoEnabled: true}))
	assert.True(t, c.isDemoModeValid(response.DocSettingsResponse{DemoEnabled: true, DemoStarted: time.Now().AddDate(0, 0, -29)}))
	assert.False(t, c.isDemoModeValid(response.DocSettingsResponse{DemoEnabled: true, DemoStarted: time.Now().AddDate(0, 0, -31)}))
}

func TestCallbackRejectsInvalidRequests(t *testing.T) {
	c := NewCallbackController(nil, pclient.PipedriveApiClient{}, crypto.NewJwtManager(&config.CryptoConfig{}),
		&config.ServerConfig{}, &shared.OnlyofficeConfig{}, log.NewEmptyLogger())
	handler := c.BuildPostHandleCallback()

	cases := []struct {
		name  string
		query string
		body  string
	}{
		{"invalid json", "?cid=1&did=2&fid=3", "{"},
		{"missing token", "?cid=1&did=2&fid=3", `{"status":2}`},
		{"missing query params", "?cid=1", `{"status":2,"token":"x"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/callback"+tc.query, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.JSONEq(t, `{"error":1}`, rec.Body.String())
		})
	}
}
