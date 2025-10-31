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
import AppExtensionsSDK from "@pipedrive/app-extensions-sdk";
import { useSnapshot } from "valtio";
import { Tabs, TabList, Tab, TabPanel } from "react-tabs";
import { useTranslation } from "react-i18next";

import { OnlyofficeSpinner } from "@components/spinner";

import { AuthToken } from "@context/TokenContext";

import { Creation } from "./Creation";
import { Upload } from "./Upload";

export const CreatePage: React.FC = () => {
  const { t } = useTranslation();
  const { access_token: accessToken, error } = useSnapshot(AuthToken);
  const [selected, setSelected] = useState(0);
  const loadingError = !accessToken && error;
  const loaded = accessToken && !error;

  useEffect(() => {
    new AppExtensionsSDK().initialize({
      size: {
        height: 424,
        width: 622,
      },
    });
  });

  return (
    <div className="relative w-full h-full flex flex-col overflow-hidden bg-white dark:bg-dark-bg">
      {!loaded && !loadingError ? (
        <div className="flex justify-center items-center h-full w-full">
          <OnlyofficeSpinner />
        </div>
      ) : (
        <Tabs
          className="flex justify-center items-start flex-col h-screen"
          onSelect={(index) => setSelected(index)}
        >
          <TabList className="flex justify-start items-center min-h-[40px] w-full bg-gray-100 dark:bg-dark-bg">
            <Tab
              id="create-file"
              className={`flex justify-center items-center text-sm font-inter outline-none hover:cursor-pointer min-h-[40px] ${
                selected === 0
                  ? "text-blue-600 border-b-2 border-blue-600"
                  : "text-slate-500 dark:text-dark-muted"
              }`}
              style={{ margin: "0 0 0 1em", padding: "1em" }}
            >
              {t("button.creation.create", "Create")}
            </Tab>
            <Tab
              id="upload-file"
              className={`flex justify-center items-center text-sm font-inter outline-none hover:cursor-pointer min-h-[40px] ${
                selected === 1
                  ? "text-blue-600 border-b-2 border-blue-600"
                  : "text-slate-500 dark:text-dark-muted"
              }`}
              style={{ margin: "0 0 0 1em", padding: "1em" }}
            >
              {t("button.creation.upload", "Upload")}
            </Tab>
          </TabList>
          <TabPanel
            className={`${
              selected === 0 ? "h-[calc(100%-40px)] w-full" : "h-0 w-0"
            }`}
          >
            <Creation />
          </TabPanel>
          <TabPanel
            className={`${
              selected === 1 ? "h-[calc(100%-40px)] w-full" : "h-0 w-0"
            }`}
          >
            <Upload />
          </TabPanel>
        </Tabs>
      )}
    </div>
  );
};
