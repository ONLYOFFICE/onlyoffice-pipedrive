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

import { OnlyofficeButton } from "@components/button";
import { OnlyofficeSubtitle } from "@components/title";
import { OnlyofficeError } from "@components/error/Error";

type ErrorProps = {
  Icon: React.ReactNode;
  title: string;
  subtitle: string;
  button?: string;
  onClick?: React.MouseEventHandler<HTMLButtonElement> | undefined;
};

export const OnlyofficeBackgroundError: React.FC<ErrorProps> = ({
  Icon,
  title,
  subtitle,
  button,
  onClick,
}) => (
  <div className="w-full h-full flex justify-center flex-col items-center overflow-hidden bg-white dark:bg-dark-bg">
    <div className="flex justify-center items-center overflow-hidden">
      {Icon}
    </div>
    <div>
      <OnlyofficeError text={title} />
    </div>
    <div className="w-1/2 pt-2">
      <OnlyofficeSubtitle text={subtitle} />
    </div>
    {onClick && button && (
      <div className="pt-5 z-[100]">
        <OnlyofficeButton primary text={button} onClick={onClick} />
      </div>
    )}
  </div>
);
