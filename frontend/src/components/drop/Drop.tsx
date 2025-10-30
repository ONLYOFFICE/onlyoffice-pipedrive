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

/* eslint-disable react/jsx-props-no-spreading */
import React, { useState } from "react";
import { useDropzone, DropEvent, FileRejection } from "react-dropzone";
import cx from "classnames";

type DragDropProps = {
  onDrop: <T extends File>(
    acceptedFiles: T[],
    fileRejections: FileRejection[],
    event: DropEvent,
  ) => Promise<void>;
  errorText?: string;
  uploadingText?: string;
  selectText?: string;
  dragdropText?: string;
  subtext?: string;
  errorTimeout?: number;
};

const UploadIcon: React.FC = () => (
  <svg
    width="49"
    height="48"
    viewBox="0 0 49 48"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M32.5 32L24.5 24L16.5 32"
      stroke="#333333"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M24.5 24V42"
      stroke="#333333"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M41.2799 36.78C43.2306 35.7165 44.7716 34.0337 45.6597 31.9972C46.5477 29.9607 46.7323 27.6864 46.1843 25.5334C45.6363 23.3803 44.3869 21.471 42.6333 20.1069C40.8796 18.7427 38.7216 18.0014 36.4999 18H33.9799C33.3745 15.6585 32.2462 13.4846 30.6798 11.642C29.1134 9.79927 27.1496 8.33567 24.9361 7.36118C22.7226 6.3867 20.317 5.92669 17.9002 6.01573C15.4833 6.10478 13.1181 6.74057 10.9823 7.8753C8.84649 9.01003 6.99574 10.6142 5.56916 12.5671C4.14259 14.5201 3.1773 16.771 2.74588 19.1508C2.31446 21.5305 2.42813 23.977 3.07835 26.3065C3.72856 28.636 4.8984 30.7877 6.49992 32.6"
      stroke="#333333"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M32.5 32L24.5 24L16.5 32"
      stroke="#333333"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const OnlyofficeDragDrop: React.FC<DragDropProps> = ({
  onDrop,
  errorText = "Could not upload your file. Please contact ONLYOFFICE support.",
  uploadingText = "Uploading...",
  selectText = "Select a file",
  dragdropText = "or drag and drop here",
  subtext = "File size is limited",
  errorTimeout = 4000,
}) => {
  const [uploading, setUploading] = useState<boolean>(() => false);
  const [error, setError] = useState<boolean>(false);
  const uploadRef = React.useRef<HTMLInputElement | null>(null);

  const uploadFile = async (
    file: File | undefined,
    event: DropEvent,
    rejection?: FileRejection,
  ) => {
    setError(false);
    setUploading(true);
    if (file) {
      try {
        await onDrop([file], rejection ? [rejection] : [], event);
      } catch {
        setError(true);
        setTimeout(() => setError(false), errorTimeout);
      } finally {
        setUploading(false);
      }
    }
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop: (files, rejections, event) => {
      uploadFile(files[0], event, rejections[0]);
    },
    noClick: true,
    noKeyboard: true,
  });

  const style = cx({
    "flex flex-col items-center justify-center p-5": true,
    "border-2 border-slate-300 dark:border-dark-border border-dashed rounded-lg": true,
    "bg-transparent bg-opacity-20 text-black dark:text-dark-text": true,
    "transition-all transition-timing-function: ease-in-out": true,
    "transition-duration: 300ms": true,
    "bg-sky-100 dark:bg-sky-900": isDragActive,
    "bg-emerald-100 dark:bg-emerald-900": uploading,
    "bg-red-100 dark:bg-red-900": error,
  });

  return (
    <div className={`${style} w-full h-full`} {...getRootProps()}>
      <UploadIcon />
      {error && (
        <span className="font-sans font-semibold text-sm text-center text-black dark:text-dark-text">
          {errorText}
        </span>
      )}
      {uploading && !error && (
        <span className="font-sans font-semibold text-sm text-center text-black dark:text-dark-text">
          {uploadingText}
        </span>
      )}
      {!uploading && !error && (
        <>
          <input {...getInputProps()} />
          <input
            type="file"
            id="file"
            ref={uploadRef}
            style={{ display: "none" }}
            onChange={(e) => uploadFile(e.target?.files?.[0], e)}
          />
          <div className="font-sans font-semibold text-sm flex flex-wrap justify-center w-full">
            <button
              type="button"
              className="cursor-pointer outline-none border-b-2 border-dashed border-blue-500 text-blue-500 mr-1 max-w-max text-ellipsis truncate"
              onClick={() => uploadRef.current?.click()}
            >
              {selectText}
            </button>
            <span className="text-center max-w-max text-ellipsis truncate text-black dark:text-dark-text">
              {dragdropText}
            </span>
          </div>
          <span className="font-sans font-normal text-xs text-gray-400 dark:text-dark-muted text-center max-w-max text-ellipsis truncate">
            {subtext}
          </span>
        </>
      )}
    </div>
  );
};
