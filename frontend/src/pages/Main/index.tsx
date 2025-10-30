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
import { useSnapshot } from "valtio";
import { useTranslation } from "react-i18next";

import { OnlyofficeSpinner } from "@components/spinner";

import { OnlyofficeBackgroundError } from "@layouts/ErrorBackground";

import { AuthToken } from "@context/TokenContext";

import { getCurrentURL } from "@utils/url";

import TokenError from "@assets/token-error.svg";

import { Main } from "./Main";

export const MainPage: React.FC = () => {
  const { t } = useTranslation();
  const { access_token: accessToken, error, status } = useSnapshot(AuthToken);
  const loading = !accessToken && !error;
  const loadingError = !accessToken && error;
  const loaded = accessToken && !error;
  return (
    <div className="relative w-full h-full flex flex-col my-0 mx-auto bg-white dark:bg-dark-bg">
      {loading && (
        <div className="w-full h-full flex justify-center items-center">
          <OnlyofficeSpinner />
        </div>
      )}
      {loadingError && (
        <OnlyofficeBackgroundError
          Icon={<TokenError className="mb-5" />}
          title={t(
            status === 401
              ? "background.reinstall.title"
              : "background.error.title.main",
            status === 401
              ? "The document security token has expired"
              : "Something went wrong",
          )}
          subtitle={t(
            status === 401
              ? "background.reinstall.subtitle.token"
              : "background.reinstall.subtitle",
            status === 401
              ? "Something went wrong. Please reinstall the app."
              : "Something went wrong. Please reload the app.",
          )}
          button={t("background.reinstall.button", "Reinstall") || "Reinstall"}
          onClick={
            status === 401
              ? () =>
                  window.open(
                    `${getCurrentURL().url}settings/marketplace`,
                    "_blank",
                  )
              : undefined
          }
        />
      )}
      {loaded && <Main />}
    </div>
  );
};
