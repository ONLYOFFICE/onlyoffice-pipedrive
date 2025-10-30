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

package response

import (
	"encoding/json"
	"time"
)

type DocSettingsResponse struct {
	DocAddress  string    `json:"doc_address"`
	DocSecret   string    `json:"doc_secret"`
	DocHeader   string    `json:"doc_header"`
	DemoEnabled bool      `json:"demo_enabled"`
	DemoStarted time.Time `json:"demo_started"`
}

func (r DocSettingsResponse) ToJSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

type SettingsConfiguredResponse struct {
	Configured bool `json:"configured"`
}

func (r SettingsConfiguredResponse) ToJSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}
