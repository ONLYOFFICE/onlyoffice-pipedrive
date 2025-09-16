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

package web

import (
	"net/http"

	shttp "github.com/ONLYOFFICE/onlyoffice-integration-adapters/service/http"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/controller"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/web/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type PipedriveHTTPService struct {
	apiController     controller.ApiController
	authController    controller.AuthController
	fileController    controller.FileController
	authMiddleware    middleware.AuthMiddleware
	contextMiddleware middleware.ContextMiddleware
	mux               *chi.Mux
}

// NewService initializes http server with options.
func NewServer(
	apiController controller.ApiController,
	authController controller.AuthController,
	fileController controller.FileController,
	authMiddleware middleware.AuthMiddleware,
	contextMiddleware middleware.ContextMiddleware,
) shttp.ServerEngine {
	return PipedriveHTTPService{
		apiController:     apiController,
		authController:    authController,
		fileController:    fileController,
		authMiddleware:    authMiddleware,
		contextMiddleware: contextMiddleware,
		mux:               chi.NewRouter(),
	}
}

// ApplyMiddleware useed to apply http server middlewares.
func (s PipedriveHTTPService) ApplyMiddleware(middlewares ...func(http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}

// NewHandler returns http server engine.
func (s PipedriveHTTPService) NewHandler() interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
} {
	return s.InitializeServer()
}

// InitializeServer sets all injected dependencies.
func (s *PipedriveHTTPService) InitializeServer() *chi.Mux {
	s.InitializeRoutes()
	return s.mux
}

// InitializeRoutes builds all http routes.
func (s *PipedriveHTTPService) InitializeRoutes() {
	s.mux.Group(func(r chi.Router) {
		r.Use(chimiddleware.Recoverer)
		r.NotFound(func(rw http.ResponseWriter, cr *http.Request) {
			http.Redirect(rw, cr.WithContext(cr.Context()), "https://onlyoffice.com", http.StatusMovedPermanently)
		})

		r.Route("/oauth", func(cr chi.Router) {
			cr.Use(chimiddleware.NoCache)
			cr.Get("/install", s.authController.BuildGetInstall())
			cr.Get("/auth", s.authController.BuildGetAuth())
			cr.Delete("/auth", s.authMiddleware.Protect(s.authController.BuildDeleteAuth()))
		})

		r.Route("/api", func(cr chi.Router) {
			cr.Get("/data", s.apiController.BuildGetData())
			cr.Group(func(pr chi.Router) {
				pr.Use(func(h http.Handler) http.Handler {
					return s.contextMiddleware.Protect(h)
				})
				pr.Get("/me", s.apiController.BuildGetMe())
				pr.Get("/config", s.apiController.BuildGetConfig())
				pr.Post("/settings", s.apiController.BuildPostSettings())
				pr.Get("/settings", s.apiController.BuildGetSettings())
			})
		})

		r.Route("/files", func(fr chi.Router) {
			fr.Get("/download", s.fileController.BuildGetDownloadUrl())
			fr.Get("/create", s.contextMiddleware.Protect(s.fileController.BuildGetFile()))
		})
	})
}
