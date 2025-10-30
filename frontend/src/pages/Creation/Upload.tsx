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

import { FileRejection, DropEvent } from "react-dropzone";
import React, { useEffect, useState } from "react";
import AppExtensionsSDK, { Command } from "@pipedrive/app-extensions-sdk";
import { useTranslation } from "react-i18next";

import { OnlyofficeButton } from "@components/button";
import { OnlyofficeDragDrop } from "@components/drop";

import { uploadFile } from "@services/file";

import { getCurrentURL } from "@utils/url";

const onDrop = <T extends File>(acceptedFiles: T[]): Promise<void> => {
  const { url, parameters } = getCurrentURL();
  return uploadFile(
    `${url}api/v1/files`,
    parameters.get("selectedIds") || "",
    acceptedFiles[0],
  );
};

export const Upload: React.FC = () => {
  const { t } = useTranslation();
  const [sdk, setSDK] = useState<AppExtensionsSDK | null>();
  useEffect(() => {
    new AppExtensionsSDK()
      .initialize()
      .then((s) => setSDK(s))
      .catch(() => setSDK(null));
  }, []);

  return (
    <div className="h-full bg-white dark:bg-dark-bg">
      <div className="h-[calc(100%-3rem)] overflow-hidden">
        <div className="px-5 py-20 flex flex-col justify-center items-start h-full">
          <OnlyofficeDragDrop
            errorText={
              t(
                "upload.error",
                "Could not upload your file. Please contact ONLYOFFICE support.",
              ) ||
              "Could not upload your file. Please contact ONLYOFFICE support."
            }
            uploadingText={
              t("upload.uploading", "Uploading...") || "Uploading..."
            }
            selectText={t("upload.select", "Select a file") || "Select a file"}
            dragdropText={
              t("upload.dragdrop", "or drag and drop here") ||
              "or drag and drop here"
            }
            subtext={
              t("upload.subtext", "File size is limited") ||
              "File size is limited"
            }
            onDrop={async (
              files,
              // eslint-disable-next-line @typescript-eslint/no-unused-vars
              _rejections: FileRejection[],
              // eslint-disable-next-line @typescript-eslint/no-unused-vars
              _event: DropEvent,
            ) => {
              try {
                await new Promise((resolve) => {
                  setTimeout(resolve, 1000);
                });
                await onDrop(files);
                await sdk?.execute(Command.SHOW_SNACKBAR, {
                  message: t(
                    "snackbar.uploaded.ok",
                    "File {{file}} has been uploaded",
                    { file: files[0].name },
                  ),
                });
                return Promise.resolve();
              } catch {
                await sdk?.execute(Command.SHOW_SNACKBAR, {
                  message: t(
                    "snackbar.uploaded.error",
                    "Could not upload file {{file}}",
                    { file: files[0].name },
                  ),
                });
                return Promise.reject();
              }
            }}
          />
        </div>
      </div>
      <div className="h-[48px] flex items-center w-full bg-white dark:bg-dark-bg border-t dark:border-dark-border">
        <div className="flex justify-between items-center w-full">
          <div className="mx-5">
            <OnlyofficeButton
              text={t("button.cancel", "Cancel")}
              onClick={async () => {
                await sdk?.execute(Command.CLOSE_MODAL);
              }}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
