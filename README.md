# Telemetry Platform

Monorepo for vehicle telemetry: backend microservices in Go (Clean / Hexagonal Architecture + Uber FX) and frontend apps in React and React Native.

## Prerequisites

| Tool | Version |
|------|---------|
| Go   | 1.26    |

> Node, Docker, PostgreSQL, Redis and Kafka are not required for the current phase. They will be introduced in future specs.

## Services

Each microservice is an independent Go module with its own `go.mod`. Resources are not shared between services.

| Service          | Module path                                         | Default port | Health endpoint |
|------------------|-----------------------------------------------------|--------------|-----------------|
| vehicle-service  | `github.com/telemetry-platform/vehicle-service`     | `8080`       | `GET /health`   |
| geo-service      | `github.com/telemetry-platform/geo-service`         | `8081`       | `GET /health`   |
| alert-service    | `github.com/telemetry-platform/alert-service`       | `8082`       | `GET /health`   |

## Run a service

From the service directory:

```bash
cd telemetry/vehicle-service
go run ./cmd/api
```

Override the listen port with the `HTTP_PORT` environment variable:

```bash
HTTP_PORT=9090 go run ./cmd/api
```

## Health check

Each service exposes a liveness endpoint:

```bash
curl http://localhost:8080/health
# {"status":"ok","service":"vehicle-service"}
```

## Project structure

```
telemetry/
├── vehicle-service/   # Backend Go — vehicle management
├── geo-service/       # Backend Go — geolocation per vehicle
├── alert-service/     # Backend Go — alerts (geofences, speed, ...)
├── telemetry-web/     # Frontend React (pending)
└── telemetry-mobile/   # App React Native (pending)
```

## Pending (future specs)

- Frontend apps: `telemetry-web`, `telemetry-mobile`.
- Infrastructure: PostgreSQL, Redis, Kafka, Docker Compose.
- Business endpoints and domain logic.
- Readiness checks with real dependencies.
