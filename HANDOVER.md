# HANDOVER - ONLYOFFICE app for Pipedrive

## 1. Integration type

### What it is
A **separately hosted service** embedded into Pipedrive as a public Marketplace app: https://www.pipedrive.com/en/marketplace/app/onlyoffice/ab1362af3cc9f7be

User-facing parts (Pipedrive App Extensions, iframe + `@pipedrive/app-extensions-sdk`):

| Pipedrive surface | Frontend page | What it does |
|---|---|---|
| Deal panel ("ONLYOFFICE Documents") | `/` - `pages/Main` | Deal file list, search, rename, delete, download, open |
| Custom modal | `/create` - `pages/Creation` | Create a blank file or upload your own |
| App settings page | `/settings` - `pages/Settings` | Document Server address and secret, demo, plugin flags |
| New browser tab | `/editor` - `pages/Editor` | ONLYOFFICE editor (`@onlyoffice/document-editor-react`) |

The editor opens **not in a Pipedrive modal** but via `window.open`. The `PIPEDRIVE_EDITOR_MODAL_ID` build arg is a leftover from the old scheme and is not used in the code.

### Backend
5 Go services communicating via go-micro RPC with a service registry (`registry.type`: 1 = kubernetes, 2 = consul, 3 = etcd, 4 = mdns). Storage - MongoDB.

| Service | Protocol | Port (Docker) | Role |
|---|---|---|---|
| gateway | HTTP | 4044 | Public entry point: OAuth, frontend API, file creation and download, serving the AI Autofill plugin |
| auth | RPC | 5052 | User OAuth tokens, refresh |
| settings | RPC | 5150 | Per-company Document Server settings |
| builder | RPC | 6260 | Building and signing the editor config |
| callback | HTTP | 5044 | Save callback from Document Server -> upload file to Pipedrive |

Images: Docker Hub `onlyoffice/pipedrive-{frontend,gateway,auth,builder,callback,settings}`.

Environments (from `build.yaml`):
- **prod** - `https://pipedrive-gateway.onlyoffice.com`
- **dev** - `https://pipedrive-gateway.teamlab.info`

### Interaction with ONLYOFFICE
- **ONLYOFFICE Docs (Document Server)** via the Docs API:
  - builder generates the editor config and signs it with the company's Document Server JWT secret;
  - Document Server downloads the file via a direct Pipedrive link;
  - on save (status 2) it calls `callback`, which verifies the JWT and uploads the file to the deal.
- Each company specifies its own Document Server (self-hosted or Docs Cloud).
- **Demo mode**: ONLYOFFICE public demo Document Server for 30 days from activation. URL, secret and header are set in the builder and callback config (`ONLYOFFICE_DEMO_*`).
- **AI Autofill** (`ONLYOFFICE/plugin-aiautofill`, submodule):
  - the plugin is built into the gateway binary and served at `/plugins/aiautofill/*`;
  - auto-starts only for ONLYOFFICE PDF forms;
  - receives deal data (deal, person, organization, products) via `GET /api/data?code=...`, where code is a one-time access code (Mongo, TTL 30 days).
- Shared libs: `ONLYOFFICE/onlyoffice-integration-adapters` (public Go module on GitHub, v0.1.6): DI bootstrap, configs, JWT, logger, registry, storage. Also used by other ONLYOFFICE connectors (Monday, etc.).

### Versions
- Pipedrive: cloud product, no versions. Uses REST API `v1` (`/api/v1/files`, `/users/me`, deals, persons, organizations, products) and App Extensions SDK `^0.13.1`.
- ONLYOFFICE Docs: the minimum version is **not pinned in the code**.
  - JWT and plugins in the config are required.
  - AI Autofill and PDF forms require Docs with PDF form support, i.e. 8.x+ **[unconfirmed]**.
  - For support articles, refer to the current Docs version.

---

## 2. Main features (user scenarios)

1. **Installation**:
   - Marketplace -> `GET /oauth/install` - redirect to Pipedrive OAuth;
   - -> `GET /oauth/auth?code=` - tokens are stored in auth - redirect to Pipedrive.
2. **Uninstall**: Pipedrive calls `DELETE /oauth/auth` (Basic `client_id:client_secret`), the user's tokens are deleted.
3. **Settings** (Pipedrive global admin only):
   - Document Server address, secret and header; the server is validated on save;
   - or enabling demo;
   - "plugins" and "AI Autofill" flags.
4. **Deal file list**: pagination, search, rename, delete, download. The frontend calls the Pipedrive API directly with the user's token.
5. **Creating** blank DOCX/XLSX/PPTX/PDF from `document-templates` (by language) and **uploading** your own files (drag & drop).
6. **Opening in the editor**:
   - editing for formats from `document-formats`, everything else - view only;
   - co-editing, comments, chat;
   - printing is disabled (`Print: false`);
   - mobile mode is detected by User-Agent.
7. **Saving**: Pipedrive does not support file versions, so callback **uploads a new file** `Name (N).ext` to the same deal. The old file remains.
8. **AI Autofill** for ONLYOFFICE PDF forms.
9. Dark theme, 9 locales (`frontend/src/assets/locales`), user-voice banner and links.

---

## 3. Build

### Toolchain
- Go **1.25** (`go.mod`: 1.25.8; CI `licenses.yml`: 1.25; Docker: `golang:alpine` - floating tag!).
- Node **22** (Docker: `node:22.16.0-alpine`), **Yarn 1** (`yarn install --immutable`).
- Building aiautofill requires `make` and Node (`make install-tools && make`).

### Non-standard steps
1. `git submodule update --init --recursive` is mandatory: without it `//go:embed` won't compile.
2. Before building gateway you must build `backend/services/gateway/assets/aiautofill/build` (`make install-tools && make`). The directory is gitignored. The Dockerfile does **not** have this step, so CI builds it on the host before `docker build`. A local `docker-compose build` without this step will fail.
3. The frontend receives `BACKEND_GATEWAY` and `PIPEDRIVE_CREATE_MODAL_ID` **at build time**, so prod and dev need separate frontend images. The modal ID comes from the Pipedrive Developer Hub (App extensions - custom modal) and is stored in CI secrets `CREATE_MODAL_ID_PROD/DEV`.
4. Backend configs: `services/<svc>/config/config.example.yml` -> `config.yml` (gitignored). Any field can be overridden via env, see `shared/config.go` and env tags in `onlyoffice-integration-adapters/config`:
   - `CREDENTIALS_*`, `ONLYOFFICE_*`;
   - `STORAGE_URL`, `REGISTRY_TYPE`, `REGISTRY_ADDRESSES`, `SERVER_ADDRESS`, `LOGGER_LEVEL`.

   Running all services locally requires a shared registry: consul (`type: 2`) or mdns (`type: 4`).
5. Local development inside Pipedrive requires a public HTTPS gateway address (ngrok) and a test app in the Developer Hub with redirect URL `https://<ngrok>/oauth/auth`.

### Private dependencies and access
There are **no** private Go or npm packages: `GOPRIVATE` is empty, everything is installed from public GitHub, npm and proxy.golang.org. Access is needed to:

| What | Where | Owner |
|---|---|---|
| Main git | `https://git.onlyoffice.com/ONLYOFFICE/onlyoffice-pipedrive` (Gitea; PRs in the "Merge pull request ... from X into Y" style) | ONLYOFFICE integrations |
| GitHub mirror | `github.com/ONLYOFFICE/onlyoffice-pipedrive`: issues, Actions, Releases | same |
| CI secrets | GitHub Actions secrets: `DOCKER_HUB_USERNAME`, `DOCKER_HUB_ACCESS_TOKEN`, `CREATE_MODAL_ID_PROD`, `CREATE_MODAL_ID_DEV`, `PIPEDRIVE_WEBHOOK_URL`, `PIPEDRIVE_WEBHOOK_TOKEN` | repository admins **[unconfirmed]** |
| Docker Hub org `onlyoffice` | image push | DevOps |
| OAuth client_id/secret for prod and dev, demo Document Server secret, Mongo URL | environment secrets in the cluster (k8s secrets / env), locally - in the developer's `config.yml` | DevOps + app owner in Pipedrive **[unconfirmed]** |
| Pipedrive Developer Hub (prod and dev apps, sandbox companies) | developers.pipedrive.com | previous maintainer, transfer access |

Secrets are not stored in the repository and must not be written here.

---

## 4. Publishing

### Pipeline

- The version is taken from the first `## X.Y.Z` line in `CHANGELOG.md` (script `get_connector_version.py` from `ONLYOFFICE/ga-common`). `frontend/package.json` - `version` is updated in sync.
- After the build, `build.yaml` sends an HMAC-signed POST to `PIPEDRIVE_WEBHOOK_URL`. This triggers a stand update (CD on the DevOps side). Which environment gets updated with `type=prod` is **[unconfirmed]** - check with DevOps before the first prod release.
- `hooks/` - legacy Docker Hub autobuild hooks, not used by CI.

### Pipedrive Marketplace
The app is already published. **A regular release (code on our side) does not require a Pipedrive review**: rolling out new images is enough, since we host the frontend and backend.

Resubmission for review **[general practice]** is needed when changing what is stored on Pipedrive's side:
- OAuth scopes, callback/redirect URL;
- App extensions (panels, modals, settings page, their URLs);
- listing: description, screenshots, icons, support links, Privacy policy, Terms;
- publishing in new regions and languages.

Process **[general practice]**:
1. Developer Hub -> app -> edits in draft.
2. Test on a sandbox company.
3. "Submit for review".
4. The Pipedrive Marketplace team checks installation and OAuth, extensions, uninstall, listing.
5. Usually takes a few business days; remarks may come by email.

Typical rejection reasons **[general practice]**:
- installation or OAuth fails;
- uninstall is not handled;
- unnecessary scopes requested without justification;
- broken support and policy links;
- screenshots and description don't match the functionality;
- extension doesn't work for a non-admin user.

Manual steps when changing extensions:
- a new custom modal ID must be added to CI secrets `CREATE_MODAL_ID_PROD/DEV` - it is baked into the frontend;
- extension URLs in the Developer Hub must point to the frontend of the right environment.

---

## 5. Implicit knowledge

| What | Why / purpose | Can it be changed |
|---|---|---|
| User key = `fmt.Sprint(userID + companyID)`, an **arithmetic sum** | Historical decision. Collisions are theoretically possible | **No** without a Mongo migration. Used in auth (insert/select/delete), builder (`editorConfig.user.id`, token lookup), callback (`users[0]` from Docs -> token for upload), uninstall |
| `client_secret` is also the JWT key for `X-Pipedrive-App-Context` | That's how Pipedrive signs it | All 5 services must have the same credentials. There was already an incident: gateway with a different secret -> OAuth `invalid_client` |
| `ExpiresAt = now + expires_in - 700s` (gateway `auth.go`); auth cache = 1/6 of the remaining TTL | Early token refresh | Don't reduce the margin |
| Builder doesn't follow the Pipedrive download redirect and takes `Location` | Docs needs a direct link without authorization | Don't change `CheckRedirect` in `defaultClient` |
| Before uploading, callback does `PUT /files/{id}` with the **same** name | Probably an early check that the file exists and the token works (since the first commit) **[unconfirmed]** | Don't remove without testing on live Pipedrive |
| Save = new file `Name (N).ext` | Pipedrive has no file versions | The name format is visible to users, covered by `callback_test.go` |
| Callback uses the token of the first entry in `users` | The upload to Pipedrive is made on behalf of the last editor | - |
| The editor window opens empty, the URL is set after `GET_SIGNED_TOKEN` | Popup blocker workaround: `window.open` must be synchronous in the click | Don't call `window.open` after an await |
| A PDF form is detected by the `ONLYOFFICEFORM` signature in the first 300 bytes (Range request) | Cheap check without downloading the file | Covered by `build_test.go` |
| Demo is valid if `DemoStarted` is empty or younger than 30 days | Demo business rule | The logic is **duplicated** in builder, callback and gateway `api.go` - change in sync |
| `PluginsEnabled`/`AutofillEnabled` are `*bool`, `nil` = enabled | Compatibility with old settings records without these fields | Don't change to `bool` |
| `python3` + `kubernetes` in runtime images | DevOps utilities for checking the k8s API (commit ff62cb0) | Don't remove without DevOps |
| Gateway `NotFound` -> 301 to onlyoffice.com | To avoid exposing the API | - |
| The frontend gets the Pipedrive OAuth access token via `/api/me` and calls the Pipedrive API itself | File list and operations without proxying | The token is visible in the user's browser - this is intentional |

---

## 6. Known issues and tech debt

### Tech debt
- Tests:
  - two existing auth tests fail on develop (`TestSelectCaching/get_user`, `TestUserService/*_with_timeout`);
  - Mongo tests require a local MongoDB and hang without it;
  - there are no frontend tests.
- `middleware/context.go` and `middleware/auth.go` write `client_secret` to the debug log. Debug is off in prod (`debug: false`), but this should be removed.
- Error handling is inconsistent: `http.Error` in some places, bare `WriteHeader` in others, go-micro errors are parsed from a JSON string.
- Demo logic is duplicated in 3 places - move it to `shared`.
- `PIPEDRIVE_EDITOR_MODAL_ID` is a dead build arg; `hooks/` are dead hooks.
- Dockerfile:
  - `golang:alpine` without a version;
  - the runtime image is based on `golang:alpine` instead of `alpine`/distroless, so it's large;
  - the aiautofill step is outside the Dockerfile.
- `build.yaml`: the "Handle tag" step checks `$TAG`, which is not in env (should be `$MANUAL_TAG`). Currently RC_TAG is always set, and with a manual tag `inputs.tag` is used anyway, so it works by accident.
- Callback: every save adds a new file `(1)`, `(2)`, ... to the deal; old ones are not cleaned up. This is a Pipedrive limitation (no versions), but deals grow cluttered.

---

## 7. Plans and commitments

**No** commitments to customers or partners were found in the repository, issues or branches. The obvious next steps:
1. Merge `feature/error-page` and `fix/configuration`, release **1.2.2** (or 1.3.0 - the new error pages are user-visible).
2. Update AI Autofill following `plugin-aiautofill` releases (latest: v1.1.0, v1.1.1).
3. Update `document-formats` and `document-templates` following Docs releases.
4. Fix vulnerable dependencies (Dependabot, `Licenses` workflow).

Check customer or partner commitments and the roadmap with the ONLYOFFICE integrations product owner **[unconfirmed]**.

---

## Manual verification checklist (smoke before release)

Automated (`go test ./services/gateway/web/middleware/ ./services/callback/... ./services/builder/...`):
- iframe context JWT verification (401/403/200);
- uninstall Basic auth;
- rejection of malformed callbacks;
- `(N)` version names;
- 30-day demo rule;
- PDF form detection and AI Autofill auto-start.

Manually on the dev stand and a Pipedrive sandbox company:

- [ ] **Installation**: Marketplace / `/oauth/install` link -> OAuth -> return to Pipedrive without errors.
- [ ] **Settings** (admin): enter Document Server address and secret -> save with server validation; wrong address -> clear error. A non-admin user doesn't see the settings or can't save them.
- [ ] **Demo**: enable demo -> editor opens with the public demo server banner.
- [ ] **Deal panel**: file list, search, pagination (>50 files), rename, delete, download.
- [ ] **Creation**: DOCX, XLSX, PPTX and PDF form -> the file appears in the deal and opens in the editor (new tab, popup not blocked).
- [ ] **Upload** of your own file (drag & drop) -> appears in the list.
- [ ] **Editing**: make a change, close -> `Name (1).ext` with the change appears in the deal; again -> `(2)`.
- [ ] **Co-editing**: two users of the same company in the same file simultaneously.
- [ ] **Viewing**: a non-editable format (e.g. DJVU/XPS) opens view-only.
- [ ] **AI Autofill**: open an ONLYOFFICE PDF form in a deal with a person, organization and products -> the plugin starts and fills in the data; with the autofill flag off - it doesn't start.
- [ ] **Theme and language**: Pipedrive dark theme -> dark editor; user language ≠ en -> localized UI.
- [ ] **Mobile**: open the editor from a phone -> mobile mode.
- [ ] **App uninstall** -> `DELETE /oauth/auth` request returns 200, reinstalling works.
- [ ] Gateway, builder and callback logs have no errors; OAuth error pages (after merging `feature/error-page`) are displayed.
