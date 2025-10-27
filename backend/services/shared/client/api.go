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

package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client/model"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/response"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type PipedriveApiClient struct {
	client *resty.Client
}

func NewPipedriveApiClient() PipedriveApiClient {
	otelClient := otelhttp.DefaultClient
	otelClient.Transport = otelhttp.NewTransport(&http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	return PipedriveApiClient{
		client: resty.NewWithClient(otelClient).
			SetRetryCount(3).
			SetRetryWaitTime(120 * time.Millisecond).
			SetRetryMaxWaitTime(900 * time.Millisecond).
			SetLogger(log.NewEmptyLogger()).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			}),
	}
}

func (p *PipedriveApiClient) GetOrganization(ctx context.Context, id string, token model.Token) (map[string]any, error) {
	var resp model.OrganizationResponse

	res, err := p.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetResult(&resp).
		Get(fmt.Sprintf("%s/api/v2/organizations/%s", token.ApiDomain, id))

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, &UnexpectedStatusCodeError{
			Action: "get organization",
			Code:   res.StatusCode(),
		}
	}

	return resp.Data, nil
}

func (p *PipedriveApiClient) GetPerson(ctx context.Context, id string, token model.Token) (model.PersonData, error) {
	var resp model.PersorResponse

	res, err := p.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetResult(&resp).
		Get(fmt.Sprintf("%s/api/v2/persons/%s", token.ApiDomain, id))

	if err != nil {
		return model.PersonData{}, err
	}

	if res.StatusCode() != http.StatusOK {
		return model.PersonData{}, &UnexpectedStatusCodeError{
			Action: "get person",
			Code:   res.StatusCode(),
		}
	}

	return resp.Data, nil
}

func (p *PipedriveApiClient) GetDealProducts(ctx context.Context, id string, token model.Token) ([]model.Product, error) {
	var resp model.ProductsResponse

	res, err := p.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetResult(&resp).
		Get(fmt.Sprintf("%s/api/v2/deals/%s/products", token.ApiDomain, id))

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, &UnexpectedStatusCodeError{
			Action: "get deal products",
			Code:   res.StatusCode(),
		}
	}

	if !resp.Success {
		return nil, &UnexpectedStatusCodeError{
			Action: "get deal products",
			Code:   http.StatusInternalServerError,
		}
	}

	// TODO: Handle pagination more gracefully
	if len(resp.Data) > 5 {
		return resp.Data[:5], nil
	}

	return resp.Data, nil
}

func (p *PipedriveApiClient) GetDeal(ctx context.Context, id string, token model.Token) (model.Deal, error) {
	var deal model.Deal
	var resp model.DealResponse

	res, err := p.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetResult(&resp).
		Get(fmt.Sprintf("%s/api/v2/deals/%s", token.ApiDomain, id))

	if err != nil {
		return deal, err
	}

	if res.StatusCode() != http.StatusOK {
		return deal, &UnexpectedStatusCodeError{
			Action: "get deal",
			Code:   res.StatusCode(),
		}
	}

	if !resp.Success {
		return deal, &UnexpectedStatusCodeError{
			Action: "get deal",
			Code:   http.StatusInternalServerError,
		}
	}

	return resp.Data, nil
}

func (p *PipedriveApiClient) GetMe(ctx context.Context, token model.Token) (model.User, error) {
	var usr model.User
	var resp interface{}

	res, err := p.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetResult(&resp).
		Get(fmt.Sprintf("%s/api/v1/users/me", token.ApiDomain))

	if err != nil {
		return usr, err
	}

	if res.StatusCode() != http.StatusOK {
		return usr, &UnexpectedStatusCodeError{
			Action: "get me",
			Code:   res.StatusCode(),
		}
	}

	m, ok := resp.(map[string]interface{})
	if !ok {
		return usr, &UnexpectedStatusCodeError{
			Action: "get me",
			Code:   http.StatusInternalServerError,
		}
	}

	if err := mapstructure.Decode(m["data"], &usr); err != nil {
		return usr, err
	}

	return usr, nil
}

func (p *PipedriveApiClient) UpdateFile(ctx context.Context, id, name string, token model.Token) error {
	res, err := p.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetFormData(map[string]string{
			"name": name,
		}).
		Put(fmt.Sprintf("%s/api/v1/files/%s", token.ApiDomain, id))

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return &UnexpectedStatusCodeError{
			Action: "update file",
			Code:   res.StatusCode(),
		}
	}

	return nil
}

func (c *PipedriveApiClient) ValidateFileSize(ctx context.Context, limit int64, url string) (int64, error) {
	headResp, err := c.client.R().
		SetContext(ctx).
		Head(url)

	if err != nil {
		return 0, err
	}

	if val, err := strconv.ParseInt(headResp.Header().Get("Content-Length"), 10, 0); val > limit || err != nil {
		return 0, ErrInvalidContentLength
	} else {
		return val, nil
	}
}

func (p PipedriveApiClient) getFile(ctx context.Context, url string) (io.ReadCloser, error) {
	fileResp, err := p.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return nil, err
	}

	if fileResp.RawResponse.StatusCode != http.StatusOK {
		fileResp.RawBody().Close()
		return nil, fmt.Errorf("unexpected status code: %d", fileResp.RawResponse.StatusCode)
	}

	return fileResp.RawBody(), nil
}

func (p *PipedriveApiClient) UploadFile(ctx context.Context, url, deal, fileID, filename string, size int64, token model.Token) error {
	if err := p.UpdateFile(ctx, fileID, filename, token); err != nil {
		return err
	}

	file, err := p.getFile(ctx, url)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = p.client.R().
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetFileReader("file", filename, file).
		SetFormData(map[string]string{
			"deal_id": deal,
		}).
		Post(fmt.Sprintf("%s/api/v1/files", token.ApiDomain))

	if err != nil {
		return err
	}

	return nil
}

func (p *PipedriveApiClient) CreateFile(ctx context.Context, deal, filename string, file io.ReadCloser, token model.Token) (response.AddFileResponse, error) {
	var body response.AddFileResponse

	_, err := p.client.R().
		SetResult(&body).
		SetContext(ctx).
		SetAuthToken(token.AccessToken).
		SetFileReader("file", filename, file).
		SetFormData(map[string]string{
			"deal_id": deal,
		}).
		Post(fmt.Sprintf("%s/api/v1/files", token.ApiDomain))

	if err != nil {
		return body, err
	}

	return body, nil
}
