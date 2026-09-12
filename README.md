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

## Vehicle Deletion — Saga Pattern (Choreographed)

### Problem

When a vehicle is deleted via `DELETE /api/vehicles/{id}`, the vehicle-service performs a **soft delete** (`deleted_at = NOW()`) in its own PostgreSQL database. However, the geo-service and alert-service maintain their own databases and Redis caches with data referencing that `vehicle_id`. Without a cleanup mechanism, this data becomes **orphaned**:

| Data store | Orphaned data | TTL |
|---|---|---|
| geo-service PostgreSQL (`gps_positions`) | All historical positions for the vehicle | None |
| geo-service Redis (`geo:pos:{vid}:{lat}:{lng}`) | Anti-duplicate cache keys | Auto-expires (~60s) |
| alert-service PostgreSQL (`alerts`) | All historical alerts for the vehicle | None |
| alert-service Redis (`alert:vehicle:{vid}`) | Position tracker (last known position) | **No TTL** — persists forever |

### Solution: Choreographed Saga with Eventual Consistency

Vehicle deletion uses a **choreographed saga** — no central orchestrator. Each service reacts independently to a `vehicle.deleted` event published by vehicle-service to a RabbitMQ fanout exchange. Reliability is provided by RabbitMQ's dead-letter exchange (DLX) and retry queue mechanism.

```
DELETE /api/vehicles/{id}
         │
         ▼
┌─────────────────────┐
│   vehicle-service    │
│  1. SoftDelete (PG)  │
│  2. Publish event    │──┐
└─────────────────────┘  │ vehicle.deleted {vehicleId, deletedAt}
                         │
              ┌──────────┴──────────┐
              │  vehicle_events      │  (fanout exchange)
              └──────────┬──────────┘
                    ┌─────┴─────┐
                    ▼           ▼
┌──────────────────────┐  ┌──────────────────────┐
│    geo-service        │  │   alert-service       │
│  (worker consumer)    │  │  (api consumer)       │
│                       │  │                       │
│  DELETE FROM          │  │  DELETE FROM          │
│    gps_positions      │  │    alerts             │
│    WHERE vehicle_id   │  │    WHERE vehicle_id   │
│                       │  │                       │
│  (Redis cache auto-   │  │  Redis DEL            │
│   expires via TTL)    │  │    alert:vehicle:{id} │
└──────────────────────┘  └──────────────────────┘
```

### Step-by-step flow

1. Client sends `DELETE /api/vehicles/{id}` to vehicle-service.
2. vehicle-service soft-deletes the vehicle in PostgreSQL (synchronous local transaction).
3. vehicle-service publishes a `VehicleDeletedEvent{vehicleId, deletedAt}` to the `vehicle_events` fanout exchange (best-effort — if publishing fails, the error is logged but the request still succeeds).
4. vehicle-service returns `204 No Content` to the client.
5. **geo-service** (worker) consumes the event from the `geo_vehicle_deletions` queue, executes `DELETE FROM gps_positions WHERE vehicle_id = $1`, and acks the message. Redis anti-duplicate cache keys auto-expire via TTL.
6. **alert-service** (API) consumes the event from the `alert_vehicle_deletions` queue, executes `DELETE FROM alerts WHERE vehicle_id = $1` and `Redis DEL alert:vehicle:{vehicleId}`, then acks the message.

### RabbitMQ topology

```
Exchange: vehicle_events (fanout, durable)

Queue: geo_vehicle_deletions (durable, DLX → geo_vehicle_deletions_retry)
Queue: geo_vehicle_deletions_retry (durable, TTL=10s, DLX → geo_vehicle_deletions)

Queue: alert_vehicle_deletions (durable, DLX → alert_vehicle_deletions_retry)
Queue: alert_vehicle_deletions_retry (durable, TTL=10s, DLX → alert_vehicle_deletions)
```

Event payload (JSON):
```json
{"vehicleId": 42, "deletedAt": "2026-09-11T22:00:00Z"}
```

### Consistency guarantees

- **Synchronous:** The vehicle is immediately hidden from users (soft-delete makes it return 404).
- **Eventual:** Cleanup of geo-service and alert-service data happens asynchronously. During the brief window between soft-delete and cleanup, stale data may exist in geo/alert, but **no new data is generated** (geo-service validates vehicle existence via HTTP and rejects new positions for soft-deleted vehicles).
- **Idempotency:** All cleanup operations are idempotent — `DELETE WHERE vehicle_id = $1` and Redis `DEL` are safe to retry. If a message is redelivered (via DLX retry), re-executing is a no-op.

### Failure handling

| Failure | Behavior |
|---|---|
| RabbitMQ is down when publishing | Soft-delete still succeeds. Event is lost (best-effort). A transactional outbox pattern is recommended as a future improvement. |
| Consumer fails to process (e.g., DB error) | Message is nacked (`requeue=false`) → goes to retry queue via DLX → after TTL, dead-lettered back to main queue for another attempt. |
| Consumer is down | Message is durably stored in the queue and processed when the consumer restarts. |
| All retries exhausted | Message remains in the retry loop. Inspect via RabbitMQ management UI (`http://localhost:15672`) for manual intervention. |

### Manual verification

```bash
# 1. Create a vehicle
curl -X POST http://localhost:8090/api/vehicles -H "Content-Type: application/json" -d '{"plate":"ABC-123"}'

# 2. Ingest a position (triggers geo + alert data)
curl -X POST http://localhost:8081/api/positions -H "Content-Type: application/json" \
  -d '{"vehicleId":1,"latitude":4.7,"longitude":-74.1,"recordedAt":"2026-09-11T22:00:00Z"}'

# 3. Verify alerts exist
curl http://localhost:8092/api/alerts

# 4. Delete the vehicle
curl -X DELETE http://localhost:8090/api/vehicles/1

# 5. Wait for async cleanup
sleep 2

# 6. Verify alerts are gone
curl http://localhost:8092/api/alerts  # → {"data":[],"total":0}

# 7. Verify Redis tracker key is gone
docker exec telemetry-redis redis-cli EXISTS alert:vehicle:1  # → 0

# 8. Verify positions are gone
docker exec telemetry-postgres psql -U postgres -d geo_service \
  -c "SELECT COUNT(*) FROM gps_positions WHERE vehicle_id = 1"  # → 0
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
