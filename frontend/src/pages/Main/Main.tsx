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

import React, { useCallback, useEffect, useRef, useState } from "react";
import AppExtensionsSDK, {
  Command,
  Modal,
} from "@pipedrive/app-extensions-sdk";
import { useTranslation } from "react-i18next";

import { OnlyofficeButton } from "@components/button";
import { OnlyofficeFile } from "@components/file";
import { OnlyofficeFileInfo } from "@components/info";
import { OnlyofficeNoFile } from "@components/nofile";
import { OnlyofficeSpinner } from "@components/spinner";

import { OnlyofficeBackgroundError } from "@layouts/ErrorBackground";

import { useFileSearch } from "@hooks/useFileSearch";

import { checkSettings } from "@services/settings";

import { formatBytes, getFileIcon, isFileSupported } from "@utils/file";
import { getCurrentURL } from "@utils/url";

import SettingsError from "@assets/settings-error.svg";
import { OnlyofficeFileActions } from "./Actions";

export const Main: React.FC = () => {
  const { t } = useTranslation();
  const { url, parameters } = getCurrentURL();
  const [sdk, setSDK] = useState<AppExtensionsSDK | null>();
  const [settingsConfigured, setSettingsConfigured] = useState<boolean | null>(
    null,
  );
  const { isLoading, fetchNextPage, isFetchingNextPage, files, hasNextPage } =
    useFileSearch(
      `${url}api/v1/deals/${parameters.get("selectedIds")}/files`,
      20,
    );

  const observer = useRef<IntersectionObserver | null>(null);
  const lastItem = useCallback(
    (node: Element | null) => {
      if (isLoading) return;
      if (observer.current) observer.current.disconnect();
      observer.current = new IntersectionObserver(async (entries) => {
        if (entries[0].isIntersecting && hasNextPage) {
          fetchNextPage();
        }
      });
      if (node) observer.current.observe(node);
    },
    [isLoading, fetchNextPage, hasNextPage],
  );

  useEffect(() => {
    new AppExtensionsSDK()
      .initialize()
      .then(async (s) => {
        setSDK(s);
        const configured = await checkSettings(s);
        setSettingsConfigured(configured);
      })
      .catch(() => {
        setSDK(null);
        setSettingsConfigured(false);
      });
  }, []);

  if (settingsConfigured === null) {
    return (
      <div className="w-full h-full flex justify-center items-center bg-white dark:bg-dark-bg">
        <OnlyofficeSpinner />
      </div>
    );
  }

  if (settingsConfigured === false) {
    return (
      <OnlyofficeBackgroundError
        Icon={<SettingsError className="mb-5" />}
        title={t("background.settings.title", "Document Server not configured")}
        subtitle={t(
          "background.settings.subtitle",
          "Please configure the Document Server in the settings.",
        )}
        button={t("background.settings.button", "Go to Settings")}
        onClick={async () => {
          window.open(`${url}settings/marketplace`, "_blank");
        }}
      />
    );
  }

  return (
    <div className="table-shadow h-full bg-white dark:bg-dark-bg">
      <div className="overflow-x-hidden overflow-y-auto custom-scroll px-5 h-[85%] md:justify-between mr-5">
        {isLoading && (
          <div className="h-[85%] w-full flex justify-center items-center">
            <OnlyofficeSpinner />
          </div>
        )}
        {!isLoading && (!files || files.length === 0) && (
          <OnlyofficeNoFile
            title={t("files.error.nofiles", "Could not find Pipedrive files")}
          />
        )}
        {!isLoading &&
          files &&
          files.length > 0 &&
          files.map((file, index) => {
            if (files.length === index + 1) {
              return (
                <div key={file.id + file.add_time} ref={lastItem}>
                  <OnlyofficeFile
                    Icon={getFileIcon(file.name)}
                    name={file.name}
                    supported={isFileSupported(file.name)}
                    actions={<OnlyofficeFileActions file={file} />}
                  >
                    <OnlyofficeFileInfo
                      info={{
                        // "Created by": file.person_name,
                        [t("files.info.workspace", "Workspace")]:
                          file.remote_location,
                        [t("files.info.type", "Type")]: file.file_type,
                        [t("files.info.modified", "Date modified")]:
                          file.update_time,
                        [t("files.info.creation", "Creation date")]:
                          file.add_time,
                        [t("files.info.size", "Size")]: formatBytes(
                          file.file_size,
                        ),
                      }}
                    />
                  </OnlyofficeFile>
                </div>
              );
            }
            return (
              <div key={file.id + file.add_time}>
                <OnlyofficeFile
                  Icon={getFileIcon(file.name)}
                  name={file.name}
                  supported={isFileSupported(file.name)}
                  actions={<OnlyofficeFileActions file={file} />}
                >
                  <OnlyofficeFileInfo
                    info={{
                      // "Created by": file.person_name,
                      [t("files.info.workspace", "Workspace")]:
                        file.remote_location,
                      [t("files.info.type", "Type")]: file.file_type,
                      [t("files.info.modified", "Date modified")]:
                        file.update_time,
                      [t("files.info.creation", "Creation date")]:
                        file.add_time,
                      [t("files.info.size", "Size")]: formatBytes(
                        file.file_size,
                      ),
                    }}
                  />
                </OnlyofficeFile>
              </div>
            );
          })}
        {isFetchingNextPage && (
          <div
            className={`relative w-full ${
              isLoading ? "h-full" : "h-fit"
            } my-5 flex justify-center items-center`}
          >
            <OnlyofficeSpinner />
          </div>
        )}
      </div>
      <div className="h-[15%] w-3/4 text-ellipsis flex justify-start items-center px-5">
        <OnlyofficeButton
          text={t("button.upload", "Create or upload document")}
          primary
          onClick={async () => {
            await sdk?.execute(Command.OPEN_MODAL, {
              type: Modal.CUSTOM_MODAL,
              action_id: process.env.PIPEDRIVE_CREATE_MODAL_ID || "",
            });
          }}
        />
      </div>
    </div>
  );
};
