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

import axios from "axios";
import axiosRetry from "axios-retry";

import { ConfigResponse } from "src/types/config";

export const fetchConfig = async (
  token: string,
  id: string,
  name: string,
  key: string,
  dealID: string,
  dark?: boolean,
  signal?: AbortSignal,
) => {
  const client = axios.create();
  axiosRetry(client, {
    retries: 2,
    retryCondition: (error) => error.status !== 200,
    retryDelay: (count) => count * 50,
    shouldResetTimeout: true,
  });

  const res = await axios<ConfigResponse>({
    method: "GET",
    url: `${process.env.BACKEND_GATEWAY}/api/config`,
    params: {
      id,
      name,
      key,
      deal_id: dealID,
      dark: dark?.toString() || "false",
    },
    headers: {
      "Content-Type": "application/json",
      "X-Pipedrive-App-Context": token,
    },
    signal,
  });
  return res.data;
};
