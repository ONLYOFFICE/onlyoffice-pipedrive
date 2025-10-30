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

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useMemo,
} from "react";
import AppExtensionsSDK from "@pipedrive/app-extensions-sdk";

type Theme = "light" | "dark";

interface ThemeContextType {
  theme: Theme;
  toggleTheme: () => void;
  isDark: boolean;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
};

interface ThemeProviderProps {
  children: React.ReactNode;
}

const initializePipedriveSDK = async (): Promise<Theme> => {
  try {
    const sdk = await new AppExtensionsSDK().initialize();
    const pipedriveTheme = sdk.userSettings.theme;
    return pipedriveTheme === "light" ? "light" : "dark";
  } catch {
    return "light";
  }
};

const applyThemeToDOM = (theme: Theme) => {
  const root = document.documentElement;
  root.classList.toggle("dark", theme === "dark");
};

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const [theme, setTheme] = useState<Theme>("light");

  const toggleTheme = () => {
    setTheme((prev) => (prev === "light" ? "dark" : "light"));
  };

  const isDark = theme === "dark";

  useEffect(() => {
    initializePipedriveSDK().then(setTheme);
  }, []);

  useEffect(() => {
    applyThemeToDOM(theme);
  }, [theme]);

  const value = useMemo(
    () => ({
      theme,
      toggleTheme,
      isDark,
    }),
    [theme, isDark],
  );

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
};
