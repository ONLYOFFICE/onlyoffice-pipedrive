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
import cx from "classnames";

type SubtitleProps = {
  text: string;
  large?: boolean;
  center?: boolean;
};

export const OnlyofficeSubtitle: React.FC<SubtitleProps> = ({
  text,
  large = false,
  center = true,
}) => {
  const style = cx({
    "text-slate-700 dark:text-dark-muted font-normal": !!text,
    "text-center": center,
    "text-sm": !large,
    "text-base": large,
  });

  return <p className={style}>{text}</p>;
};
