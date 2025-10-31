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

import { OnlyofficeButton } from "@components/button";

export const Footer: React.FC<{
  closeButtonText: string;
  onClose: () => void;
}> = ({ closeButtonText, onClose }) => {
  return (
    <div className="h-[48px] flex items-center w-full bg-white dark:bg-dark-bg border-t dark:border-dark-border">
      <div className="flex justify-between items-center w-full px-5">
        <OnlyofficeButton text={closeButtonText} onClick={onClose} />
      </div>
    </div>
  );
};
