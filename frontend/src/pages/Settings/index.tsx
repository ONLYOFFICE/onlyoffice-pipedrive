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

import React, { useEffect, useState } from "react";
import AppExtensionsSDK, { Command } from "@pipedrive/app-extensions-sdk";
import { useSnapshot } from "valtio";
import { useTranslation } from "react-i18next";

import { OnlyofficeButton } from "@components/button";
import { OnlyofficeInput } from "@components/input";
import { OnlyofficeTitle } from "@components/title";
import { OnlyofficeSpinner } from "@components/spinner";
import { OnlyofficeBackgroundError } from "@layouts/ErrorBackground";
import { Banner } from "@layouts/Banner";

import { postSettings, getSettings } from "@services/settings";
import { getPipedriveMe } from "@services/me";

import { AuthToken } from "@context/TokenContext";

import OnlyofficeLogo from "@assets/onlyoffice-logo.svg";
import SettingsError from "@assets/settings-error.svg";
import { getCurrentURL } from "@utils/url";

function SettingsErrorIcon() {
  return (
    <div className="flex flex-col items-center justify-center">
      <OnlyofficeLogo />
      <SettingsError />
    </div>
  );
}

export const SettingsPage: React.FC = () => {
  const { t } = useTranslation();
  const [sdk, setSDK] = useState<AppExtensionsSDK | null>();
  const { access_token: accessToken, error, status } = useSnapshot(AuthToken);
  const [admin, setAdmin] = useState(false);
  const [loading, setLoading] = useState(true);
  const [address, setAddress] = useState<string | undefined>(undefined);
  const [secret, setSecret] = useState<string | undefined>(undefined);
  const [header, setHeader] = useState<string | undefined>(undefined);
  const [demoEnabled, setDemoEnabled] = useState(false);
  const [demoStarted, setDemoStarted] = useState<string | undefined>(undefined);
  const [saving, setSaving] = useState(false);

  const isDemoValid = (): boolean => {
    if (!demoEnabled) return false;

    if (
      !demoStarted ||
      demoStarted === "" ||
      demoStarted.startsWith("0001-01-01")
    )
      return true;

    const startDate = new Date(demoStarted);
    if (Number.isNaN(startDate.getTime())) return true;

    const thirtyDaysAgo = new Date();
    thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

    return startDate > thirtyDaysAgo;
  };

  const isDemoExpired = (): boolean => {
    if (
      !demoStarted ||
      demoStarted === "" ||
      demoStarted.startsWith("0001-01-01")
    )
      return false;

    const startDate = new Date(demoStarted);
    if (Number.isNaN(startDate.getTime())) return false;

    const thirtyDaysAgo = new Date();
    thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

    return startDate <= thirtyDaysAgo;
  };

  const getDemoStatus = (): string => {
    if (!demoEnabled) return "";

    if (
      !demoStarted ||
      demoStarted === "" ||
      demoStarted.startsWith("0001-01-01")
    )
      return t(
        "settings.demo.status.notstarted",
        "Demo will start when first used",
      );

    const startDate = new Date(demoStarted);
    if (Number.isNaN(startDate.getTime()))
      return t(
        "settings.demo.status.notstarted",
        "Demo will start when first used",
      );

    const daysAgo = Math.floor(
      (Date.now() - startDate.getTime()) / (1000 * 60 * 60 * 24),
    );
    const daysLeft = 30 - daysAgo;

    if (daysLeft > 0)
      return t(
        "settings.demo.status.active",
        "Demo active - {{days}} day(s) remaining",
        { days: daysLeft },
      );
    return t(
      "settings.demo.status.expired",
      "Demo has expired - please provide credentials",
    );
  };

  useEffect(() => {
    new AppExtensionsSDK()
      .initialize()
      .then((s) => {
        setSDK(s);
      })
      .catch(() => setSDK(null));
  }, []);

  useEffect(() => {
    if (accessToken && !error && !!sdk) {
      const { url } = getCurrentURL();
      getPipedriveMe(`${url}api/v1/users/me`)
        .then(async (ures) => {
          try {
            if (ures.data.access.find((a) => a.app === "global" && a.admin)) {
              const res = await getSettings(sdk);
              setAddress(res.doc_address);
              setSecret(res.doc_secret);
              setHeader(res.doc_header);
              setDemoEnabled(res.demo_enabled);
              setDemoStarted(res.demo_started);
              setAdmin(true);
            }
          } catch {
            setAdmin(false);
          } finally {
            setLoading(false);
          }
        })
        .catch(() => {
          setLoading(false);
        });
    }

    if (error) setLoading(false);
  }, [sdk, accessToken, error]);

  const handleSettings = async () => {
    if (sdk) {
      const hasCredentials = address && secret && header;
      const canSave = hasCredentials || isDemoValid();

      if (!canSave) {
        await sdk.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "settings.validation.error",
            "Please provide Document Server credentials or enable valid demo mode",
          ),
        });
        return;
      }

      if (address && !address.trim().toLowerCase().startsWith("https://")) {
        await sdk.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "settings.validation.https",
            "Document Server must use https protocol for Pipedrive integration",
          ),
        });
        return;
      }

      try {
        setSaving(true);
        const finalAddress =
          address && !address.endsWith("/") ? `${address}/` : address;
        await postSettings(
          sdk,
          finalAddress || "",
          secret || "",
          header || "",
          demoEnabled,
        );
        setDemoStarted(demoStarted || new Date().toISOString());
        await sdk.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "settings.saving.ok",
            "ONLYOFFICE settings have been saved",
          ),
        });
      } catch {
        await sdk.execute(Command.SHOW_SNACKBAR, {
          message: t(
            "settings.saving.error",
            "Could not save ONLYOFFICE settings",
          ),
        });
      } finally {
        setSaving(false);
      }
    }
  };

  return (
    <div className="custom-scroll w-screen h-screen overflow-y-scroll overflow-x-hidden bg-white dark:bg-dark-bg">
      {loading && !error && (
        <div className="h-full w-full flex justify-center items-center">
          <OnlyofficeSpinner />
        </div>
      )}
      {!loading && error && (
        <OnlyofficeBackgroundError
          Icon={<SettingsErrorIcon />}
          title={t("background.error.title.settings", "Something went wrong")}
          subtitle={t(
            status !== 401
              ? "background.error.subtitle"
              : "background.error.subtitle.token",
            status !== 401
              ? "Could not fetch plugin settings. Something went wrong. Please reload the pipedrive window"
              : "Could not fetch plugin settings. Something went wrong with your access token. Please reinstall the app",
          )}
          button={
            status === 401
              ? t("background.reinstall.button", "Reinstall") || "Reinstall"
              : t("button.reload", "Reload") || "Reload"
          }
          onClick={
            status === 401
              ? () => {
                  if (status === 401)
                    window.open(
                      `${getCurrentURL().url}settings/marketplace`,
                      "_blank",
                    );
                }
              : () => window.location.reload()
          }
        />
      )}
      {!loading && !error && !admin && (
        <OnlyofficeBackgroundError
          Icon={<SettingsErrorIcon />}
          title={t("background.access.title", "Access Denied")}
          subtitle={t(
            "background.access.subtitle",
            "Something went wrong or access denied",
          )}
          button={t("button.reload", "Reload") || "Reload"}
          onClick={() => window.location.reload()}
        />
      )}
      {!loading && !error && admin && (
        <>
          <div className="flex flex-col items-start pl-5 pr-5 pt-5 pb-3">
            <div className="pb-2">
              <OnlyofficeTitle
                text={t("settings.title", "Configure ONLYOFFICE app settings")}
              />
            </div>
            <p className="text-slate-800 dark:text-dark-text font-normal text-base text-left">
              {t(
                "settings.text",
                `
                The plugin which enables the users to edit office documents from
                Pipedrive using ONLYOFFICE Document Server, allows multiple users
                to collaborate in real time and to save back those changes to
                Pipedrive
              `,
              )}
            </p>
            <div
              className="flex items-center gap-4"
              style={{ marginTop: "10px" }}
            >
              <a
                href="https://helpcenter.onlyoffice.com/integration/pipedrive.aspx"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 text-sm font-medium transition-colors duration-200"
              >
                {t("settings.links.learnmore", "Learn more")} ↗
              </a>
              <a
                href="https://feedback.onlyoffice.com/forums/966080-your-voice-matters?category_id=519288"
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 text-sm font-medium transition-colors duration-200"
              >
                {t("settings.links.suggest", "Suggest a feature")} ↗
              </a>
            </div>
          </div>
          <div className="max-w-[320px]">
            <div className="pl-5 pr-5 pb-2">
              <OnlyofficeInput
                text={t("settings.inputs.address", "Document Server Address")}
                valid={
                  !address || address.trim() === ""
                    ? demoEnabled && isDemoValid()
                    : address.trim().toLowerCase().startsWith("https://")
                }
                errorText={t(
                  "settings.validation.https",
                  "Document Server must use https protocol for Pipedrive integration",
                )}
                disabled={saving}
                value={address}
                onChange={(e) => setAddress(e.target.value)}
              />
            </div>
            <div className="pl-5 pr-5 pb-2">
              <OnlyofficeInput
                text={t("settings.inputs.secret", "Document Server Secret")}
                valid={
                  !secret || secret.trim() === ""
                    ? demoEnabled && isDemoValid()
                    : true
                }
                errorText={
                  t(
                    "settings.inputs.error.secret",
                    "Document Server Secret is required",
                  ) || "Document Server Secret is required"
                }
                disabled={saving}
                value={secret}
                onChange={(e) => setSecret(e.target.value)}
                type="password"
              />
            </div>
            <div className="pl-5 pr-5">
              <OnlyofficeInput
                text={t("settings.inputs.header", "Document Server Header")}
                valid={
                  !header || header.trim() === ""
                    ? demoEnabled && isDemoValid()
                    : true
                }
                errorText={
                  t(
                    "settings.inputs.error.header",
                    "Document Server Header is required",
                  ) || "Document Server Header is required"
                }
                disabled={saving}
                value={header}
                onChange={(e) => setHeader(e.target.value)}
              />
            </div>
            <div className="pl-5 pr-5 mt-4">
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="demo-enabled"
                  checked={demoEnabled}
                  onChange={(e) => setDemoEnabled(e.target.checked)}
                  disabled={saving || isDemoExpired()}
                  className="w-4 h-4 text-blue-600 bg-gray-100 dark:bg-dark-bg border-gray-300 dark:border-dark-border rounded focus:ring-blue-500 focus:ring-2 disabled:opacity-50 disabled:cursor-not-allowed"
                />
                <label
                  htmlFor="demo-enabled"
                  className="ml-2 text-sm font-medium text-gray-900 dark:text-dark-text"
                >
                  {t("settings.inputs.demo", "Enable Demo Mode")}
                </label>
              </div>
              <p className="text-xs text-gray-500 dark:text-dark-muted mt-1 ml-6">
                {demoEnabled
                  ? getDemoStatus()
                  : t(
                      "settings.inputs.demo.description",
                      "Enable demo mode to test the integration without a Document Server",
                    )}
              </p>
            </div>
            <div className="flex justify-start items-center mt-4 ml-5">
              <OnlyofficeButton
                text={t("button.save", "Save")}
                primary
                disabled={
                  saving ||
                  (!address || address.trim() === ""
                    ? !(demoEnabled && isDemoValid())
                    : !address.trim().toLowerCase().startsWith("https://")) ||
                  (!secret || secret.trim() === ""
                    ? !(demoEnabled && isDemoValid())
                    : false) ||
                  (!header || header.trim() === ""
                    ? !(demoEnabled && isDemoValid())
                    : false)
                }
                onClick={handleSettings}
              />
            </div>
            <div className="relative bottom-0 ml-5 w-[568px]">
              <Banner />
            </div>
          </div>
        </>
      )}
    </div>
  );
};
