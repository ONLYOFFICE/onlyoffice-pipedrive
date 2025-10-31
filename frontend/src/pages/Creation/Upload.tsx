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

import React, { useEffect, useState } from "react";
import AppExtensionsSDK, { Command } from "@pipedrive/app-extensions-sdk";
import { useTranslation } from "react-i18next";

import { OnlyofficeDragDrop } from "@components/drop";
import { UploadItem, Footer } from "@components/upload";

import { uploadFile, deleteFile } from "@services/file";

import { getCurrentURL } from "@utils/url";

import type { FileStatus, UploadedFile } from "@objects/file";

const MAX_FILES_LIMIT = 5;

const DropZone: React.FC<{
  errorText: string;
  uploadingText: string;
  selectText: string;
  dragdropText: string;
  subtext: string;
  onDrop: (files: File[]) => Promise<void>;
}> = ({
  errorText,
  uploadingText,
  selectText,
  dragdropText,
  subtext,
  onDrop,
}) => (
  <div className="h-48 flex-shrink-0">
    <OnlyofficeDragDrop
      errorText={errorText}
      uploadingText={uploadingText}
      selectText={selectText}
      dragdropText={dragdropText}
      subtext={subtext}
      onDrop={onDrop}
    />
  </div>
);

const FilesList: React.FC<{
  files: UploadedFile[];
  onRemoveFile: (fileId: string) => void;
  onCancelUpload: (fileId: string) => void;
}> = ({ files, onRemoveFile, onCancelUpload }) => {
  if (files.length === 0) return null;

  return (
    <div className="mt-4 flex-1 overflow-y-auto">
      {files.map((file) => (
        <UploadItem
          key={file.id}
          file={file}
          onRemove={onRemoveFile}
          onCancel={onCancelUpload}
        />
      ))}
    </div>
  );
};

export const Upload: React.FC = () => {
  const { t } = useTranslation();
  const [sdk, setSDK] = useState<AppExtensionsSDK | null>();
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [abortControllers, setAbortControllers] = useState<
    Map<string, AbortController>
  >(new Map());

  useEffect(() => {
    new AppExtensionsSDK()
      .initialize()
      .then((s) => setSDK(s))
      .catch(() => setSDK(null));
  }, []);

  const handleCancelUpload = (fileId: string) => {
    const controller = abortControllers.get(fileId);
    if (controller) {
      controller.abort();
      setAbortControllers((prev) => {
        const updated = new Map(prev);
        updated.delete(fileId);
        return updated;
      });
      setUploadedFiles((prev) =>
        prev.map((f) =>
          f.id === fileId ? { ...f, status: "cancelled" as FileStatus } : f,
        ),
      );
    }
  };

  const handleRemoveFile = async (fileId: string) => {
    const file = uploadedFiles.find((f) => f.id === fileId);
    if (file?.backendFileId && file.status === "success") {
      const { url } = getCurrentURL();
      try {
        await deleteFile(`${url}api/v1/files/${file.backendFileId}`);
        await sdk?.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "snackbar.fileremoved.ok",
            "File {{file}} has been removed",
            { file: file.name },
          ),
        });
      } catch {
        await sdk?.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "snackbar.fileremoved.error",
            "Could not remove file {{file}}",
            { file: file.name },
          ),
        });
        return;
      }
    }

    setUploadedFiles((prev) => {
      const updated = prev.filter((f) => f.id !== fileId);
      if (updated.length === 0) {
        sdk?.execute(Command.RESIZE, {
          height: 424,
          width: 622,
        });
      }

      return updated;
    });
  };

  const validateFileCount = async (fileCount: number): Promise<boolean> => {
    if (fileCount > MAX_FILES_LIMIT) {
      await sdk?.execute(Command.SHOW_SNACKBAR, {
        message: t(
          "upload.limit.error",
          "You can upload a maximum of {{count}} files at once",
          { count: MAX_FILES_LIMIT },
        ),
      });
      return false;
    }
    return true;
  };

  const createUploadedFiles = (files: File[]): UploadedFile[] => {
    const timestamp = Date.now();
    return files.map((file, index) => ({
      id: `${timestamp}-${index}-${Math.random().toString(36).slice(2)}`,
      name: file.name,
      size: file.size,
      status: "uploading" as FileStatus,
    }));
  };

  const uploadSingleFile = async (
    file: File,
    fileId: string,
    url: string,
    dealId: string,
  ): Promise<{ success: boolean; fileName: string; cancelled?: boolean }> => {
    const controller = new AbortController();
    setAbortControllers((prev) => new Map(prev).set(fileId, controller));

    try {
      const response = await uploadFile(
        `${url}api/v1/files`,
        dealId,
        file,
        controller.signal,
      );
      setAbortControllers((prev) => {
        const updated = new Map(prev);
        updated.delete(fileId);
        return updated;
      });
      setUploadedFiles((prev) =>
        prev.map((f) =>
          f.id === fileId
            ? { ...f, status: "success", backendFileId: response.data?.id }
            : f,
        ),
      );
      return { success: true, fileName: file.name };
    } catch (error: unknown) {
      setAbortControllers((prev) => {
        const updated = new Map(prev);
        updated.delete(fileId);
        return updated;
      });

      if (
        (error as { code?: string; name?: string }).code === "ERR_CANCELED" ||
        (error as { code?: string; name?: string }).name === "CanceledError"
      ) {
        return { success: false, fileName: file.name, cancelled: true };
      }

      setUploadedFiles((prev) =>
        prev.map((f) => (f.id === fileId ? { ...f, status: "error" } : f)),
      );
      return { success: false, fileName: file.name };
    }
  };

  const uploadAllFiles = async (
    files: File[],
    newFiles: UploadedFile[],
    url: string,
    dealId: string,
  ) => {
    const results = await Promise.all(
      files.map((file, i) =>
        uploadSingleFile(file, newFiles[i].id, url, dealId),
      ),
    );

    return {
      successful: results.filter((r) => r.success).map((r) => r.fileName),
      failed: results
        .filter((r) => !r.success && !r.cancelled)
        .map((r) => r.fileName),
      cancelled: results.filter((r) => r.cancelled).map((r) => r.fileName),
    };
  };

  const showUploadNotifications = async (
    successful: string[],
    failed: string[],
  ) => {
    if (successful.length > 0) {
      await sdk?.execute(Command.SHOW_SNACKBAR, {
        message:
          successful.length === 1
            ? t("snackbar.uploaded.ok", "File {{file}} has been uploaded", {
                file: successful[0],
              })
            : t(
                "snackbar.uploaded.multiple.ok",
                "{{count}} files have been uploaded successfully",
                { count: successful.length },
              ),
      });
    }

    if (failed.length > 0) {
      await sdk?.execute(Command.SHOW_SNACKBAR, {
        message:
          failed.length === 1
            ? t("snackbar.uploaded.error", "Could not upload file {{file}}", {
                file: failed[0],
              })
            : t(
                "snackbar.uploaded.multiple.error",
                "Could not upload {{count}} file(s)",
                { count: failed.length },
              ),
      });
    }
  };

  const handleFilesUpload = async (files: File[]) => {
    if (!(await validateFileCount(files.length))) return;

    const { url, parameters } = getCurrentURL();
    const dealId = parameters.get("selectedIds") || "";

    const newFiles = createUploadedFiles(files);
    setUploadedFiles((prev) => [...prev, ...newFiles]);

    await sdk?.execute(Command.RESIZE, {
      height: 500,
      width: 622,
    });

    const { successful, failed } = await uploadAllFiles(
      files,
      newFiles,
      url,
      dealId,
    );

    await showUploadNotifications(successful, failed);
  };

  const handleDrop = async (files: File[]) => {
    try {
      await handleFilesUpload(files);
      return Promise.resolve();
    } catch {
      return Promise.reject();
    }
  };

  const handleClose = async () => {
    await sdk?.execute(Command.CLOSE_MODAL);
  };

  return (
    <div className="h-full bg-white dark:bg-dark-bg flex flex-col">
      <div className="flex-1 px-5 py-5 flex flex-col overflow-hidden">
        <DropZone
          errorText={
            t(
              "upload.error",
              "Could not upload your file. Please contact ONLYOFFICE support.",
            ) ||
            "Could not upload your file. Please contact ONLYOFFICE support."
          }
          uploadingText={
            t("upload.uploading", "Uploading...") || "Uploading..."
          }
          selectText={t("upload.select", "Select file(s)") || "Select file(s)"}
          dragdropText={
            t("upload.dragdrop", "or drag and drop here") ||
            "or drag and drop here"
          }
          subtext={
            t(
              "upload.subtext",
              "File size must not exceed 20MB. Max {{count}} files at once",
              { count: MAX_FILES_LIMIT },
            ) ||
            `File size must not exceed 20MB. Max ${MAX_FILES_LIMIT} files at once`
          }
          onDrop={handleDrop}
        />
        <div className="mt-4 flex-shrink-0">
          <p className="text-xs text-gray-600 dark:text-gray-400 leading-relaxed">
            {t(
              "upload.info",
              "Files will be automatically uploaded, attached to this Pipedrive deal, and ready for editing in ONLYOFFICE.",
            )}
          </p>
        </div>
        <FilesList
          files={uploadedFiles}
          onRemoveFile={handleRemoveFile}
          onCancelUpload={handleCancelUpload}
        />
      </div>
      <Footer
        closeButtonText={t("button.close", "Close")}
        onClose={handleClose}
      />
    </div>
  );
};
