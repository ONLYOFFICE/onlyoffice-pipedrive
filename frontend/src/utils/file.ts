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

import Docx from "@assets/docx.svg";
import Pdf from "@assets/pdf.svg";
import Pptx from "@assets/pptx.svg";
import Xlsx from "@assets/xlsx.svg";
import Vsd from "@assets/vsd.svg";
import Unsupported from "@assets/unsupported.svg";
import Supported from "@assets/supported.svg";

import wordIcon from "@assets/word.ico";
import slideIcon from "@assets/slide.ico";
import cellIcon from "@assets/cell.ico";
import pdfIcon from "@assets/pdf.ico";
import vsdIcon from "@assets/vsd.ico";
import genericIcon from "@assets/generic.ico";

import formatsData from "@assets/document-formats/onlyoffice-docs-formats.json";

interface Format {
  name: string;
  type: string;
  actions: string[];
  convert: string[];
  mime: string[];
}

const formats = formatsData as Format[];

const DOCUMENT_EXTS = formats
  .filter((f) => f.type === "word" && f.actions.includes("view"))
  .map((f) => f.name);

const SPREADSHEET_EXTS = formats
  .filter((f) => f.type === "cell" && f.actions.includes("view"))
  .map((f) => f.name);

const PRESENTATION_EXTS = formats
  .filter((f) => f.type === "slide" && f.actions.includes("view"))
  .map((f) => f.name);

const DIAGRAM_EXTS = formats
  .filter((f) => f.type === "diagram" && f.actions.includes("view"))
  .map((f) => f.name);

const PDF_EXTS = formats
  .filter((f) => f.type === "pdf" && f.actions.includes("view"))
  .map((f) => f.name);

const EDITABLE_EXTS = formats
  .filter((f) => f.actions.includes("edit"))
  .map((f) => f.name);

const OPENABLE_EXTS = DOCUMENT_EXTS.concat(SPREADSHEET_EXTS)
  .concat(PRESENTATION_EXTS)
  .concat(DIAGRAM_EXTS)
  .concat(PDF_EXTS);

const WORD = "word";
const SLIDE = "slide";
const CELL = "cell";
const DIAGRAM = "diagram";
const PDF = "pdf";

const getFileExt = (filename: string): string =>
  filename.split(".").pop() || "";

export const getFileParts = (filename: string): [string, string] => {
  const parts = filename.split(".");
  parts.pop();
  return [parts.join(".") || "invalid", getFileExt(filename)];
};

export const isFileEditable = (filename: string) => {
  const ext = getFileExt(filename).toLowerCase();
  return EDITABLE_EXTS.includes(ext);
};

export const isFileSupported = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();
  return OPENABLE_EXTS.includes(e);
};

export const getFileType = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();

  if (DOCUMENT_EXTS.includes(e)) return WORD;
  if (SPREADSHEET_EXTS.includes(e)) return CELL;
  if (PRESENTATION_EXTS.includes(e)) return SLIDE;
  if (DIAGRAM_EXTS.includes(e)) return DIAGRAM;
  if (PDF_EXTS.includes(e)) return PDF;

  return null;
};

export const getMimeType = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();

  const format = formats.find((f) => f.name === e);
  if (format && format.mime.length > 0) {
    return format.mime[0];
  }

  switch (e) {
    case "docx":
      return "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
    case "pptx":
      return "application/vnd.openxmlformats-officedocument.presentationml.presentation";
    case "xlsx":
      return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";
    default:
      return "application/octet-stream";
  }
};

export const getFileIcon = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();

  if (e === "pdf") return Pdf;
  if (DOCUMENT_EXTS.includes(e)) return Docx;
  if (SPREADSHEET_EXTS.includes(e)) return Xlsx;
  if (PRESENTATION_EXTS.includes(e)) return Pptx;
  if (DIAGRAM_EXTS.includes(e)) return Vsd;
  if (OPENABLE_EXTS.includes(e)) return Supported;

  return Unsupported;
};

export const formatBytes = (bytes: number, decimals = 2) => {
  if (!+bytes) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / k ** i).toFixed(dm))} ${sizes[i]}`;
};

export const getFileFavicon = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();

  if (PDF_EXTS.includes(e)) return pdfIcon;
  if (DOCUMENT_EXTS.includes(e)) return wordIcon;
  if (PRESENTATION_EXTS.includes(e)) return slideIcon;
  if (SPREADSHEET_EXTS.includes(e)) return cellIcon;
  if (DIAGRAM_EXTS.includes(e)) return vsdIcon;

  return genericIcon;
};

export const isPDF = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();
  return e === "pdf";
};
