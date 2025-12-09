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

import React, { useState } from "react";

import { useTheme } from "@context/ThemeContext";

import DetailsIcon from "@assets/arrow-down.svg";
import DetailsIconDark from "@assets/arrow-down_dark.svg";

type FileProps = {
  Icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
  name: string;
  supported?: boolean;
  actions?: React.ReactNode;
  children?: React.ReactNode;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
  isRenaming?: boolean;
  onRenameSubmit?: (newName: string) => void;
  onRenameCancel?: () => void;
};

export const OnlyofficeFile: React.FC<FileProps> = ({
  Icon,
  name,
  supported = false,
  actions,
  children,
  onClick,
  isRenaming = false,
  onRenameSubmit,
  onRenameCancel,
}) => {
  const [showDetails, setShowDetails] = useState(false);
  const [renameValue, setRenameValue] = useState(name);
  const inputRef = React.useRef<HTMLInputElement>(null);
  const { isDark } = useTheme();

  React.useEffect(() => {
    if (isRenaming && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isRenaming]);

  React.useEffect(() => {
    setRenameValue(name);
  }, [name]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      onRenameSubmit?.(renameValue.trim());
    } else if (e.key === "Escape") {
      e.preventDefault();
      onRenameCancel?.();
    }
  };

  const handleBlur = () => {
    if (isRenaming) {
      const trimmedValue = renameValue.trim();
      if (trimmedValue !== "" && trimmedValue !== name) {
        onRenameSubmit?.(trimmedValue);
      } else {
        onRenameCancel?.();
      }
    }
  };

  return (
    <>
      <div className="flex items-center w-full border-b dark:border-dark-border py-2 my-1">
        <div className="flex items-center justify-center">
          <div
            role="button"
            tabIndex={0}
            onClick={() => setShowDetails(!showDetails)}
            onKeyDown={() => setShowDetails(!showDetails)}
            className={`w-[16px] h-[16px] hover:cursor-pointer mx-1 text-black dark:text-dark-text ${
              showDetails ? "rotate-180" : "rotate-0"
            }`}
          >
            {isDark ? (
              <DetailsIconDark className="fill-current" />
            ) : (
              <DetailsIcon className="fill-current" />
            )}
          </div>
        </div>
        <div className="flex items-center justify-start flex-1 min-w-0">
          <div className="w-[32px] h-[32px] flex-shrink-0">
            <Icon />
          </div>
          {isRenaming ? (
            <input
              ref={inputRef}
              type="text"
              value={renameValue}
              onChange={(e) => setRenameValue(e.target.value)}
              onKeyDown={handleKeyDown}
              onBlur={handleBlur}
              className="text-left font-semibold font-sans text-sm px-2 flex-1 h-[32px] text-black dark:text-dark-text bg-white dark:bg-dark-bg border border-blue-500 dark:border-blue-400 rounded outline-none"
            />
          ) : (
            <button
              className={`${
                supported && onClick ? "cursor-pointer" : "cursor-default"
              } text-left font-semibold font-sans text-sm px-2 flex-1 h-[32px] overflow-hidden text-ellipsis whitespace-nowrap text-black dark:text-dark-text min-w-0`}
              type="button"
              title={name}
              onClick={onClick}
              disabled={isRenaming}
            >
              {name}
            </button>
          )}
        </div>
        <div className="flex items-center justify-end flex-shrink-0">
          {actions}
        </div>
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
