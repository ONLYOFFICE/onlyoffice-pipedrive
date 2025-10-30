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

module.exports = {
    content: ["./src/**/*.{ts,tsx}"],
    darkMode: 'class', // Enable class-based dark mode
    theme: {
        extend: {
            width: {
                1500: "1500px",
                1400: "1400px",
                1300: "1300px",
                1200: "1200px",
                1100: "1100px",
                1000: "1000px",
                900: "900px",
                800: "800px",
                700: "700px",
                600: "600px",
                536: "536px",
                454: "454px",
                400: "400px",
                327: "327px",
                300: "300px",
                200: "200px",
                197: "197px",
            },
            minWidth: {
                xxs: "14rem",
            },
            maxWidth: {
                xxs: "14rem",
            },
            height: {
                1500: "1500px",
                1400: "1400px",
                1300: "1300px",
                1200: "1200px",
                1100: "1100px",
                1000: "1000px",
                900: "900px",
                800: "800px",
                700: "700px",
                600: "600px",
                536: "536px",
                454: "454px",
                400: "400px",
                327: "327px",
                300: "300px",
                200: "200px",
                197: "197px",
            },
            colors: {
                onlyoffice: "#0F4071",
                "smoke-light": "rgba(0, 0, 0, 0.4)",
                dark: {
                    bg: "#1e2029",
                    surface: "#282C34",
                    border: "#444D56",
                    text: "#E6E6E8",
                    muted: "#B0B0B3",
                    primary: "#4CAF50",
                    hover: "#2C3038",
                    selected: "#3E444D"
                }
            },
            screens: {
                wrap: { raw: "(max-width: 713px)" },
                scrollable: { raw: "(max-height: 550px)" },
                large: { raw: "(min-height: 730px) and (min-width: 900px)" },
                small: { raw: "(max-width: 452px)" },
                xsmall: { raw: "(max-width: 300px)" },
                xxsmall: { raw: "(max-width: 150px)" },
            },
        },
        fontFamily: {
            open: ['"Open Sans"'],
        },
    },
    plugins: [],
};
