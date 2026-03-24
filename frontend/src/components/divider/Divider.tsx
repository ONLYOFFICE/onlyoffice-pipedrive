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

type DividerProps = {
  text: string;
};

export const OnlyofficeDivider: React.FC<DividerProps> = ({ text }) => (
  <div className="relative flex py-5 items-center w-full">
    <div className="flex-grow border-t border-gray-400" />
    <span className="flex-shrink mx-4 text-gray-400">{text}</span>
    <div className="flex-grow border-t border-gray-400" />
  </div>
);
