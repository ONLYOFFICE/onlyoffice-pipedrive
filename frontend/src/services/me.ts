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

import axios, { AxiosInstance } from "axios";
import axiosRetry from "axios-retry";
import AppExtensionsSDK, { Command } from "@pipedrive/app-extensions-sdk";

import { AuthToken } from "@context/TokenContext";

import { PipedriveUserResponse, UserResponse } from "src/types/user";

const setupRetry = (client: AxiosInstance, retries: number, delayMultiplier: number) => {
  axiosRetry(client as Parameters<typeof axiosRetry>[0], {
    retries,
    retryCondition: (error) => error.status !== 200,
    retryDelay: (count) => count * delayMultiplier,
    shouldResetTimeout: true,
  });
};

export const getMe = async (sdk: AppExtensionsSDK) => {
  const pctx = await sdk.execute(Command.GET_SIGNED_TOKEN);
  const client = axios.create({ baseURL: process.env.BACKEND_GATEWAY });
  setupRetry(client, 3, 50);

  const res = await client<UserResponse>({
    method: "GET",
    url: `/api/me`,
    headers: {
      "Content-Type": "application/json",
      "X-Pipedrive-App-Context": pctx.token,
    },
    timeout: 5000,
  });

  return { response: res.data };
};

export const getPipedriveMe = async (url: string, accessToken?: string) => {
  const client = axios.create();
  setupRetry(client, 3, 50);

  const res = await client<PipedriveUserResponse>({
    method: "GET",
    url,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${accessToken || AuthToken.access_token}`,
    },
    timeout: 3000,
  });

  return res.data;
};
