# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Frontend (must build before Go if embedding)
cd web && npm ci && npm run build && cd ..

# Go binary
go build -trimpath -o tinyurl .

# Run server
go run .

# Run with config file
go run . --config config.yaml

# Delete a short URL
go run . delete <code>

# Version info
go run . --version

# Build with version injection
go build -trimpath -ldflags="-X 'github.com/wxccs/tinyurl/global.Version=v0.1.0' -X 'github.com/wxccs/tinyurl/global.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)' -X 'github.com/wxccs/tinyurl/global.GitCommit=$(git rev-parse HEAD)'" -o tinyurl .
```

## Architecture

**MVC web service** built with Gin + Cobra + Viper + Gorm + Logrus. A Vue3 SPA frontend is embedded into the Go binary at compile time via `//go:embed web/dist`.

**Startup flow**: `main.go` (embeds frontend FS) → `cmd.Execute()` (cobra) → `runServer()` initializes logger → config → database → shortid generator → gin engine → `Routes.Setup()` registers all handlers → graceful shutdown on SIGINT/SIGTERM.

**Global state** (`global/`): Package-level vars (`Log`, `DB`, `Config`, `Generator`, `StaticFS`) shared across the app. Version/BuildTime/GitCommit are injected via `-ldflags` at build time.

**Config** (`internal/config/`): Viper with priority: CLI flags > env vars (`TU_` prefix, e.g. `TU_DATABASE_TYPE`) > config file > defaults. Config file searched at `~/.config/tinyurl/`, `/etc/tinyurl/`, `./config.yaml`.

**Database** (`internal/database/`): Gorm with SQLite (pure-Go, no CGO), MySQL, or PostgreSQL. Auto-migrates the `urls` table on init.

**Short ID generator** (`internal/shortid/`): 64-bit value = NodeID(4bit) | Timestamp(32bit) | Random(28bit), base62-encoded. Length ≥ 7, NodeID 0–15. Collisions retried up to 5 times at controller level.

**MVC layout**:
- `app/Controllers/` — each controller in its own subdirectory (e.g. `UrlController/`, `PageController/`). Use `NewResponse()` from `response.go` for JSON responses.
- `app/Models/` — Gorm models
- `app/Middlewares/` — Gin middleware
- `app/Routes/` — route registration in `Setup()`, static routes registered before `/:code` to avoid param capture

**Frontend** (`web/`): Vue3 + Element Plus + vue-i18n (32 locales). Vite dev server proxies `/api` to `localhost:8080`. Built output at `web/dist/` is embedded into the Go binary.

## Logging Convention

All logging uses `global.Log` with a `func` field formatted as `<package-path>.<function-name>`:
```go
global.Log.WithField("func", "app.Controllers.UrlController.Shorten").Info("url shortened")
```

## CI/CD

Tag push (`v*`) triggers `.github/workflows/release.yaml`:
1. Build 5-platform binaries (linux/darwin/windows × amd64 + linux/darwin arm64)
2. Publish to GitHub Releases (`.tar.gz` for Linux/macOS, `.zip` for Windows)
3. Build multi-arch Docker image (amd64+arm64) and push to `ghcr.io/wxccs/tinyurl`

## Docker

Multi-stage Dockerfile: Node build → Go build → minimal Alpine image with non-root `appuser`. Supports `TARGETOS`/`TARGETARCH` build args for multi-arch.

docker-compose.yaml defaults to SQLite with a named volume at `/app/data`. All config overridable via `TU_*` environment variables.
