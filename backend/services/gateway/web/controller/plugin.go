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
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/gateway/assets"
)

type PluginController struct {
	logger log.Logger
}

func NewPluginController(logger log.Logger) PluginController {
	return PluginController{
		logger: logger,
	}
}

func (c PluginController) BuildServePlugin() http.HandlerFunc {
	pluginFS, err := fs.Sub(assets.PluginFiles, "aiautofill/build")
	if err != nil {
		c.logger.Errorf("could not create plugin filesystem: %s", err.Error())
		return func(rw http.ResponseWriter, r *http.Request) {
			http.Error(rw, "Plugin not available", http.StatusInternalServerError)
		}
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/plugins/aiautofill")
		if path == "" || path == "/" {
			path = "index.html"
		} else {
			path = strings.TrimPrefix(path, "/")
		}

		file, err := pluginFS.Open(path)
		if err != nil {
			c.logger.Errorf("could not open file %s: %s", path, err.Error())
			http.NotFound(rw, r)
			return
		}

		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.logger.Errorf("could not stat file %s: %s", path, err.Error())
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".html":
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".js":
			rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".json":
			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".css":
			rw.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".png":
			rw.Header().Set("Content-Type", "image/png")
		case ".svg":
			rw.Header().Set("Content-Type", "image/svg+xml")
		}

		http.ServeContent(rw, r, stat.Name(), stat.ModTime(), file.(io.ReadSeeker))
	}
}
