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

package domain

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type DocSettings struct {
	CompanyID   string    `json:"company_id" mapstructure:"company_id"`
	DocAddress  string    `json:"doc_address" mapstructure:"doc_address"`
	DocSecret   string    `json:"doc_secret" mapstructure:"doc_secret"`
	DocHeader   string    `json:"doc_header" mapstructure:"doc_header"`
	DemoEnabled bool      `json:"demo_enabled" mapstructure:"demo_enabled"`
	DemoStarted time.Time `json:"demo_started" mapstructure:"demo_started"`
}

func (u DocSettings) ToJSON() []byte {
	buf, _ := json.Marshal(u)
	return buf
}

func (u *DocSettings) Validate() error {
	u.CompanyID = strings.TrimSpace(u.CompanyID)
	u.DocAddress = strings.TrimSpace(u.DocAddress)
	u.DocSecret = strings.TrimSpace(u.DocSecret)
	u.DocHeader = strings.TrimSpace(u.DocHeader)

	if u.CompanyID == "" {
		return &InvalidModelFieldError{
			Model:  "Docserver",
			Field:  "CompanyID",
			Reason: "Should not be empty",
		}
	}

	if u.DemoEnabled {
		if u.DemoStarted.IsZero() {
			u.DemoStarted = time.Now()
			return nil
		}

		staleDate := time.Now().AddDate(0, 0, -30)
		if u.DemoStarted.Before(staleDate) {
			return &InvalidModelFieldError{
				Model:  "Docserver",
				Field:  "Demo Started",
				Reason: "Demo period has expired (more than 30 days old)",
			}
		}
		return nil
	}

	hasCredentials := u.DocAddress != "" || u.DocSecret != "" || u.DocHeader != ""

	if hasCredentials {
		if u.DocAddress == "" {
			return &InvalidModelFieldError{
				Model:  "Docserver",
				Field:  "Document Address",
				Reason: "Required when other credentials are provided",
			}
		}

		if u.DocSecret == "" {
			return &InvalidModelFieldError{
				Model:  "Docserver",
				Field:  "Document Secret",
				Reason: "Required when other credentials are provided",
			}
		}

		if u.DocHeader == "" {
			return &InvalidModelFieldError{
				Model:  "Docserver",
				Field:  "Document Header",
				Reason: "Required when other credentials are provided",
			}
		}

		url, err := url.Parse(u.DocAddress)
		if err != nil {
			return &InvalidModelFieldError{
				Model:  "Docserver",
				Field:  "Document Address",
				Reason: err.Error(),
			}
		}

		u.DocAddress = fmt.Sprintf("%s://%s/%s", url.Scheme, url.Host, url.Path)
		for {
			if strings.LastIndex(u.DocAddress, "/") == len(u.DocAddress)-1 {
				u.DocAddress = u.DocAddress[:len(u.DocAddress)-1]
			} else {
				break
			}
		}

		u.DocAddress += "/"
	}

	return nil
}
