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

import Nofiles from "@assets/nofile.svg";

type NoFileProps = {
  title: string;
};

export const OnlyofficeNoFile: React.FC<NoFileProps> = ({ title }) => (
  <div className="h-full w-full flex flex-col justify-center items-center">
    <Nofiles />
    <span className="font-sans font-bold text-sm max-w-max text-ellipsis truncate">
      {title}
    </span>
  </div>
);
