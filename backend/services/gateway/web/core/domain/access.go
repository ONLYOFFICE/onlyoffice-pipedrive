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
	"strings"
)

type AICodeAccess struct {
	Code   string `json:"code" mapstructure:"code"`
	UserID string `json:"user_id" mapstructure:"user_id"`
	FileID string `json:"file_id" mapstructure:"file_id"`
	DealID string `json:"deal_id" mapstructure:"deal_id"`
}

func (u AICodeAccess) ToJSON() []byte {
	buf, _ := json.Marshal(u)
	return buf
}

func (u *AICodeAccess) Validate() error {
	u.Code = strings.TrimSpace(u.Code)
	u.FileID = strings.TrimSpace(u.FileID)
	u.UserID = strings.TrimSpace(u.UserID)
	u.DealID = strings.TrimSpace(u.DealID)

	if u.Code == "" {
		return &InvalidModelFieldError{
			Model:  "AICodeAccess",
			Field:  "Code",
			Reason: "Should not be empty",
		}
	}

	if u.UserID == "" {
		return &InvalidModelFieldError{
			Model:  "AICodeAccess",
			Field:  "UserID",
			Reason: "Should not be empty",
		}
	}

	if u.FileID == "" {
		return &InvalidModelFieldError{
			Model:  "AICodeAccess",
			Field:  "FileID",
			Reason: "Should not be empty",
		}
	}

	if u.DealID == "" {
		return &InvalidModelFieldError{
			Model:  "AICodeAccess",
			Field:  "DealID",
			Reason: "Should not be empty",
		}
	}

	return nil
}
