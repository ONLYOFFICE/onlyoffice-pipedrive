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

import { AuthToken } from "@context/TokenContext";

import { FileResponse } from "src/types/file";

export const fetchFiles = async (
  url: string,
  start = 0,
  limit = 50,
  signal: AbortSignal | undefined = undefined,
  sort = "add_time ASC",
) => {
  const client = axios.create();
  axiosRetry(client, {
    retries: 6,
    retryCondition: (error) => error.status !== 200,
    retryDelay: (count) => count * 100,
    shouldResetTimeout: true,
  });

  try {
    const res = await client<FileResponse>({
      method: "GET",
      url,
      params: {
        start,
        limit,
        sort,
      },
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${AuthToken.access_token}`,
      },
      signal,
      timeout: 3500,
    });

    return {
      response: res.data.data,
      pagination: res.data.additional_data.pagination,
    };
  } catch {
    return {
      response: [],
      pagination: {
        start,
        next_start: 0,
        limit,
        more_items_in_collection: false,
      },
    };
  }
};

export const downloadFile = async (domain: string, id: string) => {
  const resp = await axios.get(
    `${process.env.BACKEND_GATEWAY}/files/download?domain=${encodeURIComponent(
      domain,
    )}&file_id=${id}`,
    {
      headers: {
        Authorization: `Bearer ${AuthToken.access_token}`,
      },
    },
  );

  return resp.data;
};

export const deleteFile = async (
  url: string,
  signal: AbortSignal | undefined = undefined,
) => {
  const res = await axios({
    method: "DELETE",
    url,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${AuthToken.access_token}`,
    },
    signal,
    timeout: 4500,
  });

  return res.status === 200;
};

export const uploadFile = async (url: string, deal: string, file: File) => {
  const form = new FormData();
  form.append("file", file);
  form.append("deal_id", deal);

  const res = await axios({
    method: "POST",
    url,
    headers: {
      Accept: "application/json",
      "Content-Type": "multipart/form-data",
      Authorization: `Bearer ${AuthToken.access_token}`,
    },
    data: form,
  });

  return res.data;
};
