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

import React from "react";
import { useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { DocumentEditor } from "@onlyoffice/document-editor-react";
import { Helmet } from "react-helmet";

import { OnlyofficeButton } from "@components/button";
import { OnlyofficeError } from "@components/error";
import { OnlyofficeSpinner } from "@components/spinner";

import { useBuildConfig } from "@hooks/useBuildConfig";

import { getFileFavicon } from "@utils/file";

import Icon from "@assets/nofile.svg";

declare global {
  interface Window {
    connector: any;
  }
}

const onEditor = () => {
  const loader = document.getElementById("eloader");
  if (loader) {
    loader.classList.add("opacity-0");
    loader.classList.add("-z-100");
    loader.classList.add("hidden");
  }

  const editor = document.getElementById("editor");
  if (editor) {
    editor.classList.remove("opacity-0");
  }
};

export const OnlyofficeEditorPage: React.FC = () => {
  const { t } = useTranslation();
  const [params] = useSearchParams();

  const isDark = params.get("dark") === "true";
  const { isLoading, error, data } = useBuildConfig(
    params.get("token") || "",
    params.get("id") || "",
    params.get("name") || "new.docx",
    params.get("key") || new Date().toTimeString(),
    params.get("deal_id") || "1",
    isDark
  );

  const validConfig = !error && !isLoading && data;
  const backgroundClass = isDark ? "bg-dark-bg" : "bg-white";
  return (
    <div className={`w-full h-full overflow-hidden ${backgroundClass}`}>
      <Helmet
        link={[
          {
            rel: "shortcut icon",
            type: "image/x-icon",
            href: `${getFileFavicon(params.get("name") || "new.docx")}`,
          },
        ]}
      />
      {!error && (
        <div
          id="eloader"
          className={`relative w-full h-full flex flex-col small:justify-between justify-center items-center transition duration-250 ease-linear ${backgroundClass}`}
        >
          <div className="pb-5 small:h-full small:flex small:items-center">
            <OnlyofficeSpinner isDark={isDark} />
          </div>
          <div className="small:mb-5 small:px-5 small:w-full">
            <OnlyofficeButton
              primary
              text={t("button.cancel", "Cancel")}
              fullWidth
              onClick={() => window.close()}
            />
          </div>
        </div>
      )}
      {!!error && (
        <div
          className={`w-full h-full flex justify-center flex-col items-center mb-1 ${backgroundClass}`}
        >
          <Icon />
          <OnlyofficeError
            text={t(
              "editor.error",
              "Could not open the file. Something went wrong"
            )}
          />
          <div className="pt-5">
            <OnlyofficeButton
              primary
              text={t("button.close", "Close")}
              onClick={() => window.close()}
            />
          </div>
        </div>
      )}
      {validConfig && data.server_url && (
        <div
          id="editor"
          className="w-full h-full opacity-0 transition duration-250 ease-linear"
        >
          <DocumentEditor
            id="docxEditor"
            documentServerUrl={data.server_url}
            config={{
              document: {
                fileType: data.document.fileType,
                key: data.document.key,
                title: data.document.title,
                url: data.document.url,
                permissions: data.document.permissions,
              },
              documentType: data.documentType,
              editorConfig: {
                callbackUrl: data.editorConfig.callbackUrl,
                user: data.editorConfig.user,
                lang: data.editorConfig.lang,
                customization: {
                  hideRightMenu: data.editorConfig.customization.hideRightMenu,
                  plugins: data.editorConfig.customization.plugins,
                },
                plugins: data.editorConfig.plugins,
              },
              token: data.token,
              type: data.type,
              events: {
                onRequestClose: async () => {
                  window.close();
                },
                onAppReady: onEditor,
                onError: () => {
                  onEditor();
                },
                onWarning: onEditor,
                onDocumentReady: (event: any) => {
                  const connector = window.DocEditor.instances["docxEditor"].createConnector();
                  window.connector = connector;
                },
              },
            }}
          />
        </div>
      )}
    </div>
  );
};

export default OnlyofficeEditorPage;
