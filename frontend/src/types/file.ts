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

export type FileStatus = "uploading" | "success" | "error" | "cancelled";

export interface UploadedFile {
  id: string;
  name: string;
  size: number;
  status: FileStatus;
  backendFileId?: number;
}

export type File = {
  id: string;
  user_id: number;
  name: string;
  file_size: number;
  file_type: string;
  add_time: string;
  update_time: string;
  url: string;
  person_name: string;
  remote_location: string;
};

type Pagination = {
  pagination: {
    start: number;
    next_start: number;
    limit: number;
    more_items_in_collection: boolean;
  };
};

export type FileResponse = {
  success: boolean;
  data: File[];
  additional_data: Pagination;
};
