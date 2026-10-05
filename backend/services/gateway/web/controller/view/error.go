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

package view

import (
	"bytes"
	"embed"
	"encoding/base64"
	"html/template"
	"net/http"
)

type Illustration string

const (
	IllustrationClient        Illustration = "ant"
	IllustrationServer        Illustration = "hedgehog"
	IllustrationConfiguration Illustration = "owl"
)

//go:embed error.html
var errorHTML string

//go:embed icons/*.svg
var icons embed.FS

var errorTemplate = template.Must(template.New("error").Parse(errorHTML))

var (
	logo    = mustReadIcon("onlyoffice")
	favicon = template.URL("data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(mustReadIcon("favicon"))))
)

func mustReadIcon(name string) template.HTML {
	icon, err := icons.ReadFile("icons/" + name + ".svg")
	if err != nil {
		panic(err)
	}

	return template.HTML(icon)
}

type ErrorPage struct {
	Status       int
	Title        string
	Message      string
	RetryURL     string
	Illustration Illustration
}

type errorPageData struct {
	ErrorPage
	Logo    template.HTML
	Icon    template.HTML
	Favicon template.URL
}

var defaultErrorPages = map[int]ErrorPage{
	http.StatusBadRequest: {
		Title:   "Bad request",
		Message: "The request could not be processed. Please try again.",
	},
	http.StatusUnauthorized: {
		Title:   "Authorization failed",
		Message: "We could not verify your Pipedrive account. Please reinstall the app and try again.",
	},
	http.StatusForbidden: {
		Title:   "Access denied",
		Message: "You do not have permission to access this page.",
	},
	http.StatusNotFound: {
		Title:   "Page not found",
		Message: "The page you are looking for does not exist.",
	},
	http.StatusRequestTimeout: {
		Title:   "Request timed out",
		Message: "The request took too long to complete. Please try again.",
	},
	http.StatusInternalServerError: {
		Title:   "Something went wrong",
		Message: "An unexpected error occurred on our side. Please try again later.",
	},
	http.StatusBadGateway: {
		Title:   "Pipedrive is unavailable",
		Message: "We could not reach Pipedrive. Please try again in a few minutes.",
	},
	http.StatusServiceUnavailable: {
		Title:   "Service unavailable",
		Message: "The service is temporarily unavailable. Please try again in a few minutes.",
	},
	http.StatusGatewayTimeout: {
		Title:   "Pipedrive did not respond",
		Message: "Pipedrive took too long to respond. Please try again in a few minutes.",
	},
}

func NewErrorPage(status int) ErrorPage {
	if status < 400 || status > 599 {
		status = http.StatusInternalServerError
	}

	page, ok := defaultErrorPages[status]
	if !ok {
		page = defaultErrorPages[http.StatusInternalServerError]
		if status < 500 {
			page = defaultErrorPages[http.StatusBadRequest]
		}
	}

	page.Status = status
	page.Illustration = IllustrationServer
	if status < 500 {
		page.Illustration = IllustrationClient
	}

	return page
}

func RenderError(rw http.ResponseWriter, page ErrorPage) {
	defaults := NewErrorPage(page.Status)
	page.Status = defaults.Status
	if page.Title == "" {
		page.Title = defaults.Title
	}

	if page.Message == "" {
		page.Message = defaults.Message
	}

	if page.Illustration == "" {
		page.Illustration = defaults.Illustration
	}

	icon, err := icons.ReadFile("icons/" + string(page.Illustration) + ".svg")
	if err != nil {
		icon, _ = icons.ReadFile("icons/" + string(defaults.Illustration) + ".svg")
	}

	var buf bytes.Buffer
	if err := errorTemplate.Execute(&buf, errorPageData{
		ErrorPage: page,
		Logo:      logo,
		Favicon:   favicon,
		Icon:      template.HTML(icon),
	}); err != nil {
		http.Error(rw, http.StatusText(page.Status), page.Status)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(page.Status)
	_, _ = rw.Write(buf.Bytes())
}
