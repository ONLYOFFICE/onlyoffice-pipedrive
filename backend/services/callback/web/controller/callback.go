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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	plog "github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/util/backoff"
)

type CallbackController struct {
	client       client.Client
	pipedriveAPI pclient.PipedriveApiClient
	jwtManager   crypto.JwtManager
	config       *config.ServerConfig
	onlyoffice   *shared.OnlyofficeConfig
	logger       plog.Logger
}

func NewCallbackController(
	client client.Client,
	pipedriveAPI pclient.PipedriveApiClient,
	jwtManager crypto.JwtManager,
	config *config.ServerConfig,
	onlyoffice *shared.OnlyofficeConfig,
	logger plog.Logger,
) *CallbackController {
	return &CallbackController{
		client:       client,
		pipedriveAPI: pipedriveAPI,
		jwtManager:   jwtManager,
		config:       config,
		onlyoffice:   onlyoffice,
		logger:       logger,
	}
}

func (c CallbackController) isDemoModeValid(settings response.DocSettingsResponse) bool {
	if !settings.DemoEnabled {
		return false
	}

	if settings.DemoStarted.IsZero() {
		return true
	}

	staleDate := time.Now().AddDate(0, 0, -30)
	return settings.DemoStarted.After(staleDate)
}

func (c CallbackController) BuildPostHandleCallback() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		cid, did, fid := strings.TrimSpace(query.Get("cid")), strings.TrimSpace(query.Get("did")), strings.TrimSpace(query.Get("fid"))
		rw.Header().Set("Content-Type", "application/json")

		var body request.CallbackRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			c.logger.Errorf("could not decode a callback body")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if body.Token == "" {
			c.logger.Error("invalid callback body token")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if cid == "" || did == "" || fid == "" {
			c.logger.Error("invalid query parameter")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		req := c.client.NewRequest(fmt.Sprintf("%s:settings", c.config.Namespace), "SettingsSelectHandler.GetSettings", cid)
		var res response.DocSettingsResponse
		if err := c.client.Call(r.Context(), req, &res); err != nil {
			c.logger.Errorf("could not extract doc server settings %s", cid)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		var jwtSecret string
		if c.isDemoModeValid(res) {
			if c.onlyoffice.Onlyoffice.Demo.DocumentServerSecret == "" {
				c.logger.Errorf("demo mode is enabled but demo secret is not configured")
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			}

			jwtSecret = c.onlyoffice.Onlyoffice.Demo.DocumentServerSecret
		} else {
			if res.DocSecret == "" {
				c.logger.Errorf("no document server secret found and demo mode not valid (company %s)", cid)
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			}

			jwtSecret = res.DocSecret
		}

		if err := c.jwtManager.Verify(jwtSecret, body.Token, &body); err != nil {
			c.logger.Errorf("could not verify callback jwt (%s). Reason: %s", body.Token, err.Error())
			rw.WriteHeader(http.StatusForbidden)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if err := body.Validate(); err != nil {
			c.logger.Errorf("invalid callback body. Reason: %s", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(response.CallbackResponse{
				Error: 1,
			}.ToJSON())
			return
		}

		if body.Status == 2 {
			filename := strings.TrimSpace(r.URL.Query().Get("filename"))
			if filename == "" {
				rw.WriteHeader(http.StatusInternalServerError)
				c.logger.Errorf("callback request %s does not contain a filename", body.Key)
				rw.Write(response.CallbackResponse{
					Error: 1,
				}.ToJSON())
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(c.onlyoffice.Onlyoffice.Callback.UploadTimeout)*time.Second)
			defer cancel()

			usr := body.Users[0]
			if usr != "" {
				size, err := c.pipedriveAPI.ValidateFileSize(ctx, c.onlyoffice.Onlyoffice.Callback.MaxSize, body.URL)
				if err != nil {
					c.logger.Errorf("could not validate file %s: %s", filename, err.Error())
					rw.WriteHeader(http.StatusBadRequest)
					rw.Write(response.CallbackResponse{
						Error: 1,
					}.ToJSON())
					return
				}

				req := c.client.NewRequest(fmt.Sprintf("%s:auth", c.config.Namespace), "UserSelectHandler.GetUser", usr)
				var ures response.UserResponse
				if err := c.client.Call(ctx, req, &ures, client.WithRetries(3), client.WithBackoff(func(ctx context.Context, req client.Request, attempts int) (time.Duration, error) {
					return backoff.Do(attempts), nil
				})); err != nil {
					c.logger.Errorf("could not get user tokens: %s", err.Error())
					rw.WriteHeader(http.StatusBadRequest)
					rw.Write(response.CallbackResponse{
						Error: 1,
					}.ToJSON())
					return
				}

				if err := c.pipedriveAPI.UploadFile(ctx, body.URL, did, fid, filename, size, model.Token{
					AccessToken:  ures.AccessToken,
					RefreshToken: ures.RefreshToken,
					TokenType:    ures.TokenType,
					Scope:        ures.Scope,
					ApiDomain:    ures.ApiDomain,
				}); err != nil {
					c.logger.Debugf("could not upload an onlyoffice file to pipedrive: %s", err.Error())
					rw.WriteHeader(http.StatusBadRequest)
					rw.Write(response.CallbackResponse{
						Error: 1,
					}.ToJSON())
					return
				}
			}
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(response.CallbackResponse{
			Error: 0,
		}.ToJSON())
	}
}
