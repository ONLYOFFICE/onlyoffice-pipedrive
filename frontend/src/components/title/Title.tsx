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

type TitleProps = {
  text: string;
  large?: boolean;
  align?: "left" | "center" | "right";
};

export const OnlyofficeTitle: React.FC<TitleProps> = ({
  text,
  large = false,
  align = "center",
}) => {
  const style = cx({
    "font-semibold text-slate-800 dark:text-dark-text": !!text,
    "text-center": align === "center",
    "text-left": align === "left",
    "text-right": align === "right",
    "text-base": large,
    "text-sm": !large,
  });

  return <p className={style}>{text}</p>;
};
