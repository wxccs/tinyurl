# TinyURL

A self-hosted URL shortener service built with Go and Vue 3.

## Features

- Generate short URLs with configurable length (default 7 characters)
- Instant redirect via short code
- Support SQLite, MySQL, and PostgreSQL
- Vue 3 + Element Plus frontend with 32 languages
- Single binary with embedded frontend (no separate deployment needed)
- Graceful shutdown
- Docker support with multi-architecture images

## Quick Start

### Binary

Download from [GitHub Releases](https://github.com/wxccs/tinyurl/releases), then:

```bash
./tinyurl
```

Server starts at `http://0.0.0.0:8080` by default.

### Docker

```bash
docker compose up -d
```

## Configuration

Configuration sources in priority order: CLI flags > environment variables > config file > defaults.

### CLI Flags

```
--config      Config file path (default $HOME/.config/tinyurl/config.yaml)
--log-level   Log level: 0=Panic, 1=Fatal, 2=Error, 3=Warn, 4=Info, 5=Debug, 6=Trace (default 4)
--log-file    Log file path (default empty, logs to stdout only)
```

### Environment Variables

Prefix with `TU_`, replace `.` and `-` with `_`:

| Variable | Default | Description |
|---|---|---|
| `TU_SERVER_HOST` | `0.0.0.0` | Listen host |
| `TU_SERVER_PORT` | `8080` | Listen port |
| `TU_DATABASE_TYPE` | `sqlite` | Database type: sqlite, mysql, postgres |
| `TU_DATABASE_HOST` | | Database host (mysql/postgres) |
| `TU_DATABASE_PORT` | `3306`/`5432` | Database port (auto-defaults by type) |
| `TU_DATABASE_USER` | | Database user (mysql/postgres) |
| `TU_DATABASE_PASSWORD` | | Database password (mysql/postgres) |
| `TU_DATABASE_DBNAME` | `tinyurl` | Database name (mysql/postgres) |
| `TU_DATABASE_PATH` | `data/tinyurl.db` | Database file path (sqlite) |
| `TU_SHORTURL_LENGTH` | `7` | Short code length (minimum 7) |
| `TU_SHORTURL_NODE_ID` | `0` | Node ID for distributed deployment (0-15) |
| `TU_PAGE_TITLE` | `TinyURL - Short URL Generator` | Page title |
| `TU_BEIAN_MIIT` | | MIIT filing number |
| `TU_BEIAN_MPS` | | MPS filing number |

### Config File

Create `config.yaml` in one of the search paths (`~/.config/tinyurl/`, `/etc/tinyurl/`, `./`):

```yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  type: sqlite
  path: data/tinyurl.db

shorturl:
  length: 7
  node_id: 0
```

## API

### Shorten URL

```
POST /api/shorten
Content-Type: application/json

{"url": "https://example.com/very/long/url"}
```

Response:

```json
{
  "code": 0,
  "msg": "success",
  "data": {"short_url": "http://localhost:8080/Ab3x9Kp"}
}
```

### Redirect

```
GET /:code
```

Returns `302` redirect to the original URL.

### Public Config

```
GET /api/config
```

Returns page title and filing information.

## Build from Source

### Prerequisites

- Go 1.26+
- Node.js 22+

### Steps

```bash
# Build frontend
cd web && npm ci && npm run build && cd ..

# Build binary
go build -trimpath -o tinyurl .

# With version info
go build -trimpath \
  -ldflags="-X 'github.com/wxccs/tinyurl/global.Version=$(git describe --tags)' \
    -X 'github.com/wxccs/tinyurl/global.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
    -X 'github.com/wxccs/tinyurl/global.GitCommit=$(git rev-parse HEAD)'" \
  -o tinyurl .
```

## Docker

### docker-compose.yaml

Default uses SQLite with a named volume. Uncomment the MySQL/PostgreSQL environment variables to switch databases:

```yaml
services:
  tinyurl:
    image: ghcr.io/wxccs/tinyurl:latest
    ports:
      - "8080:8080"
    volumes:
      - tinyurl-data:/app/data
    environment:
      - TU_DATABASE_TYPE=sqlite
      - TU_DATABASE_PATH=/app/data/tinyurl.db
      # MySQL:
      # - TU_DATABASE_TYPE=mysql
      # - TU_DATABASE_HOST=mysql
      # - TU_DATABASE_PORT=3306
      # - TU_DATABASE_USER=tinyurl
      # - TU_DATABASE_PASSWORD=changeme
      # - TU_DATABASE_DBNAME=tinyurl
      # PostgreSQL:
      # - TU_DATABASE_TYPE=postgres
      # - TU_DATABASE_HOST=postgres
      # - TU_DATABASE_PORT=5432
      # - TU_DATABASE_USER=tinyurl
      # - TU_DATABASE_PASSWORD=changeme
      # - TU_DATABASE_DBNAME=tinyurl
    restart: unless-stopped

volumes:
  tinyurl-data:
```

### Multi-arch Build

```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  --build-arg VERSION=v0.1.0 \
  -t ghcr.io/wxccs/tinyurl:v0.1.0 \
  --push .
```

## License

[MIT](LICENSE)
