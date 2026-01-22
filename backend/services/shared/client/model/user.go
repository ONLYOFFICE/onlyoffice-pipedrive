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

type User struct {
	ID        int      `json:"id"`
	CompanyID int      `json:"company_id" mapstructure:"company_id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Language  Language `json:"language" mapstructure:"language"`
	Access    []Access `json:"access" mapstructure:"access"`
}

func (u *User) Validate() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)

	if u.Name == "" {
		return ErrInvalidTokenFormat
	}

	if err := u.Language.Validate(); err != nil {
		return err
	}

	return nil
}
