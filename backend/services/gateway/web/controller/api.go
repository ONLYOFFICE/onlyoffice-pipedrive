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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"go-micro.dev/v4/client"
	"golang.org/x/sync/errgroup"
)

var ErrNotAdmin = errors.New("no admin access")

type ApiController struct {
	client        client.Client
	apiClient     pclient.PipedriveApiClient
	commandClient pclient.CommandClient
	jwtManager    crypto.JwtManager
	config        *config.ServerConfig
	userCache     map[string]model.User
	dealCache     map[string]model.Deal
	logger        log.Logger
}

func NewApiController(
	client client.Client,
	apiClient pclient.PipedriveApiClient,
	commandClient pclient.CommandClient,
	jwtManager crypto.JwtManager,
	serverConfig *config.ServerConfig,
	logger log.Logger,
) ApiController {
	return ApiController{
		client:        client,
		apiClient:     apiClient,
		commandClient: commandClient,
		jwtManager:    jwtManager,
		config:        serverConfig,
		logger:        logger,
		userCache:     make(map[string]model.User),
		dealCache:     make(map[string]model.Deal),
	}
}

func (c *ApiController) getUser(ctx context.Context, id string) (response.UserResponse, int, error) {
	var ures response.UserResponse
	if err := c.client.Call(ctx, c.client.NewRequest(fmt.Sprintf("%s:auth", c.config.Namespace), "UserSelectHandler.GetUser", id), &ures); err != nil {
		c.logger.Errorf("could not get user access info: %s", err.Error())
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return ures, http.StatusRequestTimeout, err
		}

		microErr := response.MicroError{}
		if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
			return ures, http.StatusUnauthorized, err
		}

		return ures, microErr.Code, err
	}

	return ures, http.StatusOK, nil
}

func (c ApiController) BuildGetMe() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		pctx, ok := r.Context().Value("X-Pipedrive-App-Context").(request.PipedriveTokenContext)
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract pipedrive context from the context")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		ures, status, _ := c.getUser(ctx, fmt.Sprint(pctx.UID+pctx.CID))
		if status != http.StatusOK {
			rw.WriteHeader(status)
			return
		}

		rw.Write(response.UserTokenResponse{
			ID:          ures.ID,
			AccessToken: ures.AccessToken,
			ExpiresAt:   ures.ExpiresAt,
		}.ToJSON())
	}
}

type Data struct {
	Data any    `json:"data"`
	Code string `json:"code"`
}

func (c ApiController) BuildGetData() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		usr, ok := c.userCache[code]
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		deal, ok := c.dealCache[code]
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		data, err := json.Marshal(Data{
			Data: map[string]any{
				"user": usr,
				"deal": deal,
			},
			Code: code,
		})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		delete(c.userCache, code)
		delete(c.dealCache, code)
		c.userCache["1111"] = usr
		c.dealCache["1111"] = deal

		rw.Write(data)
	}
}

func (c ApiController) BuildPostSettings() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		pctx, ok := r.Context().Value("X-Pipedrive-App-Context").(request.PipedriveTokenContext)
		if !ok {
			c.logger.Error("could not extract pipedrive context from the context")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		len, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 0)
		if err != nil || (len/100000) > 10 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var settings request.DocSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Errorf(err.Error())
			return
		}

		settings.CompanyID = pctx.CID
		if err := settings.Validate(); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Errorf("invalid settings format: %s", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var companyID int64

		eg, ectx := errgroup.WithContext(ctx)

		if !settings.DemoEnabled {
			eg.Go(func() error {
				select {
				case <-ectx.Done():
					return ectx.Err()
				default:
					if err := c.commandClient.License(ectx, settings.DocAddress, settings.DocSecret); err != nil {
						c.logger.Errorf("could not validate ONLYOFFICE document server credentials: %s", err.Error())
						return err
					}
					return nil
				}
			})
		} else {
			c.logger.Debugf("skipping document server validation - demo mode enabled")
		}

		eg.Go(func() error {
			select {
			case <-ectx.Done():
				return ectx.Err()
			default:
				ures, _, err := c.getUser(ectx, fmt.Sprint(pctx.UID+pctx.CID))
				if err != nil {
					return err
				}

				urs, err := c.apiClient.GetMe(ectx, model.Token{
					AccessToken:  ures.AccessToken,
					RefreshToken: ures.RefreshToken,
					TokenType:    ures.TokenType,
					Scope:        ures.Scope,
					ApiDomain:    ures.ApiDomain,
				})
				if err != nil {
					c.logger.Errorf("could not get pipedrive user or no user has admin permissions")
					return err
				}

				for _, access := range urs.Access {
					if access.App == "global" && !access.Admin {
						return ErrNotAdmin
					}
				}

				atomic.StoreInt64(&companyID, int64(urs.CompanyID))
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			c.logger.Errorf("validation failed: %s", err.Error())
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				c.logger.Errorf("request timeout during validation")
				rw.WriteHeader(http.StatusRequestTimeout)
			} else if errors.Is(err, ErrNotAdmin) {
				c.logger.Errorf("user does not have admin permissions")
				rw.WriteHeader(http.StatusForbidden)
			} else {
				c.logger.Errorf("other validation error: %s", err.Error())
				rw.WriteHeader(http.StatusForbidden)
			}
			return
		}

		sreq := request.DocSettings{
			CompanyID:   int(atomic.LoadInt64(&companyID)),
			DocAddress:  settings.DocAddress,
			DocHeader:   settings.DocHeader,
			DocSecret:   settings.DocSecret,
			DemoEnabled: settings.DemoEnabled,
		}

		tctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var sres interface{}
		if err := c.client.Call(
			tctx,
			c.client.NewRequest(
				fmt.Sprintf("%s:settings", c.config.Namespace),
				"SettingsInsertHandler.InsertSettings",
				sreq,
			), &sres); err != nil {
			c.logger.Errorf("could not post new settings: %s", err.Error())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				rw.WriteHeader(http.StatusRequestTimeout)
				return
			}

			microErr := response.MicroError{}
			if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				c.logger.Errorf("could not post new settings: %s", err.Error())
				return
			}

			rw.WriteHeader(microErr.Code)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	}
}

func (c ApiController) BuildGetSettings() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		pctx, ok := r.Context().Value("X-Pipedrive-App-Context").(request.PipedriveTokenContext)
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract pipedrive context from the context")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		ures, status, _ := c.getUser(ctx, fmt.Sprint(pctx.UID+pctx.CID))
		if status != http.StatusOK {
			rw.WriteHeader(status)
			return
		}

		urs, _ := c.apiClient.GetMe(ctx, model.Token{
			AccessToken:  ures.AccessToken,
			RefreshToken: ures.RefreshToken,
			TokenType:    ures.TokenType,
			Scope:        ures.Scope,
			ApiDomain:    ures.ApiDomain,
		})

		for _, access := range urs.Access {
			if access.App == "global" && !access.Admin {
				rw.WriteHeader(http.StatusForbidden)
				return
			}
		}

		var docs response.DocSettingsResponse
		if err := c.client.Call(
			ctx,
			c.client.NewRequest(
				fmt.Sprintf("%s:settings", c.config.Namespace),
				"SettingsSelectHandler.GetSettings",
				fmt.Sprint(pctx.CID),
			),
			&docs,
		); err != nil {
			c.logger.Errorf("could not get settings: %s", err.Error())
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

		rw.Write(docs.ToJSON())
	}
}

func (c ApiController) BuildGetConfig() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		id, filename, key, dealID, dark := strings.TrimSpace(query.Get("id")), strings.TrimSpace(query.Get("name")),
			strings.TrimSpace(query.Get("key")), strings.TrimSpace(query.Get("deal_id")), query.Get("dark") == "true"

		pctx, ok := r.Context().Value("X-Pipedrive-App-Context").(request.PipedriveTokenContext)
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract pipedrive context from the context")
			return
		}

		if filename == "" {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Error("could not extract file name from URL Query")
			return
		}

		if len(filename) > 200 {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Error("file length is greater than 200")
			return
		}

		if key == "" {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Error("could not extract doc key from URL Query")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
		defer cancel()

		ures, status, _ := c.getUser(ctx, fmt.Sprint(pctx.UID+pctx.CID))
		if status != http.StatusOK {
			rw.WriteHeader(status)
			return
		}

		urs, err := c.apiClient.GetMe(ctx, model.Token{
			AccessToken:  ures.AccessToken,
			RefreshToken: ures.RefreshToken,
			TokenType:    ures.TokenType,
			Scope:        ures.Scope,
			ApiDomain:    ures.ApiDomain,
		})

		if err != nil {
			c.logger.Errorf("could not get pipedrive user or no user has admin permissions: %s", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		deal, err := c.apiClient.GetDeal(ctx, dealID, model.Token{
			AccessToken:  ures.AccessToken,
			RefreshToken: ures.RefreshToken,
			TokenType:    ures.TokenType,
			Scope:        ures.Scope,
			ApiDomain:    ures.ApiDomain,
		})

		if err != nil {
			c.logger.Errorf("could not get pipedrive deal: %s", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		code := "1111"
		c.userCache[code] = urs
		c.dealCache[code] = deal

		var resp response.BuildConfigResponse
		if err := c.client.Call(
			ctx,
			c.client.NewRequest(
				fmt.Sprintf("%s:builder", c.config.Namespace),
				"ConfigHandler.BuildConfig",
				request.BuildConfigRequest{
					UID:       pctx.UID,
					CID:       pctx.CID,
					Deal:      dealID,
					UserAgent: r.UserAgent(),
					Filename:  filename,
					FileID:    id,
					DocKey:    key,
					Dark:      dark,
					Code:      code,
				},
			),
			&resp,
		); err != nil {
			c.logger.Errorf("could not build onlyoffice config: %s", err.Error())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				rw.WriteHeader(http.StatusRequestTimeout)
				return
			}

			microErr := response.MicroError{}
			if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			c.logger.Errorf("build config micro error: %s", microErr.Detail)
			rw.WriteHeader(microErr.Code)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(resp.ToJSON())
	}
}
