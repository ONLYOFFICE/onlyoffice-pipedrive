/**
 *
 * (c) Copyright Ascensio System SIA 2026
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

import axios, { AxiosInstance } from "axios";
import axiosRetry from "axios-retry";
import AppExtensionsSDK, { Command } from "@pipedrive/app-extensions-sdk";

import { SettingsResponse } from "src/types/settings";

const setupRetry = (
  client: AxiosInstance,
  retries: number,
  options?: { retryOn429?: boolean; delayMultiplier?: number },
) => {
  axiosRetry(client as Parameters<typeof axiosRetry>[0], {
    retries,
    retryCondition: options?.retryOn429
      ? (error) => error.status === 429
      : (error) => error.status !== 200,
    ...(options?.delayMultiplier && {
      retryDelay: (count) => count * options.delayMultiplier!,
      shouldResetTimeout: true,
    }),
  });
};

export const postSettings = async (
  sdk: AppExtensionsSDK,
  address: string,
  secret: string,
  header: string,
  demoEnabled = false,
  pluginsEnabled = true,
  autofillEnabled = true,
) => {
  const pctx = await sdk.execute(Command.GET_SIGNED_TOKEN);
  const client = axios.create({ baseURL: process.env.BACKEND_GATEWAY });
  setupRetry(client, 2, { retryOn429: true });

  return client({
    method: "POST",
    url: `/api/settings`,
    headers: {
      "Content-Type": "application/json",
      "X-Pipedrive-App-Context": pctx.token,
    },
    data: {
      doc_address: address,
      doc_secret: secret,
      doc_header: header,
      demo_enabled: demoEnabled,
      plugins_enabled: pluginsEnabled,
      autofill_enabled: autofillEnabled,
    },
    timeout: 4000,
  });
};

export const getSettings = async (sdk: AppExtensionsSDK) => {
  const pctx = await sdk.execute(Command.GET_SIGNED_TOKEN);
  const client = axios.create({ baseURL: process.env.BACKEND_GATEWAY });
  setupRetry(client, 3, { delayMultiplier: 50 });

  const settings = await client<SettingsResponse>({
    method: "GET",
    url: `/api/settings`,
    headers: {
      "Content-Type": "application/json",
      "X-Pipedrive-App-Context": pctx.token,
    },
    timeout: 3000,
  });

  return settings.data;
};

export const checkSettings = async (sdk: AppExtensionsSDK) => {
  const pctx = await sdk.execute(Command.GET_SIGNED_TOKEN);
  const client = axios.create({ baseURL: process.env.BACKEND_GATEWAY });
  setupRetry(client, 2, { delayMultiplier: 50 });

  try {
    const response = await client<{ configured: boolean }>({
      method: "GET",
      url: `/api/settings/check`,
      headers: {
        "Content-Type": "application/json",
        "X-Pipedrive-App-Context": pctx.token,
      },
      timeout: 3000,
    });

    return response.data.configured;
  } catch {
    return false;
  }
};
