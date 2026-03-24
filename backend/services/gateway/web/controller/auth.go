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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"go-micro.dev/v4/client"
	"golang.org/x/oauth2"
)

type AuthController struct {
	client        client.Client
	pipedriveAuth pclient.PipedriveAuthClient
	pipedriveAPI  pclient.PipedriveApiClient
	config        *config.ServerConfig
	credentials   *oauth2.Config
	logger        log.Logger
}

func NewAuthController(
	client client.Client,
	pipedriveAuth pclient.PipedriveAuthClient,
	pipedriveAPI pclient.PipedriveApiClient,
	config *config.ServerConfig,
	credentials *oauth2.Config,
	logger log.Logger,
) AuthController {
	return AuthController{
		client:        client,
		pipedriveAuth: pipedriveAuth,
		pipedriveAPI:  pipedriveAPI,
		config:        config,
		credentials:   credentials,
		logger:        logger,
	}
}

func (c AuthController) BuildGetInstall() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.logger.Debug("a new install request")
		http.Redirect(
			rw, r,
			fmt.Sprintf(
				"https://oauth.pipedrive.com/oauth/authorize?client_id=%s&redirect_uri=%s",
				c.credentials.ClientID,
				url.QueryEscape(c.credentials.RedirectURL),
			),
			http.StatusMovedPermanently,
		)
	}
}

func (c AuthController) BuildGetAuth() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.logger.Debug("a new auth request")
		code := strings.TrimSpace(r.URL.Query().Get("code"))
		if code == "" {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Debug("empty auth code parameter")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()

		token, err := c.pipedriveAuth.GetAccessToken(ctx, code, c.credentials.RedirectURL)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Errorf("could not get pipedrive access token: %s", err.Error())
			return
		}

		usr, err := c.pipedriveAPI.GetMe(ctx, token)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Errorf("could not get pipedrive user: %s", err.Error())
			return
		}

		var ures response.UserResponse
		if err := c.client.Call(
			ctx,
			c.client.NewRequest(
				fmt.Sprintf("%s:auth", c.config.Namespace),
				"UserInsertHandler.InsertUser",
				response.UserResponse{
					ID:           fmt.Sprint(usr.ID + usr.CompanyID),
					AccessToken:  token.AccessToken,
					RefreshToken: token.RefreshToken,
					TokenType:    token.TokenType,
					Scope:        token.Scope,
					ApiDomain:    token.ApiDomain,
					ExpiresAt:    time.Now().Local().Add(time.Second * time.Duration(token.ExpiresIn-700)).UnixMilli(),
				},
			),
			&ures,
		); err != nil {
			c.logger.Errorf("could not get user access info: %s", err.Error())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				rw.WriteHeader(http.StatusRequestTimeout)
				return
			}

			microErr := response.MicroError{}
			if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			rw.WriteHeader(microErr.Code)
			return
		}

		c.logger.Debugf("redirecting to api domain: %s", token.ApiDomain)
		http.Redirect(rw, r, token.ApiDomain, http.StatusMovedPermanently)
	}
}

func (c AuthController) BuildDeleteAuth() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.logger.Debug("a new uninstall request")
		var ureq request.UninstallRequest

		len, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 0)
		if err != nil || (len/100000) > 10 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&ureq); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			c.logger.Errorf("could not unmarshal request body: %s", err.Error())
			return
		}

		var res interface{}
		if err := c.client.Call(
			r.Context(),
			c.client.NewRequest(
				fmt.Sprintf("%s:auth", c.config.Namespace),
				"UserDeleteHandler.DeleteUser",
				fmt.Sprint(ureq.UserID+ureq.CompanyID),
			),
			&res,
		); err != nil {
			c.logger.Errorf("could not delete user: %s", err.Error())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				rw.WriteHeader(http.StatusRequestTimeout)
				return
			}

			microErr := response.MicroError{}
			if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			rw.WriteHeader(microErr.Code)
			return
		}

		c.logger.Debugf("successfully published delete-auth message for user %s", fmt.Sprint(ureq.UserID+ureq.CompanyID))
		rw.WriteHeader(http.StatusOK)
	}
}
