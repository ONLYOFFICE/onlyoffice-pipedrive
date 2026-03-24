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

package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var ErrCommandServiceError = errors.New("got a command service error 1 status")

type CommandClient struct {
	client     *resty.Client
	jwtManager crypto.JwtManager
}

func NewCommandClient(jwtManager crypto.JwtManager) CommandClient {
	otelClient := &http.Client{
		Transport: otelhttp.NewTransport(&http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 6 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}),
	}

	return CommandClient{
		client: resty.NewWithClient(otelClient).
			SetRetryCount(0).
			SetRetryWaitTime(120 * time.Millisecond).
			SetRetryMaxWaitTime(900 * time.Millisecond).
			SetLogger(log.NewEmptyLogger()).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			}),
		jwtManager: jwtManager,
	}
}

func (p *CommandClient) License(ctx context.Context, url, secret string) error {
	var resp response.BaseCommandResponse

	token, err := p.jwtManager.Sign(secret, request.BaseCommandRequest{
		C: "version",
	})

	if err != nil {
		return err
	}

	res, err := p.client.R().
		SetContext(ctx).
		SetBody(request.TokenCommandRequest{
			Token: token,
		}).
		SetResult(&resp).
		Post(fmt.Sprintf("%scommand?shardkey=%s", url, uuid.New().String()))

	if err != nil {
		return err
	}

	if res.StatusCode() >= 300 || resp.Error != 0 {
		return ErrCommandServiceError
	}

	return nil
}
