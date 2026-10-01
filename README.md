# URL Shortener

A lightweight URL shortener built with Go, PostgreSQL, and Redis. Shorten URLs, redirect with caching, track click analytics, and manage links from a web panel.

## Features

- Auto-generated short codes with configurable expiration
- Fast redirects via Redis caching
- Click analytics (IP, user agent, referer, timestamp)
- Web panel to search, view, and delete links
- Delete endpoint for link removal

## Tech Stack

- **Go 1.24** — HTTP server with [chi](https://github.com/go-chi/chi) router
- **PostgreSQL 16** — Persistent storage
- **Redis 7** — In-memory cache for fast redirects
- **Docker Compose** — Container orchestration (includes Adminer on port 8081)

## Quick Start

```bash
git clone <repo-url>
cd url-shortener
docker compose up -d
```

The server starts on `http://localhost:8080`, Adminer on `http://localhost:8081`.

### Apply database migrations

```bash
docker exec -i url-shortener-postgres psql -U admin -d shortener < migrations/001_create_urls_table.sql
docker exec -i url-shortener-postgres psql -U admin -d shortener < migrations/002_create_clicks_table.sql
```

## API Endpoints

### `POST /shorten`

Create a shortened URL.

```json
{
  "url": "https://example.com/very/long/url",
  "expires_at": "2026-12-31T23:59:59Z"
}
```

`expires_at` is optional (RFC 3339). Response `200`:

```json
{
  "short_code": "aB3xK9",
  "short_url": "http://localhost:8080/aB3xK9"
}
```

### `GET /{shortCode}`

Redirects with HTTP 302 to the original URL. Returns `404` if the code doesn't exist or has expired.

### `DELETE /{shortCode}`

Deletes the link and its cache entry. Returns `204 No Content`.

### `GET /analytics/{shortCode}`

Click analytics: total count plus the 50 most recent clicks. Returns `404` if the code doesn't exist.

### `GET /health`

Health check. Returns `200 OK`.

## Web Panel

A simple HTML interface at `/panel`:

| Route | Description |
|---|---|
| `GET /panel` | List all links with search (`q`) and pagination |
| `GET/POST /panel/create` | Create a new link |
| `GET /panel/{shortCode}` | Link detail with analytics and delete button |
| `POST /panel/{shortCode}/delete` | Delete a link |

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `postgres` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `admin` | PostgreSQL user |
| `DB_PASSWORD` | `admin` | PostgreSQL password |
| `DB_NAME` | `shortener` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `REDIS_ADDR` | `redis:6379` | Redis address |

## Run Locally (without Docker)

```bash
# Start PostgreSQL and Redis on localhost, then:
go run ./cmd/api
```
