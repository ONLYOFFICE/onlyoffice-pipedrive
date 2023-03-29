/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
    "hover:shadow-lg duration-200": !disabled,
    "bg-green-700 text-white": primary,
    "bg-white text-black border-2 border-slate-300 border-solid": !primary,
    "min-w-[62px] h-[32px]": true,
    "w-full": fullWidth,
    "bg-opacity-50 cursor-not-allowed": disabled,
  });

  return (
    <button
      type="button"
      disabled={disabled}
      className={`flex justify-center items-center p-3 text-sm lg:text-base font-bold rounded-md cursor-pointer ${classes} truncate text-ellipsis`}
      onClick={onClick}
    >
      {text}
      {Icon ? <div className="pl-1">{Icon}</div> : null}
    </button>
  );
};
