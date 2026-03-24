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

package middleware

import (
	"context"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/crypto"
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared/request"
	"go-micro.dev/v4/logger"
	"golang.org/x/oauth2"
)

type ContextMiddleware struct {
	jwtManager  crypto.JwtManager
	credentials *oauth2.Config
	logger      log.Logger
}

func BuildHandleContextMiddleware(
	jwtManager crypto.JwtManager,
	credentials *oauth2.Config,
	logger log.Logger,
) ContextMiddleware {
	return ContextMiddleware{
		jwtManager:  jwtManager,
		credentials: credentials,
		logger:      logger,
	}
}

func (m ContextMiddleware) Protect(next http.Handler) http.HandlerFunc {
	m.logger.Debugf("pipedrive context middleware has been built with client_secret %s", m.credentials.ClientSecret)
	return func(rw http.ResponseWriter, r *http.Request) {
		var ctx request.PipedriveTokenContext
		token := r.Header.Get("X-Pipedrive-App-Context")
		if token == "" {
			logger.Errorf("unauthorized access to an api endpoint")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := m.jwtManager.Verify(m.credentials.ClientSecret, token, &ctx); err != nil {
			logger.Errorf("could not verify X-Pipedrive-App-Context: %s", err.Error())
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), "X-Pipedrive-App-Context", ctx)))
	}
}
