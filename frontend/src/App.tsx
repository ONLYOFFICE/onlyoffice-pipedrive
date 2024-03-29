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
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
  useLocation,
} from "react-router-dom";

import { MainPage } from "@pages/Main";
import { CreatePage } from "@pages/Creation";
import { SettingsPage } from "@pages/Settings";
import { OnlyofficeEditorPage } from "@pages/Editor";

import { OnlyofficeSpinner } from "@components/spinner";
import { TokenProvider } from "@context/TokenContext";

const CenteredOnlyofficeSpinner = () => (
  <div className="w-full h-full flex justify-center items-center">
    <OnlyofficeSpinner />
  </div>
);

const LazyRoutes: React.FC = () => {
  const location = useLocation();
  return (
    <Routes location={location} key={location.pathname}>
      <Route path="/">
        <Route
          index
          element={
            <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
              <TokenProvider>
                <MainPage />
              </TokenProvider>
            </React.Suspense>
          }
        />
        <Route
          path="create"
          element={
            <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
              <TokenProvider>
                <CreatePage />
              </TokenProvider>
            </React.Suspense>
          }
        />
        <Route
          path="settings"
          element={
            <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
              <TokenProvider>
                <SettingsPage />
              </TokenProvider>
            </React.Suspense>
          }
        />
      </Route>
      <Route
        path="/editor"
        element={
          <React.Suspense fallback={<CenteredOnlyofficeSpinner />}>
            <OnlyofficeEditorPage />
          </React.Suspense>
        }
      />
      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  );
};

function App() {
  return (
    <div className="w-full h-full flex justify-center items-center">
      <Router>
        <LazyRoutes />
      </Router>
    </div>
  );
}

export default App;
