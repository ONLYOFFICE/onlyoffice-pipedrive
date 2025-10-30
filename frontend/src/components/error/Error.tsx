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
import cx from "classnames";

type ErrorProps = {
  text: string;
  isDark?: boolean;
};

export const OnlyofficeError: React.FC<ErrorProps> = ({ text, isDark }) => {
  const textStyle = cx("font-semibold text-center flex items-center", {
    "text-white": isDark === true,
    "text-gray-900": isDark === false,
    "text-gray-900 dark:text-white": isDark === undefined,
  });

  return (
    <div className="flex flex-col justify-center items-center">
      <span className={textStyle}>{text}</span>
    </div>
  );
};
