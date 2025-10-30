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
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/config"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/assets"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"go-micro.dev/v4/client"
)

type FileController struct {
	client     client.Client
	apiClient  pclient.PipedriveApiClient
	jwtManager crypto.JwtManager
	config     *config.ServerConfig
	onlyoffice *shared.OnlyofficeConfig
	logger     log.Logger
}

func NewFileController(
	client client.Client,
	apiClient pclient.PipedriveApiClient,
	jwtManager crypto.JwtManager,
	config *config.ServerConfig,
	onlyoffice *shared.OnlyofficeConfig,
	logger log.Logger,
) FileController {
	return FileController{
		client:     client,
		apiClient:  apiClient,
		jwtManager: jwtManager,
		config:     config,
		onlyoffice: onlyoffice,
		logger:     logger,
	}
}

func (c *FileController) getUser(ctx context.Context, id string) (response.UserResponse, int) {
	var ures response.UserResponse
	if err := c.client.Call(
		ctx,
		c.client.NewRequest(
			fmt.Sprintf("%s:auth", c.config.Namespace),
			"UserSelectHandler.GetUser",
			id,
		),
		&ures,
	); err != nil {
		c.logger.Errorf("could not get user access info: %s", err.Error())
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return ures, http.StatusRequestTimeout
		}

		microErr := response.MicroError{}
		if err := json.Unmarshal([]byte(err.Error()), &microErr); err != nil {
			return ures, http.StatusUnauthorized
		}

		return ures, microErr.Code
	}

	return ures, http.StatusOK
}

func (c FileController) BuildGetFile() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		lang, fileType, dealID, filename := strings.TrimSpace(query.Get("lang")),
			strings.TrimSpace(query.Get("type")), strings.TrimSpace(query.Get("deal")),
			strings.TrimSpace(query.Get("filename"))
		if lang == "" || fileType == "" || dealID == "" || filename == "" {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(filename) > 0 {
			lastDot := strings.LastIndex(filename, ".")
			if lastDot > 0 {
				baseName := strings.TrimSpace(filename[:lastDot])
				if baseName == "" {
					filename = fmt.Sprintf("New Document.%s", fileType)
				}
			}
		}

		pctx, ok := r.Context().Value("X-Pipedrive-App-Context").(request.PipedriveTokenContext)
		if !ok {
			rw.WriteHeader(http.StatusForbidden)
			c.logger.Error("could not extract pipedrive context from the context")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		ures, status := c.getUser(ctx, fmt.Sprint(pctx.UID+pctx.CID))
		if status != http.StatusOK {
			rw.WriteHeader(status)
			return
		}

		file, err := assets.Files.Open(fmt.Sprintf("assets/%s/new.%s", lang, fileType))
		if err != nil {
			lang = "default"
			file, err = assets.Files.Open(fmt.Sprintf("assets/%s/new.%s", lang, fileType))
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				c.logger.Errorf("could not get a new file: %s", err.Error())
				return
			}

			defer file.Close()
			res, ferr := c.apiClient.CreateFile(ctx, dealID, filename, file, model.Token{
				AccessToken:  ures.AccessToken,
				RefreshToken: ures.AccessToken,
				TokenType:    ures.TokenType,
				Scope:        ures.Scope,
				ApiDomain:    ures.ApiDomain,
			})

			if ferr != nil {
				rw.WriteHeader(http.StatusBadRequest)
				c.logger.Errorf("could not upload a pipedrive file: %s", ferr.Error())
				return
			}

			rw.Write(res.ToJSON())
			return
		}

		defer file.Close()
		res, ferr := c.apiClient.CreateFile(ctx, dealID, filename, file, model.Token{
			AccessToken:  ures.AccessToken,
			RefreshToken: ures.AccessToken,
			TokenType:    ures.TokenType,
			Scope:        ures.Scope,
			ApiDomain:    ures.ApiDomain,
		})

		if ferr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			c.logger.Errorf("could not upload a pipedrive file: %s", ferr.Error())
			return
		}

		rw.Write(res.ToJSON())
	}
}

func (c FileController) BuildGetDownloadUrl() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "plain/text")
		query := r.URL.Query()
		domain, fileID := query.Get("domain"), query.Get("file_id")

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		dreq, _ := http.NewRequest("GET", fmt.Sprintf("%s/files/%s/download", strings.TrimSuffix(domain, "/"), fileID), nil)
		dreq.Header.Add("Authorization", r.Header.Get("Authorization"))
		resp, err := client.Do(dreq)
		if err != nil {
			c.logger.Errorf("could not build a new download url: %s", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if resp != nil {
			defer resp.Body.Close()
		}

		if resp.StatusCode != 302 {
			c.logger.Errorf("unexpected status code while building a new download url: %d", resp.StatusCode)
			rw.WriteHeader(resp.StatusCode)
			return
		}

		rw.Write([]byte(resp.Header.Get("Location")))
	}
}
