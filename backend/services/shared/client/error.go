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

package client

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUrlFormat     = errors.New("url is not valid")
	ErrInvalidContentLength = errors.New("could not perform api actions due to exceeding content-length")
)

type UnexpectedStatusCodeError struct {
	Action string
	Code   int
}

func (e *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("could not perform pipedrive %s action. Status code: %d", e.Action, e.Code)
}
