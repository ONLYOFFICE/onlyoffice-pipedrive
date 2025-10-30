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

import React, { useEffect, useState } from "react";
import i18next from "i18next";
import md5 from "md5";
import AppExtensionsSDK, { Command } from "@pipedrive/app-extensions-sdk";
import { useTranslation } from "react-i18next";

import { useTheme } from "@context/ThemeContext";
import { useDeleteFile } from "@hooks/useDeleteFile";

import { downloadFile } from "@services/file";

import { getFileParts, isFileSupported } from "@utils/file";
import { getCurrentURL } from "@utils/url";

import { File } from "src/types/file";

import Pencil from "@assets/pencil.svg";
import PencilDark from "@assets/pencil_dark.svg";
import Download from "@assets/download.svg";
import DownloadDark from "@assets/download_dark.svg";
import Trash from "@assets/trash.svg";
import TrashDark from "@assets/trash_dark.svg";

type FileActionsProps = {
  file: File;
};

export const OnlyofficeFileActions: React.FC<FileActionsProps> = ({ file }) => {
  const { t } = useTranslation();
  const { url, parameters } = getCurrentURL();
  const { isDark } = useTheme();
  const [sdk, setSDK] = useState<AppExtensionsSDK | null>();
  const [disable, setDisable] = useState(false);
  const mutator = useDeleteFile(`${url}api/v1/files/${file.id}`);

  useEffect(() => {
    new AppExtensionsSDK()
      .initialize()
      .then((s) => setSDK(s))
      .catch(() => setSDK(null));
  }, []);

  const handleDelete = () => {
    setDisable(true);
    mutator
      .mutateAsync()
      .then(async () => {
        await sdk?.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "snackbar.fileremoved.ok",
            `File ${file.name} has been removed`,
            { file: file.name },
          ),
        });
      })
      .catch(async () => {
        setDisable(false);
        await sdk?.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "snackbar.fileremoved.error",
            `Could not remove file ${file.name}`,
            { file: file.name },
          ),
        });
      });
  };

  const handleEditor = async () => {
    setDisable(true);
    if (isFileSupported(file.name)) {
      const win = window.open("/editor");
      const token = await sdk?.execute(Command.GET_SIGNED_TOKEN);
      if (token) {
        const [name, ext] = getFileParts(file.name);
        if (win && win.location)
          win.location.href = `/editor?token=${token.token}&deal_id=${
            parameters.get("selectedIds") || "1"
          }&id=${file.id}&name=${`${encodeURIComponent(
            name.substring(0, 190),
          )}.${ext}`}&key=${md5(file.id + file.update_time)}&lng=${
            i18next.language
          }&dark=${isDark}`;
      }
    }
    // temporary solution
    setTimeout(() => setDisable(false), 10000);
  };

  const handleDownload = async () => {
    setDisable(true);
    try {
      const durl = await downloadFile(url, file.id);
      window.open(durl);
    } catch {
      await sdk?.execute(Command.SHOW_SNACKBAR, {
        message: t(
          "snackbar.filedownload.error",
          `Could not download file ${file.name}`,
          { file: file.name },
        ),
      });
    } finally {
      setDisable(false);
    }
  };

  const handleClick = (handler: () => void) => (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (!disable) {
      handler();
    }
  };

  const handleKeyDown = (handler: () => void) => (e: React.KeyboardEvent) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      e.stopPropagation();
      if (!disable) {
        handler();
      }
    }
  };

  const isEditorDisabled = !isFileSupported(file.name) || disable;
  const isDownloadDisabled = disable;
  const isDeleteDisabled = disable;

  return (
    <>
      <div
        role="button"
        tabIndex={isEditorDisabled ? -1 : 0}
        className={`${
          isEditorDisabled
            ? "hover:cursor-default opacity-50"
            : "hover:cursor-pointer"
        } mx-1`}
        onClick={isEditorDisabled ? undefined : handleClick(handleEditor)}
        onKeyDown={isEditorDisabled ? undefined : handleKeyDown(handleEditor)}
        aria-disabled={isEditorDisabled}
      >
        {isDark ? <PencilDark /> : <Pencil />}
      </div>
      <div
        role="button"
        tabIndex={isDownloadDisabled ? -1 : 0}
        className={`mx-1 ${
          isDownloadDisabled
            ? "hover:cursor-default opacity-50"
            : "hover:cursor-pointer"
        }`}
        onClick={handleClick(handleDownload)}
        onKeyDown={handleKeyDown(handleDownload)}
        aria-disabled={isDownloadDisabled}
      >
        {isDark ? <DownloadDark /> : <Download />}
      </div>
      <div
        role="button"
        tabIndex={isDeleteDisabled ? -1 : 0}
        className={`mx-1 ${
          isDeleteDisabled
            ? "hover:cursor-default opacity-50"
            : "hover:cursor-pointer"
        }`}
        onClick={handleClick(handleDelete)}
        onKeyDown={handleKeyDown(handleDelete)}
        aria-disabled={isDeleteDisabled}
      >
        {isDark ? <TrashDark /> : <Trash />}
      </div>
    </>
  );
};
