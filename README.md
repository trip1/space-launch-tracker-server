# space-launch-server

Scaffold for a high-performance Go server that will expose upcoming space launches and events from [thespacedevs.com](https://thespacedevs.com/) with Valkey-backed caching.

## Stack

- Go 1.25+
- chi router (`github.com/go-chi/chi/v5`)
- Valkey-compatible client (`github.com/redis/go-redis/v9`)

## Endpoints

- `GET /healthz`
- `GET /readyz`
- `GET /v1/launches/upcoming`
- `GET /v1/events/upcoming`

Both upcoming endpoints accept optional `?limit=<n>`.

Upcoming endpoints check Valkey first and only call SpaceDevs on cache miss.

## Run

```bash
make run
```

## Development

```bash
make air-install
make dev
```

## Build

```bash
make build
```

## Docker

```bash
make docker-build
make docker-export
```

`make docker-export` builds the image and streams it to `10.0.10.1` over SSH for `docker load`. Override with `SERVER_USER`, `SERVER_HOST`, `SERVER_PORT`, or `DOCKER_IMAGE`.

## Config

Environment variables:

- `HTTP_ADDR` (default `:8080`)
- `HTTP_READ_TIMEOUT` (default `5s`)
- `HTTP_READ_HEADER_TIMEOUT` (default `2s`)
- `HTTP_WRITE_TIMEOUT` (default `10s`)
- `HTTP_IDLE_TIMEOUT` (default `60s`)
- `HTTP_SHUTDOWN_TIMEOUT` (default `10s`)
- `HTTP_MAX_HEADER_BYTES` (default `1048576`)
- `VALKEY_ADDR` (default `127.0.0.1:6379`)
- `VALKEY_PASSWORD` (default empty)
- `VALKEY_DB` (default `0`)
- `VALKEY_DIAL_TIMEOUT` (default `2s`)
- `VALKEY_READ_TIMEOUT` (default `500ms`)
- `VALKEY_WRITE_TIMEOUT` (default `500ms`)
- `VALKEY_POOL_SIZE` (default `50`)
- `VALKEY_MIN_IDLE_CONNS` (default `10`)
- `VALKEY_MAX_RETRIES` (default `2`)
- `VALKEY_POOL_TIMEOUT` (default `1s`)
- `VALKEY_CONN_MAX_IDLE_TIME` (default `15m`)
- `VALKEY_CONN_MAX_LIFETIME` (default `2h`)
- `SPACEDEVS_BASE_URL` (default `https://ll.thespacedevs.com/2.3.0`)
- `SPACEDEVS_TIMEOUT` (default `8s`)
- `SPACEDEVS_CACHE_TTL` (default `60s`)
- `SPACEDEVS_USER_AGENT` (default `space-launch-server/0.1`)
- `SPACEDEVS_LAUNCHES_LIMIT` (default `25`)
- `SPACEDEVS_EVENTS_LIMIT` (default `25`)
