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
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/domain"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/core/port"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/google/uuid"
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
	accessService port.AICodeAccessService
	logger        log.Logger
}

func NewApiController(
	client client.Client,
	apiClient pclient.PipedriveApiClient,
	commandClient pclient.CommandClient,
	jwtManager crypto.JwtManager,
	serverConfig *config.ServerConfig,
	accessService port.AICodeAccessService,
	logger log.Logger,
) ApiController {
	return ApiController{
		client:        client,
		apiClient:     apiClient,
		commandClient: commandClient,
		jwtManager:    jwtManager,
		config:        serverConfig,
		logger:        logger,
		accessService: accessService,
	}
}

func (c *ApiController) writeJSONResponse(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data.(interface{ ToJSON() []byte }).ToJSON())
}

func (c *ApiController) writeErrorResponse(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status)
}

func (c *ApiController) extractPipedriveContext(r *http.Request) (request.PipedriveTokenContext, bool) {
	pctx, ok := r.Context().Value("X-Pipedrive-App-Context").(request.PipedriveTokenContext)
	if !ok {
		c.logger.Error("could not extract pipedrive context from the context")
	}

	return pctx, ok
}

func (c *ApiController) createTimeoutContext(r *http.Request, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), timeout)
}

func (c *ApiController) handleMicroError(err error, rw http.ResponseWriter, defaultStatus int) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		c.writeErrorResponse(rw, http.StatusRequestTimeout)
		return
	}

	microErr := response.MicroError{}
	if unmarshalErr := json.Unmarshal([]byte(err.Error()), &microErr); unmarshalErr != nil {
		c.writeErrorResponse(rw, defaultStatus)
		return
	}

	c.writeErrorResponse(rw, microErr.Code)
}

func (c *ApiController) createToken(ures response.UserResponse) model.Token {
	return model.Token{
		AccessToken:  ures.AccessToken,
		RefreshToken: ures.RefreshToken,
		TokenType:    ures.TokenType,
		Scope:        ures.Scope,
		ApiDomain:    ures.ApiDomain,
	}
}

func (c *ApiController) getUserID(pctx request.PipedriveTokenContext) string {
	return fmt.Sprint(pctx.UID + pctx.CID)
}

func (c *ApiController) validateUserAccess(ctx context.Context, token model.Token) error {
	user, err := c.apiClient.GetMe(ctx, token)
	if err != nil {
		return err
	}

	for _, access := range user.Access {
		if access.App == "global" && !access.Admin {
			return ErrNotAdmin
		}
	}
	return nil
}

func (c *ApiController) getUser(ctx context.Context, id string) (response.UserResponse, int, error) {
	var ures response.UserResponse
	err := c.client.Call(ctx, c.client.NewRequest(fmt.Sprintf("%s:auth", c.config.Namespace), "UserSelectHandler.GetUser", id), &ures)
	if err != nil {
		c.logger.Errorf("could not get user access info: %s", err.Error())
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return ures, http.StatusRequestTimeout, err
		}

		microErr := response.MicroError{}
		if unmarshalErr := json.Unmarshal([]byte(err.Error()), &microErr); unmarshalErr != nil {
			return ures, http.StatusUnauthorized, err
		}

		return ures, microErr.Code, err
	}

	return ures, http.StatusOK, nil
}

func (c *ApiController) getAccess(ctx context.Context, code string) (request.DataRequest, error) {
	access, err := c.accessService.GetCodeAccess(ctx, code)
	if err != nil {
		c.logger.Errorf("could not get AI access code: %s", err.Error())
		return request.DataRequest{}, err
	}

	userID, _ := strconv.Atoi(access.UserID)
	return request.DataRequest{
		UserID: userID,
		DealID: access.DealID,
	}, nil
}

func (c *ApiController) regenerateAccess(ctx context.Context, userID, dealID string) (string, error) {
	newCode := uuid.New().String()
	access := domain.AICodeAccess{
		Code:   newCode,
		UserID: userID,
		DealID: dealID,
	}

	if err := c.accessService.UpsertCodeAccess(ctx, access); err != nil {
		c.logger.Errorf("could not regenerate AI access code: %s", err.Error())
		return "", err
	}
	return newCode, nil
}

func (c *ApiController) BuildGetMe() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		pctx, ok := c.extractPipedriveContext(r)
		if !ok {
			c.writeErrorResponse(rw, http.StatusForbidden)
			return
		}

		ctx, cancel := c.createTimeoutContext(r, 5*time.Second)
		defer cancel()

		ures, status, _ := c.getUser(ctx, c.getUserID(pctx))
		if status != http.StatusOK {
			c.writeErrorResponse(rw, status)
			return
		}

		c.writeJSONResponse(rw, response.UserTokenResponse{
			ID:          ures.ID,
			AccessToken: ures.AccessToken,
			ExpiresAt:   ures.ExpiresAt,
		})
	}
}

func (c *ApiController) BuildGetData() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		ctx, cancel := c.createTimeoutContext(r, 5*time.Second)
		defer cancel()

		data, err := c.getAccess(ctx, code)
		if err != nil {
			c.writeErrorResponse(rw, http.StatusNotFound)
			return
		}

		ures, _, err := c.getUser(ctx, fmt.Sprint(data.UserID))
		if err != nil {
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		token := c.createToken(ures)
		user, err := c.apiClient.GetMe(ctx, token)
		if err != nil {
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		deal, err := c.apiClient.GetDeal(ctx, fmt.Sprint(data.DealID), token)
		if err != nil {
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		nCode, err := c.regenerateAccess(ctx, fmt.Sprint(data.UserID), data.DealID)
		if err != nil {
			c.writeErrorResponse(rw, http.StatusInternalServerError)
			return
		}

		c.writeJSONResponse(rw, response.DataResponse{
			Data: map[string]any{
				"user": user,
				"deal": deal,
			},
			Code: nCode,
		})
	}
}

func (c *ApiController) BuildPostSettings() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		pctx, ok := c.extractPipedriveContext(r)
		if !ok {
			c.writeErrorResponse(rw, http.StatusForbidden)
			return
		}

		contentLen, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 0)
		if err != nil || (contentLen/100000) > 10 {
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		var settings request.DocSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			c.logger.Errorf("decode error: %s", err.Error())
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		settings.CompanyID = pctx.CID
		if err := settings.Validate(); err != nil {
			c.logger.Errorf("invalid settings format: %s", err.Error())
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		ctx, cancel := c.createTimeoutContext(r, 10*time.Second)
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
				ures, _, err := c.getUser(ectx, c.getUserID(pctx))
				if err != nil {
					return err
				}

				token := c.createToken(ures)
				if err := c.validateUserAccess(ectx, token); err != nil {
					c.logger.Errorf("user validation failed: %s", err.Error())
					return err
				}

				urs, err := c.apiClient.GetMe(ectx, token)
				if err != nil {
					c.logger.Errorf("could not get pipedrive user: %s", err.Error())
					return err
				}

				atomic.StoreInt64(&companyID, int64(urs.CompanyID))
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			c.logger.Errorf("validation failed: %s", err.Error())
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				c.writeErrorResponse(rw, http.StatusRequestTimeout)
			} else if errors.Is(err, ErrNotAdmin) {
				c.writeErrorResponse(rw, http.StatusForbidden)
			} else {
				c.writeErrorResponse(rw, http.StatusForbidden)
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

		var sres any
		err = c.client.Call(
			tctx,
			c.client.NewRequest(
				fmt.Sprintf("%s:settings", c.config.Namespace),
				"SettingsInsertHandler.InsertSettings",
				sreq,
			), &sres)
		if err != nil {
			c.logger.Errorf("could not post new settings: %s", err.Error())
			c.handleMicroError(err, rw, http.StatusUnauthorized)
			return
		}

		c.writeErrorResponse(rw, http.StatusCreated)
	}
}

func (c *ApiController) BuildGetSettings() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		pctx, ok := c.extractPipedriveContext(r)
		if !ok {
			c.writeErrorResponse(rw, http.StatusForbidden)
			return
		}

		ctx, cancel := c.createTimeoutContext(r, 3*time.Second)
		defer cancel()

		ures, status, _ := c.getUser(ctx, c.getUserID(pctx))
		if status != http.StatusOK {
			c.writeErrorResponse(rw, status)
			return
		}

		token := c.createToken(ures)
		if err := c.validateUserAccess(ctx, token); err != nil {
			c.writeErrorResponse(rw, http.StatusForbidden)
			return
		}

		var docs response.DocSettingsResponse
		err := c.client.Call(
			ctx,
			c.client.NewRequest(
				fmt.Sprintf("%s:settings", c.config.Namespace),
				"SettingsSelectHandler.GetSettings",
				fmt.Sprint(pctx.CID),
			),
			&docs,
		)

		if err != nil {
			c.logger.Errorf("could not get settings: %s", err.Error())
			c.handleMicroError(err, rw, http.StatusUnauthorized)
			return
		}

		c.writeJSONResponse(rw, docs)
	}
}

func (c *ApiController) BuildGetConfig() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		pctx, ok := c.extractPipedriveContext(r)
		if !ok {
			c.writeErrorResponse(rw, http.StatusForbidden)
			return
		}

		query := r.URL.Query()
		id := strings.TrimSpace(query.Get("id"))
		filename := strings.TrimSpace(query.Get("name"))
		key := strings.TrimSpace(query.Get("key"))
		dealID := strings.TrimSpace(query.Get("deal_id"))
		dark := query.Get("dark") == "true"

		if filename == "" {
			c.logger.Error("could not extract file name from URL Query")
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		if len(filename) > 200 {
			c.logger.Error("file length is greater than 200")
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		if key == "" {
			c.logger.Error("could not extract doc key from URL Query")
			c.writeErrorResponse(rw, http.StatusBadRequest)
			return
		}

		ctx, cancel := c.createTimeoutContext(r, 6*time.Second)
		defer cancel()

		code := uuid.New().String()
		access := domain.AICodeAccess{
			Code:   code,
			UserID: c.getUserID(pctx),
			DealID: dealID,
		}

		if err := c.accessService.UpsertCodeAccess(ctx, access); err != nil {
			c.logger.Errorf("could not store AI access code: %s", err.Error())
			c.writeErrorResponse(rw, http.StatusInternalServerError)
			return
		}

		var resp response.BuildConfigResponse
		err := c.client.Call(
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
		)

		if err != nil {
			c.logger.Errorf("could not build onlyoffice config: %s", err.Error())
			c.handleMicroError(err, rw, http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		c.writeJSONResponse(rw, resp)
	}
}
