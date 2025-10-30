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

import AppExtensionsSDK from "@pipedrive/app-extensions-sdk";
import i18next from "i18next";
import axios, { AxiosError } from "axios";
import React, { useEffect, ReactNode } from "react";
import { proxy } from "valtio";

import { getMe, getPipedriveMe } from "@services/me";

import { getCurrentURL } from "@utils/url";

export const AuthToken = proxy({
  access_token: "",
  expires_at: Date.now(),
  error: false,
  status: 200,
});

type ProviderProps = {
  children?: ReactNode;
};

const TokenContext = React.createContext<boolean>(true);

export const TokenProvider: React.FC<ProviderProps> = ({ children }) => {
  useEffect(() => {
    let timerID: NodeJS.Timeout;
    new AppExtensionsSDK()
      .initialize()
      .then((sdk) => {
        const { url } = getCurrentURL();
        timerID = setTimeout(async function update() {
          if (
            !AuthToken.error &&
            (!AuthToken.access_token ||
              AuthToken.expires_at <= Date.now() - 1000 * 30)
          ) {
            try {
              const token = await getMe(sdk);
              const resp = await getPipedriveMe(
                `${url}api/v1/users/me`,
                token.response.access_token,
              );
              await i18next.changeLanguage(
                `${resp.data.language.language_code}-${resp.data.language.country_code}`,
              );
              AuthToken.access_token = token.response.access_token;
              AuthToken.expires_at = token.response.expires_at;
            } catch (err) {
              if (axios.isAxiosError(err) && !axios.isCancel(err)) {
                const aerr = err as AxiosError;
                if (aerr.response && aerr.response.status) {
                  AuthToken.status = aerr.response.status;
                } else {
                  AuthToken.status = 500;
                }
                AuthToken.error = true;
                AuthToken.access_token = "";
              }
            }
          }
          timerID = setTimeout(
            update,
            AuthToken.expires_at - Date.now() - 1000 * 30,
          );
        }, 0);
      })
      .catch(() => null);

    return () => clearTimeout(timerID);
  }, []);
  return <TokenContext.Provider value>{children}</TokenContext.Provider>;
};
