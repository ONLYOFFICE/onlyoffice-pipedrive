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
package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	plog "github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	shared "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mileusna/useragent"
	"go-micro.dev/v4/client"
	"golang.org/x/sync/errgroup"
)

var (
	defaultClient = &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

type ConfigHandler struct {
	client        client.Client
	apiClient     pclient.PipedriveApiClient
	jwtManager    crypto.JwtManager
	config        *config.ServerConfig
	onlyoffice    *shared.OnlyofficeConfig
	logger        plog.Logger
	formatManager shared.FormatManager
}

func NewConfigHandler(
	client client.Client,
	jwtManager crypto.JwtManager,
	apiClient pclient.PipedriveApiClient,
	config *config.ServerConfig,
	onlyoffice *shared.OnlyofficeConfig,
	formatManager shared.FormatManager,
	logger plog.Logger,
) ConfigHandler {
	return ConfigHandler{
		client:        client,
		apiClient:     apiClient,
		jwtManager:    jwtManager,
		config:        config,
		onlyoffice:    onlyoffice,
		logger:        logger,
		formatManager: formatManager,
	}
}

func (c ConfigHandler) isDemoModeValid(settings response.DocSettingsResponse) bool {
	if !settings.DemoEnabled {
		return false
	}

	if settings.DemoStarted.IsZero() {
		return true
	}

	staleDate := time.Now().AddDate(0, 0, -30)
	return settings.DemoStarted.After(staleDate)
}

func (c ConfigHandler) processConfig(user response.UserResponse, req request.BuildConfigRequest, ctx context.Context) (response.BuildConfigResponse, error) {
	var config response.BuildConfigResponse

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	g, gctx := errgroup.WithContext(tctx)

	var usr model.User
	var settings response.DocSettingsResponse

	g.Go(func() error {
		u, err := c.apiClient.GetMe(gctx, model.Token{
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
			TokenType:    user.TokenType,
			Scope:        user.Scope,
			ApiDomain:    user.ApiDomain,
		})
		if err != nil {
			c.logger.Debugf("could not get pipedrive user: %s", err.Error())
			return err
		}
		usr = u
		return nil
	})

	g.Go(func() error {
		var docs response.DocSettingsResponse
		if err := c.client.Call(
			gctx,
			c.client.NewRequest(
				fmt.Sprintf("%s:settings", c.config.Namespace),
				"SettingsSelectHandler.GetSettings",
				fmt.Sprint(req.CID),
			),
			&docs,
		); err != nil {
			c.logger.Debugf("could not document server settings: %s", err.Error())
			return err
		}

		if c.isDemoModeValid(docs) {
			if c.onlyoffice.Onlyoffice.Demo.DocumentServerURL == "" ||
				c.onlyoffice.Onlyoffice.Demo.DocumentServerSecret == "" ||
				c.onlyoffice.Onlyoffice.Demo.DocumentServerHeader == "" {
				c.logger.Errorf("demo mode is enabled but demo credentials are not configured")
				return ErrNoSettingsFound
			}

			c.logger.Debugf("using demo mode for company %d", req.CID)
			docs.DocAddress = c.onlyoffice.Onlyoffice.Demo.DocumentServerURL
			docs.DocSecret = c.onlyoffice.Onlyoffice.Demo.DocumentServerSecret
			docs.DocHeader = c.onlyoffice.Onlyoffice.Demo.DocumentServerHeader
		} else {
			if docs.DocAddress == "" || docs.DocSecret == "" || docs.DocHeader == "" {
				c.logger.Debugf("no settings found and demo mode not valid")
				return ErrNoSettingsFound
			}
			c.logger.Debugf("using regular document server settings for company %d", req.CID)
		}

		settings = docs
		return nil
	})

	if err := g.Wait(); err != nil {
		return config, err
	}

	t := "desktop"
	ua := useragent.Parse(req.UserAgent)
	if ua.Mobile || ua.Tablet {
		t = "mobile"
	}

	dreq, err := http.NewRequestWithContext(tctx, "GET", fmt.Sprintf("%s/files/%s/download", user.ApiDomain, req.FileID), nil)
	if err != nil {
		return config, fmt.Errorf("failed to create request: %w", err)
	}
	dreq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))
	resp, err := defaultClient.Do(dreq)
	if err != nil {
		return config, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	filename := c.formatManager.EscapeFileName(req.Filename)
	theme := "default-light"
	if req.Dark {
		theme = "default-dark"
	}

	plugins := response.Plugins{
		Options: map[string]any{},
	}

	autofillEnabled := true
	if settings.AutofillEnabled != nil {
		autofillEnabled = *settings.AutofillEnabled
	}

	if autofillEnabled {
		plugins.Autostart = []string{c.onlyoffice.Onlyoffice.Builder.PluginGUID}
		plugins.Options[c.onlyoffice.Onlyoffice.Builder.PluginGUID] = map[string]any{
			"code":     req.Code,
			"callback": c.onlyoffice.Onlyoffice.Builder.AutofillerCallback,
		}
		plugins.PluginsData = []string{c.onlyoffice.Onlyoffice.Builder.AutofillerConfig}
		plugins.Url = c.onlyoffice.Onlyoffice.Builder.AutofillerURL
	}

	config = response.BuildConfigResponse{
		Document: response.Document{
			Key:   req.DocKey,
			Title: filename,
			URL:   resp.Header.Get("Location"),
		},
		EditorConfig: response.EditorConfig{
			User: response.User{
				ID:   fmt.Sprint(usr.ID + usr.CompanyID),
				Name: usr.Name,
			},
			CallbackURL: fmt.Sprintf(
				"%s/callback?cid=%d&did=%s&fid=%s&filename=%s",
				c.onlyoffice.Onlyoffice.Builder.CallbackURL,
				usr.CompanyID, req.Deal, req.FileID,
				url.QueryEscape(filename),
			),
			Customization: response.Customization{
				Goback: response.Goback{
					RequestClose: false,
				},
				Plugins:       true,
				HideRightMenu: false,
				UiTheme:       theme,
			},
			Lang:    usr.Language.Lang,
			Plugins: plugins,
		},
		Type:        t,
		ServerURL:   settings.DocAddress,
		DemoEnabled: settings.DemoEnabled,
	}

	var fileType string
	var isEditable bool

	if strings.TrimSpace(filename) != "" {
		ext := strings.ReplaceAll(filepath.Ext(filename), ".", "")
		config.Document.FileType = strings.ToLower(ext)
		format, exists := c.formatManager.GetFormatByName(ext)
		if !exists {
			return config, fmt.Errorf("format not supported: %s", ext)
		}

		fileType = format.Type
		isEditable = format.IsEditable()

		config.Document.Permissions = response.Permissions{
			Edit:                 isEditable,
			Comment:              true,
			Download:             true,
			Print:                false,
			Review:               false,
			Copy:                 true,
			ModifyContentControl: true,
			ModifyFilter:         true,
		}
		config.DocumentType = fileType
	}

	config.ExpiresAt = jwt.NewNumericDate(time.Now().Add(5 * time.Minute))
	token, err := c.jwtManager.Sign(settings.DocSecret, config)
	if err != nil {
		c.logger.Debugf("could not sign document server config: %s", err.Error())
		return config, err
	}

	config.Token = token
	return config, nil
}

func (c ConfigHandler) BuildConfig(ctx context.Context, payload request.BuildConfigRequest, res *response.BuildConfigResponse) error {
	c.logger.Debugf("processing a docs config: %s", payload.Filename)

	req := c.client.NewRequest(
		fmt.Sprintf("%s:auth", c.config.Namespace), "UserSelectHandler.GetUser",
		fmt.Sprint(payload.UID+payload.CID),
	)

	var ures response.UserResponse
	if err := c.client.Call(ctx, req, &ures); err != nil {
		c.logger.Debugf("could not get user %d access info: %s", payload.UID+payload.CID, err.Error())
		return err
	}

	config, err := c.processConfig(ures, payload, ctx)
	if err != nil {
		return err
	}

	*res = config
	return nil
}
