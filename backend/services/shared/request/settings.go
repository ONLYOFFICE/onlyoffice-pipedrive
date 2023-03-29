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

package request

import (
	"encoding/json"
	"errors"
	"strings"
)

var _ErrInvalidCompanyID = errors.New("invalid company id")
var _ErrInvalidDocAddress = errors.New("invalid doc server address")
var _ErrInvalidDocSecret = errors.New("invalid doc server secret")
var _ErrInvalidDocHeader = errors.New("invalid doc server header")

type DocSettings struct {
	CompanyID  int    `json:"company_id" mapstructure:"company_id"`
	DocAddress string `json:"doc_address" mapstructure:"doc_address"`
	DocSecret  string `json:"doc_secret" mapstructure:"doc_secret"`
	DocHeader  string `json:"doc_header" mapstructure:"doc_header"`
}

func (c DocSettings) ToJSON() []byte {
	buf, _ := json.Marshal(c)
	return buf
}

func (c DocSettings) Validate() error {
	c.DocAddress = strings.TrimSpace(c.DocAddress)
	c.DocSecret = strings.TrimSpace(c.DocSecret)

	if c.CompanyID <= 0 {
		return _ErrInvalidCompanyID
	}

	if c.DocAddress == "" {
		return _ErrInvalidDocAddress
	}

	if c.DocSecret == "" {
		return _ErrInvalidDocSecret
	}

	if c.DocHeader == "" {
		return _ErrInvalidDocHeader
	}

	return nil
}
