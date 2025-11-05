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

export type TileProps = {
  Icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
  text: string;
  size?: "xs" | "sm";
  selected?: boolean;
  onClick?: React.MouseEventHandler<HTMLDivElement>;
  onDoubleClick?: React.MouseEventHandler<HTMLDivElement>;
  onKeyDown?: React.KeyboardEventHandler<HTMLDivElement>;
};

export const OnlyofficeTile: React.FC<TileProps> = ({
  Icon,
  text,
  size = "xs",
  selected = false,
  onClick,
  onDoubleClick,
  onKeyDown,
}) => {
  const card = cx({
    "px-5 py-3.5 rounded-lg transform shadow mb-0 outline-none": true,
    "transition duration-100 ease-linear": true,
    "h-[82px]": true,
    "max-h-36 flex flex-col justify-center": true,
    "border border-[#C8CCD2] dark:border-dark-border": true,
    "hover:bg-gray-100 dark:hover:bg-dark-hover cursor-pointer": !selected,
    "bg-white dark:bg-dark-bg": !selected,
    "bg-gray-200 dark:bg-dark-selected": selected,
  });

  const spn = cx({
    "text-sm": size === "sm",
    "text-xs text-[9px]": size === "xs",
    "font-semibold text-[#21232C] dark:text-dark-muted": true,
    "overflow-hidden whitespace-nowrap inline-block text-ellipsis": true,
  });

  return (
    <div
      role="button"
      tabIndex={0}
      className={card}
      onClick={onClick}
      onDoubleClick={onDoubleClick}
      onKeyDown={onKeyDown}
    >
      <div className="flex items-center justify-center px-1 py-1">
        <div className="relative flex items-end">
          <Icon />
        </div>
      </div>
      <div className="w-full flex items-center justify-center overflow-hidden">
        <span className={spn}>{text}</span>
      </div>
    </div>
  );
};
