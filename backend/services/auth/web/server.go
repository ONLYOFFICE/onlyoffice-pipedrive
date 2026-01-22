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

package web

import (
	"github.com/ONLYOFFICE/onlyoffice-integration-adapters/service/rpc"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/auth/web/handler"
)

type AuthRPCServer struct {
	selectHandler handler.UserSelectHandler
	insertHandler handler.UserInsertHandler
	deleteHandler handler.UserDeleteHandler
}

func NewAuthRPCServer(
	selectHandler handler.UserSelectHandler,
	insetHandler handler.UserInsertHandler,
	deleteHandler handler.UserDeleteHandler,
) rpc.RPCEngine {
	return AuthRPCServer{
		selectHandler: selectHandler,
		insertHandler: insetHandler,
		deleteHandler: deleteHandler,
	}
}

func (a AuthRPCServer) BuildMessageHandlers() []rpc.RPCMessageHandler {
	return nil
}

func (a AuthRPCServer) BuildHandlers() []interface{} {
	return []interface{}{a.selectHandler, a.insertHandler, a.deleteHandler}
}
