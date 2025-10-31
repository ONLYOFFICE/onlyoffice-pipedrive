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

type Permissions = {
  comment: boolean;
  copy: boolean;
  deleteCommentAuthorOnly: boolean;
  download: boolean;
  edit: boolean;
  editCommentAuthorOnly: boolean;
  fillForms: boolean;
  modifyContentControl: boolean;
  modifyFilter: boolean;
  print: boolean;
  review: boolean;
};

type Document = {
  fileType: string;
  key: string;
  title: string;
  url: string;
  permissions: Permissions;
};

type User = {
  id: string;
  name: string;
};

type Goback = {
  requestClost: boolean;
};

type Customization = {
  goback: Goback;
  hideRightMenu: boolean;
  plugins: boolean;
};

type EditorConfig = {
  user: User;
  callbackUrl: string;
  customization: Customization;
  lang: string;
  plugins: Plugins;
};

type Plugins = {
  autostart: string[];
  options: {
    all?: object;
    pluginGuid: object;
  };
  pluginsData: string[];
  url: string;
};

export type ConfigResponse = {
  document: Document;
  documentType: string;
  editorConfig: EditorConfig;
  type: string;
  token: string;
  server_url: string;
  demo_enabled: boolean;
};
