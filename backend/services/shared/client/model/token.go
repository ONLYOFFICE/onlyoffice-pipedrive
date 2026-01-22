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

package model

import "strings"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ApiDomain    string `json:"api_domain"`
	ExpiresIn    int    `json:"expires_in"`
}

func (t *Token) Validate() error {
	t.AccessToken = strings.TrimSpace(t.AccessToken)
	t.RefreshToken = strings.TrimSpace(t.RefreshToken)
	t.TokenType = strings.TrimSpace(t.TokenType)
	t.ApiDomain = strings.TrimSpace(t.ApiDomain)
	t.Scope = strings.TrimSpace(t.Scope)

	if t.AccessToken == "" {
		return ErrInvalidTokenFormat
	}

	if t.RefreshToken == "" {
		return ErrInvalidTokenFormat
	}

	if t.TokenType == "" {
		return ErrInvalidTokenFormat
	}

	if t.ApiDomain == "" {
		return ErrInvalidTokenFormat
	}

	if t.Scope == "" {
		return ErrInvalidTokenFormat
	}

	if t.ExpiresIn < 1 {
		return ErrInvalidTokenFormat
	}

	return nil
}
