/**
 *
 * (c) Copyright Ascensio System SIA 2023
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

package web

import (
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/config"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	chttp "github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/service/http"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/callback/web/controller"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/shared"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/client"
)

type CallbackService struct {
	namespace     string
	mux           *chi.Mux
	client        client.Client
	logger        log.Logger
	maxSize       int64
	uploadTimeout int
}

// ApplyMiddleware useed to apply http server middlewares.
func (s CallbackService) ApplyMiddleware(middlewares ...func(http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}

// NewService initializes http server with options.
func NewServer(
	serverConfig *config.ServerConfig,
	workerConfig *config.WorkerConfig,
	onlyofficeConfig *shared.OnlyofficeConfig,
	logger log.Logger,
) chttp.ServerEngine {
	gin.SetMode(gin.ReleaseMode)

	service := CallbackService{
		namespace:     serverConfig.Namespace,
		mux:           chi.NewRouter(),
		logger:        logger,
		maxSize:       onlyofficeConfig.Onlyoffice.Callback.MaxSize,
		uploadTimeout: onlyofficeConfig.Onlyoffice.Callback.UploadTimeout,
	}

	return service
}

// NewHandler returns http server engine.
func (s CallbackService) NewHandler(client client.Client, cache cache.Cache) interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
} {
	return s.InitializeServer(client)
}

// InitializeServer sets all injected dependencies.
func (s *CallbackService) InitializeServer(c client.Client) *chi.Mux {
	s.client = c
	s.InitializeRoutes()
	return s.mux
}

// InitializeRoutes builds all http routes.
func (s *CallbackService) InitializeRoutes() {
	callbackController := controller.NewCallbackController(s.namespace, s.maxSize, s.uploadTimeout, s.logger, s.client)
	s.mux.Group(func(r chi.Router) {
		r.Use(chimiddleware.Recoverer)
		r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
			http.Redirect(rw, r.WithContext(r.Context()), "https://onlyoffice.com", http.StatusMovedPermanently)
		})
		r.Post("/callback", callbackController.BuildPostHandleCallback())
	})
}
