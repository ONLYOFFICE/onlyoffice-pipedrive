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

type ButtonProps = {
  text: string;
  disabled?: boolean;
  primary?: boolean;
  fullWidth?: boolean;
  Icon?: React.ReactElement;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
};

export const OnlyofficeButton: React.FC<ButtonProps> = ({
  text,
  disabled = false,
  primary = false,
  fullWidth = false,
  Icon,
  onClick,
}) => {
  const classes = cx({
    "hover:shadow-lg duration-200 transition-all": !disabled,
    "bg-green-700 dark:bg-dark-primary text-white": primary,
    "bg-white dark:bg-dark-bg text-black dark:text-dark-text border-2 border-slate-300 dark:border-dark-border border-solid":
      !primary,
    "min-w-[62px] h-[32px]": true,
    "w-full": fullWidth,
    "cursor-pointer": !disabled,
    "bg-opacity-50 cursor-not-allowed": disabled,
    "dark:opacity-40 dark:brightness-75": disabled,
  });

  return (
    <button
      type="button"
      disabled={disabled}
      className={`flex justify-center items-center p-3 text-sm lg:text-base font-bold rounded-md ${classes} truncate text-ellipsis`}
      onClick={onClick}
    >
      {text}
      {Icon ? <div className="pl-1">{Icon}</div> : null}
    </button>
  );
};
