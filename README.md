# Telemetry Platform

Monorepo for vehicle telemetry: backend microservices in Go (Clean / Hexagonal Architecture + Uber FX) and frontend apps in React and React Native.

## Prerequisites

| Tool | Version           |
|------|-------------------|
| Go   | 1.26              |
| Node | 20 LTS (or 22)    |
| Expo | `npx expo@latest` |

> Docker, PostgreSQL, Redis and Kafka are needed only to run the backend services
> (see `docker-compose.yml`). The mobile app (`telemetry-movil`) needs Node + Expo only.

## Services

Each microservice is an independent Go module with its own `go.mod`. Resources are not shared between services.

| Service          | Module path                                         | Default port | Health endpoint |
|------------------|-----------------------------------------------------|--------------|-----------------|
| vehicle-service  | `github.com/telemetry-platform/vehicle-service`     | `8080`       | `GET /health`   |
| geo-service      | `github.com/telemetry-platform/geo-service`         | `8081`       | `GET /health`   |
| alert-service    | `github.com/telemetry-platform/alert-service`       | `8082`       | `GET /health`   |
| telemetry-movil  | `telemetry-movil/` (Expo + Expo Router, TypeScript) | n/a (client) | consumes `GET /health` of vehicle-service |

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
├── telemetry-movil/   # App React Native — vehicle register/search (Expo + Expo Router)
└── telemetry-web/     # Frontend React (pending)
```

## Run the mobile app (telemetry-movil)

Requires vehicle-service running (see above). From the repo root:

```bash
cd telemetry-movil
npx expo start
```

Point the app at the backend via env (the `EXPO_PUBLIC_` prefix is required by Expo for
client-side env):

| Target           | `EXPO_PUBLIC_VEHICLE_SERVICE_URL` |
|------------------|-----------------------------------|
| iOS Sim / Web    | `http://localhost:8080` (default) |
| Android Emulator | `http://10.0.2.2:8080`            |
| Physical device  | `http://<host-LAN-IP>:8080`       |

```bash
EXPO_PUBLIC_VEHICLE_SERVICE_URL=http://10.0.2.2:8080 npx expo start --android
```

## Pending (future specs)

- Frontend web app: `telemetry-web`.
- Infrastructure: PostgreSQL, Redis, Kafka, Docker Compose.
- Readiness checks with real dependencies.
