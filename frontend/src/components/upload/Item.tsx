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

import type { UploadedFile } from "@objects/file";
import { Status } from "./Status";

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / k ** i).toFixed(1))} ${sizes[i]}`;
};

export const UploadItem: React.FC<{
  file: UploadedFile;
  onRemove: (fileId: string) => void;
  onCancel?: (fileId: string) => void;
}> = ({ file, onRemove, onCancel }) => {
  const { id, name, size, status } = file;
  const isUploading = status === "uploading";

  return (
    <div className="flex items-center justify-between py-2 border-b border-gray-200 dark:border-dark-border last:border-b-0">
      <div className="flex items-center flex-1 min-w-0">
        <Status status={status} />
        <span className="ml-[8px] text-sm text-gray-900 dark:text-dark-text truncate">
          {name}
        </span>
      </div>
      <div className="flex items-center gap-3 ml-2 flex-shrink-0">
        <span className="text-xs text-gray-500 dark:text-gray-400">
          {formatFileSize(size)}
        </span>
        <button
          type="button"
          onClick={() =>
            isUploading && onCancel ? onCancel(id) : onRemove(id)
          }
          className="group transition-all"
          aria-label={isUploading ? "Cancel upload" : "Remove file"}
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 16 16"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <circle
              cx="8"
              cy="8"
              r="8"
              fill="#5A636F"
              className="group-hover:fill-[#333333] transition-colors"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M8.99049 8.00018L11.9951 11.0048L11.0052 11.9946L8.00069 8.99014L4.99741 11.9939L4.00745 11.0041L7.01081 8.00026L4.00567 4.99512L4.99555 4.00524L8.00061 7.0103L11.0051 4.00528L11.9951 4.99508L8.99049 8.00018Z"
              fill="white"
            />
          </svg>
        </button>
      </div>
    </div>
  );
};
