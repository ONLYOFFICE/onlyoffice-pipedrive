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
	"errors"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/controller/view"
	pclient "github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/client"
)

const installURL = "/oauth/install"

func upstreamErrorPage(err error) view.ErrorPage {
	if errors.Is(err, context.DeadlineExceeded) {
		return view.ErrorPage{Status: http.StatusGatewayTimeout, RetryURL: installURL}
	}

	return view.ErrorPage{Status: http.StatusBadGateway, RetryURL: installURL}
}

func tokenErrorPage(err error) view.ErrorPage {
	var serr *pclient.UnexpectedStatusCodeError
	if errors.As(err, &serr) {
		switch serr.Code {
		case http.StatusUnauthorized, http.StatusForbidden:
			return view.ErrorPage{
				Status:       http.StatusInternalServerError,
				Title:        "The app is not configured correctly",
				Message:      "Pipedrive rejected the app credentials. Please contact the app administrator.",
				Illustration: view.IllustrationConfiguration,
			}
		case http.StatusBadRequest:
			return view.ErrorPage{
				Status:   http.StatusBadRequest,
				Title:    "Authorization link expired",
				Message:  "The authorization code is invalid or has already been used. Please install the app again.",
				RetryURL: installURL,
			}
		}
	}

	if errors.Is(err, pclient.ErrInvalidUrlFormat) {
		return view.ErrorPage{
			Status:       http.StatusInternalServerError,
			Title:        "The app is not configured correctly",
			Message:      "The OAuth redirect URL is invalid. Please contact the app administrator.",
			Illustration: view.IllustrationConfiguration,
		}
	}

	return upstreamErrorPage(err)
}
