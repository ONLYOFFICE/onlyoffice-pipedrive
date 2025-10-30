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

package request

import (
	"encoding/json"
	"net/url"
	"strings"
)

type DocSettings struct {
	CompanyID   int    `json:"company_id" mapstructure:"company_id"`
	DocAddress  string `json:"doc_address" mapstructure:"doc_address"`
	DocSecret   string `json:"doc_secret" mapstructure:"doc_secret"`
	DocHeader   string `json:"doc_header" mapstructure:"doc_header"`
	DemoEnabled bool   `json:"demo_enabled" mapstructure:"demo_enabled"`
}

func (c DocSettings) ToJSON() []byte {
	buf, _ := json.Marshal(c)
	return buf
}

func (c DocSettings) Validate() error {
	c.DocAddress = strings.TrimSpace(c.DocAddress)
	c.DocSecret = strings.TrimSpace(c.DocSecret)
	c.DocHeader = strings.TrimSpace(c.DocHeader)

	if c.CompanyID <= 0 {
		return ErrInvalidCompanyID
	}

	hasCredentials := c.DocAddress != "" || c.DocSecret != "" || c.DocHeader != ""
	if hasCredentials {
		if c.DocAddress == "" {
			return ErrInvalidDocAddress
		}

		if c.DocSecret == "" {
			return ErrInvalidDocSecret
		}

		if c.DocHeader == "" {
			return ErrInvalidDocHeader
		}

		parsedURL, err := url.Parse(c.DocAddress)
		if err != nil {
			return ErrInvalidDocAddress
		}

		if parsedURL.Scheme == "http" {
			return ErrHttpNotAllowed
		}
	} else if c.DemoEnabled {
		return nil
	}

	return nil
}
