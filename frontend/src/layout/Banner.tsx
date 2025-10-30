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

import { OnlyofficeTitle, OnlyofficeSubtitle } from "@components/title";

import BannerIcon from "@assets/banner.svg";
import { useTranslation } from "react-i18next";

export const Banner: React.FC = () => {
  const { t } = useTranslation();
  return (
    <div className="flex justify-between items-center p-5 mt-5 mb-5 bg-gray-100 dark:bg-dark-bg border border-gray-300 dark:border-dark-border">
      <div className="w-2/12">
        <BannerIcon />
      </div>
      <div className="w-7/12 flex justify-center items-start">
        <div className="flex justify-start items-center flex-col cursor-default ml-5 mr-5">
          <div className="w-full h-1/2 flex">
            <OnlyofficeTitle
              text={
                t("banner.title", "ONLYOFFICE Docs Cloud") ||
                "ONLYOFFICE Docs Cloud"
              }
              large
            />
          </div>
          <div className="w-full h-[40px] overflow-hidden">
            <OnlyofficeSubtitle
              text={
                t(
                  "banner.subtitle",
                  "Easily launch the editors in the cloud without downloading and installation",
                ) ||
                "Easily launch the editors in the cloud without downloading and installation"
              }
              center={false}
            />
          </div>
        </div>
      </div>
      <div className="w-3/12">
        <button
          type="button"
          className="pl-5 pr-5 pt-2 pb-2 text-sm rounded overflow-hidden text-ellipsis inline-block max-w-[120px] cursor-pointer hover:shadow-sm duration-200 bg-gray-800 dark:bg-gray-700 text-white"
          onClick={() =>
            window.open(
              "https://www.onlyoffice.com/docs-registration.aspx?referer=pipedrive",
            )
          }
        >
          {t("button.getnow", "Get Now")}
        </button>
      </div>
    </div>
  );
};
