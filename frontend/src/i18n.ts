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

import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import ChainedBackend from "i18next-chained-backend";
import I18NextHttpBackend from "i18next-http-backend";
import LocalStorageBackend from "i18next-localstorage-backend";
import { initReactI18next } from "react-i18next";

// eslint-disable-next-line import/no-named-as-default-member
i18n
  .use(LanguageDetector)
  .use(ChainedBackend)
  .use(initReactI18next)
  .init({
    fallbackLng: "en-US",
    debug: false,
    interpolation: {
      escapeValue: false,
    },
    backend: {
      backends: [LocalStorageBackend, I18NextHttpBackend],
      backendOptions: [
        {
          expirationTime: 24 * 60 * 60 * 1000,
        },
      ],
    },
  });

export default i18n;
