# CLAUDE.md - ONLYOFFICE Pipedrive integration

## Project Overview

ONLYOFFICE app for Pipedrive: adds ONLYOFFICE Docs editors to Pipedrive deals. It is a monorepo with a Go microservices backend (`backend/`, single Go module `github.com/ONLYOFFICE/onlyoffice-pipedrive`) and a React/TypeScript frontend (`frontend/`) loaded inside Pipedrive via `@pipedrive/app-extensions-sdk` (custom modals / panels).

## Setup

Submodules are required — the backend will not compile without them because of `//go:embed` directives:

```bash
git submodule update --init --recursive
```

- `backend/services/shared/format` - document-formats (`shared/formats.go` embeds `format/onlyoffice-docs-formats.json`)
- `backend/services/gateway/assets/assets` - document-templates (`gateway/assets/embed.go`, blank files for "create")
- `backend/services/gateway/assets/aiautofill` - plugin-aiautofill (`gateway/assets/plugin.go` embeds `aiautofill/build`)
- `frontend/src/assets/document-formats` - document-formats

`aiautofill/build` is a generated (gitignored) directory; build it before compiling the gateway:

```bash
cd backend/services/gateway/assets/aiautofill && make install-tools && make
```

CI workflows (`build.yaml`, `licenses.yml`) run this same step - any new workflow that compiles the backend needs it too.

## Commands

Backend (run from `backend/`):

```bash
go build ./...
go test ./...
go test ./services/auth/web/core/service -run TestName   # single test
go run services/gateway/main.go server -c <path/to/config.yml>
```

Mongo adapter tests (`*/core/adapter/mongo_test.go`) connect to `mongodb://localhost:27017` and need a local MongoDB.

Frontend (run from `frontend/`, yarn):

```bash
yarn install
yarn start:dev        # webpack dev server on :3000
yarn build            # production build to frontend/build
yarn test             # jest
yarn test -- -t "name" # single test
yarn eslint:fix
```

Frontend build-time env (`frontend/.env`, via dotenv-webpack / Docker build args): `BACKEND_GATEWAY`, `PIPEDRIVE_CREATE_MODAL_ID`, `PIPEDRIVE_EDITOR_MODAL_ID`.

Docker: `Dockerfile` is multi-stage with one target per service (`gateway`, `auth`, `builder`, `callback`, `settings`, `frontend`); `docker-compose build` builds them all.

## Backend architecture

Five services under `backend/services/`, all bootstrapped the same way: `main.go` → `cmd.Run()` (urfave/cli) - `server` command - `pkg.NewBootstrapper(CONFIG_PATH, pkg.WithModules(...))` from `onlyoffice-integration-adapters`, which is a DI container (uber/fx-style constructors). To add a dependency, write a constructor and register it in the service's `cmd/server.go` module list.

- **gateway** (HTTP, chi) — public entry point for the frontend and Pipedrive. Routes in `gateway/web/server.go`: `/oauth/*` (install/auth/uninstall), `/api/*` (me, config, settings), `/files/*` (create from template, download, form check), `/plugins/aiautofill/*` (serves the embedded AI Autofill plugin). `middleware/context.go` validates the Pipedrive iframe JWT; `middleware/auth.go` protects uninstall. Also owns the AI access-code store (`web/core/`) used by the autofill plugin's `/api/data` endpoint.
- **auth** (go-micro RPC) — stores Pipedrive OAuth users/tokens and refreshes them. Handlers: `UserSelectHandler`, `UserInsertHandler`, `UserDeleteHandler`.
- **settings** (go-micro RPC) — per-company Document Server settings (address, secret, etc.). Handlers: `SettingsSelectHandler`, `SettingsInsertHandler`, `SettingsDeleteHandler`.
- **builder** (go-micro RPC) — builds the signed ONLYOFFICE editor config (fetches user + settings in parallel, signs with the Document Server JWT secret).
- **callback** (HTTP) — `POST /callback` receives ONLYOFFICE Document Server save callbacks and uploads the file back to Pipedrive.

Inter-service calls use go-micro: `client.Call(ctx, client.NewRequest("<namespace>:<service>", "<Handler>.<Method>", req), &resp)`, e.g. `"pipedrive:auth"`, `"UserSelectHandler.GetUser"`. Service names come from `namespace`/`name` in each service's config.

`auth`, `settings`, and gateway's access store follow a hexagonal layout under `web/core/`: `domain/` (entities), `port/` (interfaces), `service/` (business logic), `adapter/` (memory + Mongo implementations; `BuildNew*Adapter` picks Mongo when `storage.url` is set, else in-memory). RPC handlers use `singleflight` to dedupe concurrent lookups.

`services/shared/` holds cross-service code: config structs (`config.go`, YAML + env overrides via `env:` tags such as `CREDENTIALS_CLIENT_ID`, `ONLYOFFICE_*`), Pipedrive API/OAuth clients (`client/`), RPC request/response DTOs (`request/`, `response/`), and supported file formats.

Config: each service has `config/config.example.yml`; `config.yml` is gitignored and holds local values.

## Frontend architecture

React 19 + TypeScript, webpack, Tailwind, react-router, TanStack Query, valtio, i18next (locales in `src/assets/locales`). Path aliases (`@pages`, `@components`, `@context`, ...) are defined in `tsconfig.paths.json`. Pages map to Pipedrive surfaces: `Main` (deal panel file list), `Creation` (create/upload modal), `Editor` (editor modal using `@onlyoffice/document-editor-react`), `Settings`. `TokenContext` obtains the Pipedrive SDK token that is sent to the gateway; `src/services/` wraps gateway API calls.

## Conventions

- Every source file starts with the Apache-2.0 "(c) Copyright Ascensio System SIA" header — keep it on new files. The `Licenses` workflow checks dependency licenses for both `backend/` and `frontend/`.
- Branches: work on `develop`; `master` is the release branch. Pushes to `develop` build RC Docker images (version from `CHANGELOG.md`), so update `CHANGELOG.md` for user-facing changes.
