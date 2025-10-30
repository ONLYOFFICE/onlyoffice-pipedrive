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

import "errors"

var (
	ErrInvalidCompanyID  = errors.New("invalid company id")
	ErrInvalidDocAddress = errors.New("invalid doc server address")
	ErrInvalidDocSecret  = errors.New("invalid doc server secret")
	ErrInvalidDocHeader  = errors.New("invalid doc server header")
	ErrInvalidDemoPeriod = errors.New("demo period has expired")
	ErrHttpNotAllowed    = errors.New("document server must use https protocol for pipedrive integration")
)
