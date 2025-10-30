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

/* eslint-disable jsx-a11y/label-has-associated-control */
import React from "react";
import cx from "classnames";

type InputProps = {
  text: string;
  value?: string;
  placeholder?: string;
  type?: "text" | "password";
  errorText?: string;
  valid?: boolean;
  disabled?: boolean;
  textSize?: "sm" | "xs";
  labelSize?: "sm" | "xs";
  autocomplete?: boolean;
  onChange?: React.ChangeEventHandler<HTMLInputElement>;
};

export const OnlyofficeInput: React.FC<InputProps> = ({
  text,
  value,
  placeholder,
  type = "text",
  errorText = "Please fill out this field",
  valid = true,
  disabled = false,
  textSize = "sm",
  labelSize = "xs",
  autocomplete = false,
  onChange,
}) => {
  const istyle = cx({
    "font-normal text-sm text-gray-700 dark:text-dark-text appearance-none block select-auto": true,
    "text-xs": textSize === "xs",
    "w-full border rounded-sm h-10 px-4 bg-white dark:bg-dark-bg": true,
    "border-gray-light dark:border-dark-border": valid,
    "border-red-600": !valid,
    "bg-slate-200 dark:bg-dark-bg": disabled,
  });

  const pstyle = cx({
    hidden: valid,
  });

  return (
    <div>
      <label
        className={`font-semibold text-${labelSize} text-gray-700 dark:text-dark-text py-2`}
      >
        {text}
      </label>
      <input
        value={value}
        placeholder={placeholder}
        className={istyle}
        required
        autoCorrect={autocomplete ? undefined : "off"}
        autoComplete={autocomplete ? undefined : "off"}
        type={type}
        onChange={onChange}
        disabled={disabled}
      />
      <p className={`text-red-600 text-xs ${pstyle}`}>{errorText}</p>
    </div>
  );
};
