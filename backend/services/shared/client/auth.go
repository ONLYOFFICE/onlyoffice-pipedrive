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
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
)

type PipedriveAuthClient struct {
	client       *resty.Client
	clientID     string
	clientSecret string
}

func NewPipedriveAuthClient(credentials *oauth2.Config) PipedriveAuthClient {
	otelClient := &http.Client{
		Transport: otelhttp.NewTransport(&http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}),
	}

	return PipedriveAuthClient{
		client: resty.NewWithClient(otelClient).
			SetBaseURL("https://oauth.pipedrive.com").
			SetRetryCount(0).
			SetRetryWaitTime(1000 * time.Millisecond).
			SetRetryMaxWaitTime(1500 * time.Millisecond).
			SetLogger(log.NewEmptyLogger()).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			}),
		clientID:     credentials.ClientID,
		clientSecret: credentials.ClientSecret,
	}
}

func (c PipedriveAuthClient) GetAccessToken(ctx context.Context, code, redirectURI string) (model.Token, error) {
	var resp model.Token
	if _, err := url.ParseRequestURI(redirectURI); err != nil {
		return resp, ErrInvalidUrlFormat
	}

	res, err := c.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetContext(ctx).
		SetBody(strings.NewReader(url.Values{
			"grant_type":   []string{"authorization_code"},
			"code":         []string{code},
			"redirect_uri": []string{redirectURI},
		}.Encode())).
		SetBasicAuth(c.clientID, c.clientSecret).
		SetResult(&resp).
		Post("/oauth/token")

	if err != nil {
		return resp, err
	}

	if res.StatusCode() != http.StatusOK {
		return resp, &UnexpectedStatusCodeError{
			Action: "get access token",
			Code:   res.StatusCode(),
		}
	}

	return resp, resp.Validate()
}

func (c PipedriveAuthClient) RefreshAccessToken(ctx context.Context, refreshToken string) (model.Token, error) {
	var resp model.Token

	res, err := c.client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetContext(ctx).
		SetBody(strings.NewReader(url.Values{
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refreshToken},
		}.Encode())).
		SetBasicAuth(c.clientID, c.clientSecret).
		SetResult(&resp).
		Post("/oauth/token")

	if err != nil {
		return resp, err
	}

	if res.StatusCode() != http.StatusOK {
		return resp, &UnexpectedStatusCodeError{
			Action: "refresh access token",
			Code:   res.StatusCode(),
		}
	}

	return resp, resp.Validate()
}
