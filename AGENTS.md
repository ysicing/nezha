# Repository Guidelines

This repository hosts Nezha Monitoring (dashboard + agent). Use this guide to navigate the codebase and contribute effectively.

## Project Structure & Module Organization
- `cmd/dashboard`, `cmd/agent` — main entrypoints (HTTP+gRPC dashboard, agent CLI).
- `service/` — core services and orchestrators (e.g., `singleton/`, `rpc/`).
- `model/` — domain models, configuration, DB migrations via GORM.
- `pkg/` — shared packages (e.g., `monitor/`, `utils/`, `utls/`, `pty/`).
- `proto/` — protobuf definitions; generate stubs with `script/proto.sh`.
- `resource/` — embedded static assets, templates, i18n.
- `script/`, `.github/workflows/` — helper scripts and CI.

## Build, Test, and Development Commands
- Prerequisites: Go 1.23+ (go.mod uses 1.24), protoc if touching `proto/`.
- Build dashboard: `go build -o bin/dashboard ./cmd/dashboard`
- Build agent: `go build -o bin/nezha-agent ./cmd/agent`
- Run locally: `go run ./cmd/dashboard -c data/config.yaml --db data/sqlite.db`
- Run tests: `go test -v ./...` (CI also runs `gosec` on Linux).
- Generate protobuf: `bash script/proto.sh` (requires `protoc-gen-go`, `protoc-gen-go-grpc`).

## Coding Style & Naming Conventions
- Use standard Go formatting: `go fmt ./...`, `go vet ./...` (optional `goimports`).
- Packages: short, lowercase; files: `snake_case.go`.
- Exported identifiers in PascalCase; receivers consistent and meaningful.
- Group imports: stdlib, third‑party, then internal (`github.com/naiba/nezha/...`).

## Testing Guidelines
- Place tests in `*_test.go` with `TestXxx` names; prefer table‑driven tests.
- Keep tests deterministic; networked tests should have tight timeouts or be isolated with `-run` filters.
- Aim to cover critical paths in `model/`, `service/`, and `pkg/`.

## Commit & Pull Request Guidelines
- Use conventional prefixes when possible: `feat:`, `fix:`, `docs:`, `ci:`, `refactor:`, `security:`, `chore:`; keep messages imperative.
- Target branch: `next`. Ensure `go test` and a dashboard build succeed locally.
- PRs: concise description, linked issues, and screenshots/GIFs for UI/template changes in `resource/template/`.
- If protobuf or public APIs change, re‑generate stubs and include them in the PR.

## Security & Configuration Tips
- Do not commit secrets or built binaries (e.g., `agent`, `dashboard`).
- Dashboard config: `data/config.yaml` with environment overrides via `NZ_` prefix (example: `NZ_HTTPPORT=8080 go run ./cmd/dashboard`).
- Review OIDC/TLS settings under `model/config.go` before enabling in production.
