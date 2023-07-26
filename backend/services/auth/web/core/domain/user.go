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

package domain

import (
	"encoding/json"
	"strings"
)

type UserAccess struct {
	ID           string `json:"id" mapstructure:"id"`
	AccessToken  string `json:"access_token" mapstructure:"access_token"`
	RefreshToken string `json:"refresh_token" mapstructure:"refresh_token"`
	TokenType    string `json:"token_type" mapstructure:"token_type"`
	Scope        string `json:"scope" mapstructure:"scope"`
	ExpiresAt    int64  `json:"expires_at" mapstructure:"expires_at"`
	ApiDomain    string `json:"api_domain" mapstructure:"api_domain"`
}

func (u UserAccess) ToJSON() []byte {
	buf, _ := json.Marshal(u)
	return buf
}

func (u *UserAccess) Validate() error {
	u.ID = strings.TrimSpace(u.ID)
	u.AccessToken = strings.TrimSpace(u.AccessToken)
	u.RefreshToken = strings.TrimSpace(u.RefreshToken)
	u.TokenType = strings.TrimSpace(u.TokenType)
	u.Scope = strings.TrimSpace(u.Scope)
	u.ApiDomain = strings.TrimSpace(u.ApiDomain)

	if u.ID == "" {
		return &InvalidModelFieldError{
			Model:  "User",
			Field:  "ID",
			Reason: "Should not be empty",
		}
	}

	if u.AccessToken == "" {
		return &InvalidModelFieldError{
			Model:  "User",
			Field:  "OAuth Access Token",
			Reason: "Should not be empty",
		}
	}

	if u.RefreshToken == "" {
		return &InvalidModelFieldError{
			Model:  "User",
			Field:  "OAuth Refresh Token",
			Reason: "Should not be empty",
		}
	}

	if u.TokenType == "" {
		return &InvalidModelFieldError{
			Model:  "User",
			Field:  "OAuth Token Type",
			Reason: "Should not be empty",
		}
	}

	if u.Scope == "" {
		return &InvalidModelFieldError{
			Model:  "User",
			Field:  "OAuth Scope",
			Reason: "Should not be empty",
		}
	}

	if u.ExpiresAt < 1 {
		return &InvalidModelFieldError{
			Model:  "User",
			Field:  "OAuth ExpiresAt",
			Reason: "Invalid expiresAt value",
		}
	}

	if u.ApiDomain == "" {
		return &InvalidModelFieldError{
			Model:  "User",
			Field:  "OAuth Domain",
			Reason: "Should not be empty",
		}
	}

	return nil
}
