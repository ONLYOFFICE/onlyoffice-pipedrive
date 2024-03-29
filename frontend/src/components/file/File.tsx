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

import React, { useState } from "react";

import DetailsIcon from "@assets/arrow-down.svg";

type FileProps = {
  Icon: any;
  name: string;
  supported?: boolean;
  actions?: React.ReactNode;
  children?: React.ReactNode;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
};

export const OnlyofficeFile: React.FC<FileProps> = ({
  Icon,
  name,
  supported = false,
  actions,
  children,
  onClick,
}) => {
  const [showDetails, setShowDetails] = useState(false);
  return (
    <>
      <div className="flex items-center w-full border-b py-2 my-1">
        <div className="flex items-center justify-center">
          <div
            role="button"
            tabIndex={0}
            onClick={() => setShowDetails(!showDetails)}
            onKeyDown={() => setShowDetails(!showDetails)}
            className={`w-[16px] h-[16px] hover:cursor-pointer mx-1 ${
              showDetails ? "rotate-180" : "rotate-0"
            }`}
          >
            <DetailsIcon />
          </div>
        </div>
        <div className="flex items-center justify-start w-3/4">
          <div className="w-[32px] h-[32px]">
            <Icon />
          </div>
          <button
            className={`${
              supported && onClick ? "cursor-pointer" : "cursor-default"
            } text-left font-semibold font-sans md:text-sm text-xs px-2 w-[170px] h-[32px] overflow-hidden text-ellipsis whitespace-nowrap`}
            type="button"
            title={name}
            onClick={onClick}
          >
            {name}
          </button>
        </div>
        <div className="flex items-center justify-end w-1/3">{actions}</div>
      </div>
      <div
        className={`overflow-hidden transition-all ${
          showDetails ? "h-[200px]" : "h-[0px]"
        }`}
      >
        {children}
      </div>
    </>
  );
};
