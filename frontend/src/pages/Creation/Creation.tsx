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

import React, { useEffect, useRef, useState } from "react";
import md5 from "md5";
import AppExtensionsSDK, { Command } from "@pipedrive/app-extensions-sdk";
import { useTranslation } from "react-i18next";
import i18next from "i18next";

import { OnlyofficeButton } from "@components/button";
import { OnlyofficeInput } from "@components/input";
import { OnlyofficeTile } from "@components/tile";
import { OnlyofficeTitle } from "@components/title";

import { createFile } from "@services/file";

import { getFileIcon } from "@utils/file";
import { getCurrentURL } from "@utils/url";

import Redirect from "@assets/redirect.svg";

const HelpText: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div className="text-xs text-gray-600 dark:text-gray-400 mt-4 text-left leading-relaxed">
      {t("creation.help.text", "Help us improve the ONLYOFFICE app:")}{" "}
      <a
        href="https://feedback.onlyoffice.com/forums/966080-your-voice-matters?category_id=519288"
        target="_blank"
        rel="noopener noreferrer"
        className="text-gray-600 dark:text-blue-500 underline hover:text-gray-800 dark:hover:text-blue-400"
      >
        {t("creation.help.feedback", "share your feedback")}
      </a>
      . {t("creation.help.visit", "Visit our")}{" "}
      <a
        href="https://helpcenter.onlyoffice.com/integration/pipedrive.aspx"
        target="_blank"
        rel="noopener noreferrer"
        className="text-gray-600 dark:text-blue-500 underline hover:text-gray-800 dark:hover:text-blue-400"
      >
        {t("creation.help.helpcenter", "Help Center")}
      </a>{" "}
      {t("creation.help.details", "for more details")}.
    </div>
  );
};

export const Creation: React.FC = () => {
  const { t } = useTranslation();
  const [creating, setCreating] = useState(false);
  const [sdk, setSDK] = useState<AppExtensionsSDK | null>();
  const [file, setFile] = useState(
    t("document.new", "New Document") || "New Document",
  );
  const [fileType, setFileType] = useState<"docx" | "pptx" | "xlsx">("docx");
  const handleChangeFile = (newType: "docx" | "pptx" | "xlsx") => {
    if (!creating) setFileType(newType);
  };

  const isFileNameValid = file.trim().length > 0 && file.length <= 190;
  const fileTypeRef = useRef(fileType);

  useEffect(() => {
    new AppExtensionsSDK()
      .initialize()
      .then((s) => setSDK(s))
      .catch(() => setSDK(null));
  }, []);

  useEffect(() => {
    if (fileTypeRef.current !== fileType) {
      const defaultDocx = t("document.new", "New Document") || "New Document";
      const defaultPptx =
        t("document.new.presentation", "New Presentation") ||
        "New Presentation";
      const defaultXlsx =
        t("document.new.spreadsheet", "New Spreadsheet") || "New Spreadsheet";

      const isDefault =
        file === defaultDocx || file === defaultPptx || file === defaultXlsx;
      if (isDefault) {
        switch (fileType) {
          case "docx":
            setFile(defaultDocx);
            break;
          case "pptx":
            setFile(defaultPptx);
            break;
          case "xlsx":
            setFile(defaultXlsx);
            break;
          default:
            break;
        }
      }

      fileTypeRef.current = fileType;
    }
  }, [fileType, file, t]);

  const handleCreateDocument = async () => {
    if (!isFileNameValid || creating) return;

    setCreating(true);
    const token = await sdk?.execute(Command.GET_SIGNED_TOKEN);
    if (!token) return;
    const { parameters } = getCurrentURL();

    try {
      const filename = `${file
        .replaceAll("/", ":")
        .replaceAll("\\", ":")
        .substring(0, 190)}.${fileType}`;

      const fres = await createFile(
        token.token,
        i18next.language,
        fileType,
        parameters.get("selectedIds") || "",
        filename,
      );

      window.open(
        `/editor?token=${token.token}&id=${fres.data.id}&deal_id=${
          fres.data.deal_id
        }&name=${`${encodeURIComponent(
          file.substring(0, 190),
        )}.${fileType}`}&key=${md5(
          fres.data.id + fres.data.update_time,
        )}&lng=${i18next.language}`,
      );
      await sdk?.execute(Command.CLOSE_MODAL);
    } catch {
      await sdk?.execute(Command.SHOW_SNACKBAR, {
        message: t("creation.error", "Could not create a new file"),
      });
    } finally {
      setCreating(false);
    }
  };

  return (
    <div className="h-full w-full bg-white dark:bg-dark-bg flex flex-col">
      <div className="flex-1 w-full overflow-hidden">
        <div className="px-5 py-5 flex flex-col h-full">
          <OnlyofficeTitle
            text={t("creation.title", "Create with ONLYOFFICE")}
            large
            align="left"
          />
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
            {t(
              "creation.subtitle",
              "The file will open in a new ONLYOFFICE tab and will be automatically saved to this Pipedrive deal.",
            )}
          </p>
          <div className="w-full pt-6">
            <OnlyofficeInput
              text={t("creation.inputs.filename", "File name")}
              labelSize="sm"
              valid={isFileNameValid}
              errorText={
                file.trim().length === 0
                  ? t(
                      "creation.inputs.error.empty",
                      "File name cannot be empty",
                    ) || "File name cannot be empty"
                  : t(
                      "creation.inputs.error",
                      "File name length should be less than 190 characters",
                    ) || "File name length should be less than 190 characters"
              }
              value={file}
              onChange={(e) => setFile(e.target.value)}
              disabled={creating}
            />
          </div>
          <div className="w-full flex gap-4 pt-4">
            <div className="flex-1">
              <OnlyofficeTile
                Icon={getFileIcon("sample.docx")}
                text={t("creation.tiles.doc", "Document")}
                onClick={() => handleChangeFile("docx")}
                onKeyDown={() => handleChangeFile("docx")}
                selected={fileType === "docx"}
              />
            </div>
            <div className="flex-1">
              <OnlyofficeTile
                Icon={getFileIcon("sample.xlsx")}
                text={t("creation.tiles.spreadsheet", "Spreadsheet")}
                onClick={() => handleChangeFile("xlsx")}
                onKeyDown={() => handleChangeFile("xlsx")}
                selected={fileType === "xlsx"}
              />
            </div>
            <div className="flex-1">
              <OnlyofficeTile
                Icon={getFileIcon("sample.pptx")}
                text={t("creation.tiles.presentation", "Presentation")}
                onClick={() => handleChangeFile("pptx")}
                onKeyDown={() => handleChangeFile("pptx")}
                selected={fileType === "pptx"}
              />
            </div>
          </div>
          <HelpText />
        </div>
      </div>
      <div className="h-[48px] flex items-center w-full bg-white dark:bg-dark-bg border-t dark:border-dark-border">
        <div className="flex justify-between items-center w-full">
          <div className="mx-5">
            <OnlyofficeButton
              text={t("button.close", "Close")}
              onClick={async () => {
                await sdk?.execute(Command.CLOSE_MODAL);
              }}
              disabled={creating}
            />
          </div>
          <div className="mx-5">
            <OnlyofficeButton
              disabled={!isFileNameValid || creating}
              text={t("button.create", "Create")}
              primary
              Icon={<Redirect />}
              onClick={handleCreateDocument}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
