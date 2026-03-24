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
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"go-micro.dev/v4/logger"
	"golang.org/x/oauth2"
)

type AuthMiddleware struct {
	credentials *oauth2.Config
	logger      log.Logger
}

func BuildHandleAuthMiddleware(
	credentials *oauth2.Config,
	logger log.Logger,
) AuthMiddleware {
	return AuthMiddleware{
		credentials: credentials,
		logger:      logger,
	}
}

func (m AuthMiddleware) Protect(next http.Handler) http.HandlerFunc {
	m.logger.Debugf("pipedrive event middleware has been built with client_id %s and client_secret %s", m.credentials.ClientID, m.credentials.ClientSecret)
	return func(rw http.ResponseWriter, r *http.Request) {
		signature := strings.ReplaceAll(r.Header.Get("Authorization"), "Basic ", "")
		if signature == "" {
			m.logger.Errorf("an unauthorized access to deauth endpoint")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		validSignature := base64.StdEncoding.
			EncodeToString([]byte(fmt.Sprintf("%s:%s", m.credentials.ClientID, m.credentials.ClientSecret)))
		if signature != validSignature {
			logger.Errorf("invalid uninstall signature")
			logger.Debugf("valid signature is %s whereas signature is %s", validSignature, signature)
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	}
}
