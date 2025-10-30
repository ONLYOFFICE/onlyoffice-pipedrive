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

/* eslint-disable */
const path = require("path");
const { FlatCompat } = require("@eslint/eslintrc");
const js = require("@eslint/js");
const typescriptEslint = require("@typescript-eslint/eslint-plugin");
const typescriptParser = require("@typescript-eslint/parser");
const react = require("eslint-plugin-react");
const prettier = require("eslint-plugin-prettier");

const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
});

module.exports = [
    // Ignored files (from .eslintignore)
    {
        ignores: [
            "**/assets/**",
            "**/node_modules/**",
            "**/__tests__/**",
            "**/__snapshots__/**",
            "**/__mocks__/**",
            "**/setupTests.ts",
        ],
    },
    // Use FlatCompat to apply old-style extends configs
    ...compat.extends(
        "airbnb",
        "airbnb/hooks",
        "plugin:@typescript-eslint/recommended",
        "plugin:react/recommended",
        "plugin:import/errors",
        "plugin:import/warnings",
        "plugin:import/typescript",
        "prettier"
    ),
    // Apply custom config to TypeScript files (must come after extends to override)
    {
        files: ["**/*.ts", "**/*.tsx"],
        languageOptions: {
            parser: typescriptParser,
            parserOptions: {
                ecmaVersion: 2020,
                sourceType: "module",
                ecmaFeatures: {
                    jsx: true,
                },
                project: path.resolve(__dirname, "tsconfig.json"),
            },
            globals: {
                browser: true,
                node: true,
                es2020: true,
            },
        },
        plugins: {
            "@typescript-eslint": typescriptEslint,
            react: react,
            prettier: prettier,
        },
        settings: {
            "import/parsers": {
                "@typescript-eslint/parser": [".ts", ".tsx"],
            },
            "import/resolver": {
                "babel-module": {
                    extensions: [".js", ".jsx", ".ts", ".tsx"],
                },
                node: {
                    extensions: [".js", ".jsx", ".ts", ".tsx"],
                    paths: ["src"],
                },
                typescript: {
                    alwaysTryTypes: true,
                    project: path.resolve(__dirname, "tsconfig.json"),
                },
            },
        },
        rules: {
            "react/jsx-filename-extension": [1, { extensions: [".ts", ".tsx"] }],
            "import/extensions": "off",
            "react/prop-types": "off",
            "jsx-a11y/anchor-is-valid": "off",
            "react/jsx-props-no-spreading": ["error", { custom: "ignore" }],
            "prettier/prettier": "error",
            "react/no-unescaped-entities": "off",
            "import/no-cycle": [0, { ignoreExternal: true }],
            "prefer-const": "off",
            "no-use-before-define": "off",
            "react/function-component-definition": "off",
            "import/prefer-default-export": "off",
            "react/require-default-props": "off",
            "@typescript-eslint/no-use-before-define": [
                "error",
                { functions: false, classes: false, variables: true },
            ],
        },
    },
];

