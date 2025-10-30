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

export function setWithExpiry<T>(key: string, value: T, expiry: number): void {
  localStorage.setItem(
    key,
    JSON.stringify({
      value,
      expiry,
    }),
  );
}

export function getWithExpiry<T>(key: string): T | null {
  const sitem = localStorage.getItem(key);
  if (!sitem) {
    return null;
  }

  const item = JSON.parse(sitem);
  const now = new Date();
  if (now.getTime() > item.expiry) {
    localStorage.removeItem(key);
    return null;
  }

  return item.value as T;
}
