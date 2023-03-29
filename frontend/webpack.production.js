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
const { merge } = require("webpack-merge");
const webpack = require('webpack');
const dotenv = require('dotenv');
const common = require("./webpack.common.js");

module.exports = merge(common, {
    mode: "production",
    plugins: [
        new webpack.DefinePlugin({
            'process.env.BACKEND_GATEWAY': JSON.stringify(process.env.BACKEND_GATEWAY),
            'process.env.PIPEDRIVE_CREATE_MODAL_ID': JSON.stringify(process.env.PIPEDRIVE_CREATE_MODAL_ID),
            'process.env.PIPEDRIVE_EDITOR_MODAL_ID': JSON.stringify(process.env.PIPEDRIVE_EDITOR_MODAL_ID),
        }),
    ],
});
