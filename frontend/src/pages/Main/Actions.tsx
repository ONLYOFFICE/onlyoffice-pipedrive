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

import More from "@assets/more.svg";
import MoreDark from "@assets/more_dark.svg";

type FileActionsProps = {
  file: File;
  onRenameClick: () => void;
  isRenaming?: boolean;
};

export const OnlyofficeFileActions: React.FC<FileActionsProps> = ({
  file,
  onRenameClick,
  isRenaming = false,
}) => {
  const { t } = useTranslation();
  const { url, parameters } = getCurrentURL();
  const { isDark } = useTheme();
  const [sdk, setSDK] = useState<AppExtensionsSDK | null>();
  const [disable, setDisable] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [openDirectionReverse, setOpenDirectionReverse] = useState(false);

  const dropdownRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const deleteMutator = useDeleteFile(`${url}api/v1/files/${file.id}`);

  useEffect(() => {
    new AppExtensionsSDK()
      .initialize()
      .then((s) => setSDK(s))
      .catch(() => setSDK(null));
  }, []);

  const handleClickOutside = (event: MouseEvent) => {
    if (
      dropdownRef.current &&
      !dropdownRef.current.contains(event.target as Node)
    ) {
      setIsDropdownOpen(false);
    }
  };

  const checkPosition = () => {
    if (buttonRef.current) {
      const rect = buttonRef.current.getBoundingClientRect();
      const dropdownHeight = 128;
      const spaceBelow = window.innerHeight - rect.bottom;
      const shouldReverse = spaceBelow < dropdownHeight + 20;
      setOpenDirectionReverse(shouldReverse);
    }
  };

  useEffect(() => {
    if (isDropdownOpen) {
      checkPosition();
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isDropdownOpen]);

  const handleMenuItemClick =
    (handler: () => void) => (e: React.MouseEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDropdownOpen(false);
      if (!disable) {
        handler();
      }
    };

  const handleDelete = () => {
    setDisable(true);
    deleteMutator
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

  const handleRename = () => {
    if (!disable) {
      onRenameClick();
    }
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

  const isEditorDisabled = !isFileSupported(file.name) || disable || isRenaming;
  const isDownloadDisabled = disable || isRenaming;
  const isRenameDisabled = disable || isRenaming;
  const isDeleteDisabled = disable || isRenaming;

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        ref={buttonRef}
        type="button"
        className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          setIsDropdownOpen(!isDropdownOpen);
        }}
        aria-label={t("files.actions.menu", "Actions menu")}
        aria-expanded={isDropdownOpen}
      >
        {isDark ? <MoreDark /> : <More />}
      </button>

      {isDropdownOpen && (
        <div
          className={`absolute right-0 bg-white dark:bg-dark-bg border border-gray-200 dark:border-dark-border rounded-lg shadow-lg z-50 ${
            openDirectionReverse ? "bottom-full mb-1" : "top-full mt-1"
          }`}
          style={{ width: "8.375rem", maxWidth: "12.5rem" }}
        >
          <div>
            <button
              type="button"
              className={`w-full px-4 text-sm text-left ${
                isEditorDisabled
                  ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                  : "text-gray-700 dark:text-dark-text hover:bg-gray-100 dark:hover:bg-gray-700"
              }`}
              style={{ height: "2rem" }}
              onClick={
                isEditorDisabled ? undefined : handleMenuItemClick(handleEditor)
              }
              disabled={isEditorDisabled}
            >
              {t("files.actions.edit", "Open")}
            </button>

            <button
              type="button"
              className={`w-full px-4 text-sm text-left ${
                isRenameDisabled
                  ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                  : "text-gray-700 dark:text-dark-text hover:bg-gray-100 dark:hover:bg-gray-700"
              }`}
              style={{ height: "2rem" }}
              onClick={handleMenuItemClick(handleRename)}
              disabled={isRenameDisabled}
            >
              {t("files.actions.rename", "Rename")}
            </button>

            <button
              type="button"
              className={`w-full px-4 text-sm text-left ${
                isDownloadDisabled
                  ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                  : "text-gray-700 dark:text-dark-text hover:bg-gray-100 dark:hover:bg-gray-700"
              }`}
              style={{ height: "2rem" }}
              onClick={handleMenuItemClick(handleDownload)}
              disabled={isDownloadDisabled}
            >
              {t("files.actions.download", "Download")}
            </button>

            <button
              type="button"
              className={`w-full px-4 text-sm text-left ${
                isDeleteDisabled
                  ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                  : "text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700"
              }`}
              style={{ height: "2rem" }}
              onClick={handleMenuItemClick(handleDelete)}
              disabled={isDeleteDisabled}
            >
              {t("files.actions.delete", "Delete")}
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
