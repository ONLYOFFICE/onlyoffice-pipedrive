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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/stretchr/testify/assert"
)

const formPDFHeader = "%PDF-1.7\n%\xCD\xCA\xD2\xA9\x0D1 0 obj\x0A<<\x0A" +
	"/Type /EmbeddedFile /ONLYOFFICEFORM /Length 1234 5678\x0A>>\x0Astream\x0D\x0A"

const plainPDFHeader = "%PDF-1.7\n1 0 obj\n<< /Type /Catalog >>\nendobj\n"

func newTestHandler() ConfigHandler {
	cfg := &shared.OnlyofficeConfig{}
	cfg.Onlyoffice.Builder.PluginGUID = "asc.{plugin-guid}"
	cfg.Onlyoffice.Builder.AutofillerURL = "https://gateway/plugins/aiautofill/"
	cfg.Onlyoffice.Builder.AutofillerConfig = "https://gateway/plugins/aiautofill/config.json"
	cfg.Onlyoffice.Builder.AutofillerCallback = "https://gateway/api/data"
	return ConfigHandler{onlyoffice: cfg, logger: log.NewEmptyLogger()}
}

func pdfServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusPartialContent)
		_, _ = rw.Write([]byte(body))
	}))
}

func TestIsExtendedPDF(t *testing.T) {
	assert.True(t, isExtendedPDF([]byte(formPDFHeader)))
	assert.False(t, isExtendedPDF([]byte(plainPDFHeader)))
	assert.False(t, isExtendedPDF(nil))
}

func TestGetExtension(t *testing.T) {
	c := ConfigHandler{}
	assert.Equal(t, "docx", c.getExtension("Report.DOCX"))
	assert.Equal(t, "gz", c.getExtension("a.tar.gz"))
	assert.Equal(t, "", c.getExtension("README"))
	assert.Equal(t, "", c.getExtension("   "))
}

func TestDemoModeValidity(t *testing.T) {
	c := ConfigHandler{}
	assert.False(t, c.isDemoModeValid(response.DocSettingsResponse{}))
	assert.True(t, c.isDemoModeValid(response.DocSettingsResponse{DemoEnabled: true}))
	assert.True(t, c.isDemoModeValid(response.DocSettingsResponse{DemoEnabled: true, DemoStarted: time.Now().AddDate(0, 0, -29)}))
	assert.False(t, c.isDemoModeValid(response.DocSettingsResponse{DemoEnabled: true, DemoStarted: time.Now().AddDate(0, 0, -31)}))
}

// Smoke: AI Autofill autostarts only for ONLYOFFICE PDF forms and only when
// plugins/autofill are not disabled in company settings (nil means enabled).
func TestGeneratePluginConfig(t *testing.T) {
	form := pdfServer(formPDFHeader)
	defer form.Close()
	plain := pdfServer(plainPDFHeader)
	defer plain.Close()

	c := newTestHandler()
	ctx := context.Background()
	disabled := false

	p := c.generatePluginConfig(ctx, response.DocSettingsResponse{}, "pdf", "code", form.URL)
	assert.Equal(t, []string{"asc.{plugin-guid}"}, p.Autostart)
	assert.Equal(t, "https://gateway/plugins/aiautofill/", p.Url)
	assert.Equal(t, []string{"https://gateway/plugins/aiautofill/config.json"}, p.PluginsData)
	assert.Equal(t, map[string]any{"code": "code", "callback": "https://gateway/api/data"}, p.Options["asc.{plugin-guid}"])

	empty := func(p response.Plugins) {
		assert.Empty(t, p.Autostart)
		assert.Empty(t, p.Url)
		assert.Empty(t, p.Options)
	}

	empty(c.generatePluginConfig(ctx, response.DocSettingsResponse{}, "pdf", "code", plain.URL))
	empty(c.generatePluginConfig(ctx, response.DocSettingsResponse{}, "docx", "code", form.URL))
	empty(c.generatePluginConfig(ctx, response.DocSettingsResponse{AutofillEnabled: &disabled}, "pdf", "code", form.URL))
	empty(c.generatePluginConfig(ctx, response.DocSettingsResponse{PluginsEnabled: &disabled}, "pdf", "code", form.URL))
	empty(c.generatePluginConfig(ctx, response.DocSettingsResponse{}, "pdf", "code", ""))
}
