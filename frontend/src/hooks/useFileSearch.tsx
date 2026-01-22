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

import { useInfiniteQuery } from "react-query";

import { fetchFiles } from "@services/file";

export function useFileSearch(url: string, limit: number) {
  const {
    data,
    isLoading,
    error,
    refetch,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery({
    queryKey: ["filesData", url],
    queryFn: ({ signal, pageParam }) =>
      fetchFiles(url, pageParam, limit, signal),
    getNextPageParam: (lastPage) =>
      lastPage?.pagination?.more_items_in_collection
        ? lastPage.pagination.next_start
        : undefined,
    staleTime: 3500,
    cacheTime: 4000,
    refetchInterval: 3500,
  });

  return {
    files:
      data?.pages
        .map((page) => page.response)
        .filter(Boolean)
        .flat() || [],
    isLoading,
    error,
    refetch,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  };
}
