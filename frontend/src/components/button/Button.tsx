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
    "duration-200 transition-all": true,
    "bg-green-700 dark:bg-dark-primary text-white": primary,
    "hover:bg-green-800 dark:hover:bg-dark-primary-hover": primary && !disabled,
    "bg-white dark:bg-dark-bg text-black dark:text-dark-text border border-[#C8CCD2] dark:border-dark-border":
      !primary,
    "hover:bg-[#F6F7F8] dark:hover:bg-dark-hover": !primary && !disabled,
    "active:bg-[#E4E6E9] dark:active:bg-dark-selected": !primary && !disabled,
    "shadow-[0px_1px_2px_0px_rgba(42,54,71,0.05)]": !primary && !disabled,
    "min-w-[62px] h-[32px]": true,
    "w-full": fullWidth,
    "cursor-pointer": !disabled,
    "cursor-not-allowed opacity-50 dark:opacity-40": disabled,
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
