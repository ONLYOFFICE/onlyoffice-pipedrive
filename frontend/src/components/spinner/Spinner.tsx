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

interface SpinnerProps {
  isDark?: boolean;
}

export const OnlyofficeSpinner: React.FC<SpinnerProps> = ({
  isDark = false,
}) => {
  if (isDark) {
    return (
      <div className="relative w-14 h-14">
        <div className="absolute inset-0 rounded-full border-4 border-gray-600" />
        <div className="absolute inset-0 rounded-full border-4 border-transparent border-t-white animate-spin" />
      </div>
    );
  }

  return (
    <div className="relative w-14 h-14">
      <div className="absolute inset-0 rounded-full border-4 border-gray-300" />
      <div className="absolute inset-0 rounded-full border-4 border-transparent border-t-gray-600 animate-spin" />
    </div>
  );
};
