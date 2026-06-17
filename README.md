# URL Shortener

A lightweight URL shortener built with Go, PostgreSQL, and Redis. Shorten long URLs, redirect with caching, and track click analytics.

## Features

- Shorten URLs with auto-generated short codes
- Redirect to original URLs with Redis caching for fast lookups
- Click analytics (IP, user agent, referer, timestamp)
- Configurable URL expiration
- Retry logic for short code collisions
- Health check endpoint

## Tech Stack

- **Go 1.24** — HTTP server with [chi](https://github.com/go-chi/chi) router
- **PostgreSQL 16** — Persistent storage
- **Redis 7** — In-memory cache for fast redirects
- **Docker Compose** — Container orchestration

## Quick Start

```bash
git clone <repo-url>
cd url-shortener
docker compose up -d
```

The server starts on `http://localhost:8080`.

### Apply database migrations

```bash
docker exec -i url-shortener-postgres psql -U admin -d shortener < migrations/001_create_urls_table.sql
docker exec -i url-shortener-postgres psql -U admin -d shortener < migrations/002_create_clicks_table.sql
```

## API Endpoints

### `POST /shorten`

Create a shortened URL.

**Request:**
```json
{
  "url": "https://example.com/very/long/url",
  "expires_at": "2026-12-31T23:59:59Z"
}
```

`expires_at` is optional (RFC 3339 format).

**Response:** `200 OK`
```json
{
  "short_code": "aB3xK9",
  "short_url": "http://localhost:8080/aB3xK9"
}
```

### `GET /{shortCode}`

Redirect to the original URL.

- Redirects with HTTP 302 if found.
- Returns `404` if the code doesn't exist or has expired.

### `GET /analytics/{shortCode}`

Get click analytics for a short URL.

**Response:** `200 OK`
```json
{
  "short_code": "aB3xK9",
  "expires_at": "2026-12-31T23:59:59Z",
  "total_clicks": 5,
  "clicks": [
    {
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0 ...",
      "referer": "https://example.com",
      "clicked_at": "2026-06-17T12:00:00Z"
    }
  ]
}
```

Returns `404` if the short code doesn't exist.

### `GET /health`

Health check. Returns `200 OK`.

## Environment Variables

| Variable      | Default        | Description          |
|---------------|----------------|----------------------|
| `DB_HOST`     | `postgres`     | PostgreSQL host      |
| `DB_PORT`     | `5432`         | PostgreSQL port      |
| `DB_USER`     | `admin`        | PostgreSQL user      |
| `DB_PASSWORD` | `admin`        | PostgreSQL password  |
| `DB_NAME`     | `shortener`    | Database name        |
| `DB_SSLMODE`  | `disable`      | SSL mode             |
| `REDIS_ADDR`  | `redis:6379`   | Redis address        |

## Run Locally (without Docker)

```bash
# Start PostgreSQL and Redis on localhost, then:
go run ./cmd/api
```

## Project Structure

```
├── cmd/api/main.go           — Application entry point
├── internal/
│   ├── handler/              — HTTP handlers
│   ├── service/              — Business logic
│   ├── repository/           — Database access
│   └── model/                — Data models
├── migrations/               — SQL migration files
├── pkg/
│   ├── code/                 — Short code generator
│   └── db/                   — Database clients (PostgreSQL, Redis)
├── docker-compose.yml        — Container setup
└── Dockerfile                — App image
```
