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

import { useQuery } from "react-query";

import { fetchConfig } from "@services/config";

export function useBuildConfig(
  token: string,
  id: string,
  name: string,
  key: string,
  dealID: string,
  dark = false,
) {
  const { isLoading, error, data } = useQuery({
    queryKey: ["config", id, key, dark],
    queryFn: ({ signal }) =>
      fetchConfig(token, id, name, key, dealID, dark, signal),
    staleTime: 0,
    cacheTime: 0,
    refetchOnWindowFocus: false,
  });

  return { isLoading, error, data };
}
