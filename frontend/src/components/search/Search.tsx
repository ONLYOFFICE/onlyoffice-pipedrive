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

import React from "react";

type SearchProps = {
  value?: string;
  placeholder?: string;
  disabled?: boolean;
  autocomplete?: boolean;
  onChange?: React.ChangeEventHandler<HTMLInputElement>;
};

export const OnlyofficeSearchBar: React.FC<SearchProps> = ({
  value,
  placeholder,
  disabled,
  autocomplete = false,
  onChange,
}) => (
  <div className="font-sans text-black dark:text-dark-text bg-white dark:bg-dark-bg w-screen">
    <div className="border dark:border-dark-border rounded overflow-hidden flex">
      <input
        type="text"
        className="py-2 px-2 w-full select-auto outline-none bg-white dark:bg-dark-bg text-black dark:text-dark-text"
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        disabled={disabled}
        autoCorrect={autocomplete ? undefined : "off"}
        autoComplete={autocomplete ? undefined : "off"}
      />
      <button
        type="button"
        className={`px-6 ${disabled && "bg-gray-50 dark:bg-dark-bg"}`}
        disabled={disabled}
        aria-label="Search"
      >
        <svg
          className="h-4 w-4 text-grey-dark dark:text-dark-muted"
          fill="currentColor"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
        >
          <path d="M16.32 14.9l5.39 5.4a1 1 0 0 1-1.42 1.4l-5.38-5.38a8 8 0 1 1 1.41-1.41zM10 16a6 6 0 1 0 0-12 6 6 0 0 0 0 12z" />
        </svg>
      </button>
    </div>
  </div>
);
