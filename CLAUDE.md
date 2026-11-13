# CLAUDE.md - ONLYOFFICE Pipedrive integration

ONLYOFFICE app for Pipedrive (Marketplace app): ONLYOFFICE Docs editors for files attached to Pipedrive deals. Monorepo: Go microservices in `backend/` (one module `github.com/ONLYOFFICE/onlyoffice-pipedrive`) + React/TS frontend in `frontend/`, which runs inside Pipedrive iframes via `@pipedrive/app-extensions-sdk`. Background, history, deploy and gotchas: see `HANDOVER.md`. Releases: `/release` skill.

## Setup

```bash
git submodule update --init --recursive          # required: //go:embed needs submodules
cd backend/services/gateway/assets/aiautofill && make install-tools && make   # builds gitignored aiautofill/build, needed by gateway
```

Submodules: `backend/services/shared/format` + `frontend/src/assets/document-formats` (document-formats), `backend/services/gateway/assets/assets` (document-templates, branch `main/new`), `backend/services/gateway/assets/aiautofill` (plugin-aiautofill).

## Commands

Backend (from `backend/`, Go 1.25):

```bash
go build ./...
go test ./services/gateway/web/middleware/ ./services/callback/... ./services/builder/...   # smoke tests, no external deps
go test ./services/auth/web/core/service -run TestName                                     # single test
go run services/gateway/main.go server -c services/gateway/config/config.yml
```

- `*/core/adapter/mongo_test.go` need MongoDB on `localhost:27017`, otherwise they hang - don't run `go test ./...` without it.
- `TestSelectCaching` (auth handler) and `TestUserService/*timeout*` currently fail on `develop` - pre-existing, not caused by your change.

Frontend (from `frontend/`, Node 22, Yarn 1): `yarn install`, `yarn start:dev` (:3000), `yarn build`, `yarn eslint:fix`. No frontend tests exist (`yarn test` passes with none).

Frontend build-time env (`frontend/.env` / Docker build args): `BACKEND_GATEWAY`, `PIPEDRIVE_CREATE_MODAL_ID`. They are baked into the bundle - a frontend image is tied to one environment.

## Architecture (short)

- **gateway** (HTTP, chi, `:6044` local / `:4044` in Docker) - routes in `gateway/web/server.go`: `/oauth/*`, `/api/*`, `/files/*`, `/plugins/aiautofill/*`. `middleware/context.go` verifies the Pipedrive iframe JWT (`X-Pipedrive-App-Context`, signed with `client_secret`); `middleware/auth.go` - Basic auth for the uninstall webhook.
- **auth** (RPC) - OAuth tokens per user; **settings** (RPC) - per-company Document Server settings; **builder** (RPC) - signed editor config; **callback** (HTTP) - Docs save callback, uploads the file back to Pipedrive.
- Bootstrapping: `main.go` -> `cmd/server.go` -> `pkg.NewBootstrapper(..., pkg.WithModules(...))` from `onlyoffice-integration-adapters` (fx-style DI). New dependency = constructor + register in that module list.
- RPC: `client.Call(ctx, client.NewRequest("<namespace>:<service>", "<Handler>.<Method>", req), &resp)`, e.g. `"pipedrive:auth"`, `"UserSelectHandler.GetUser"`.
- `auth`, `settings`, gateway access store: hexagonal `web/core/{domain,port,service,adapter}`; Mongo adapter when `storage.url` is set, otherwise in-memory.
- `services/shared/`: config (YAML + `env:` overrides), Pipedrive clients, RPC DTOs, formats.
- Frontend pages: `Main` (deal panel), `Creation` (custom modal), `Settings`, `Editor` (opens in a new browser tab, not a modal). Path aliases in `tsconfig.paths.json`.

## Rules

- Do not change the user key format `fmt.Sprint(userID + companyID)` (a numeric **sum**) - it is the primary key of stored tokens and is used in auth, builder, callback, uninstall. Changing it needs a data migration.
- All services must get the same `CREDENTIALS_CLIENT_ID/SECRET` - the secret is also the key for the iframe JWT.
- Keep the demo-mode rule (enabled and `DemoStarted` empty or < 30 days) identical in builder, callback and gateway `api.go`.
- Do not remove `python3`/`kubernetes` from runtime images in `Dockerfile` - DevOps tooling uses them.
- Never commit `config.yml` or `frontend/.env` (gitignored, contain secrets).
- Every source file starts with the Apache-2.0 "(c) Copyright Ascensio System SIA" header.
- New dependencies must pass the `Licenses` workflow; manual approvals go to `frontend/.check-licenses.yml`.
- Branch from `develop`; `master` = releases. Pushes to `develop` build RC images (version from `CHANGELOG.md`) - update `CHANGELOG.md` for user-facing changes.
- Any new CI workflow that compiles the backend must init submodules and build aiautofill first.
