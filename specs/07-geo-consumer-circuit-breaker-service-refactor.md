# SPEC 07 — Refactor del consumer de geo-service con servicio intermedio y Circuit Breaker

> **Status:** Approved
> **Depends on:** SPEC 02 (geo-gps-position-ingestion), SPEC 04 (alert-anomaly-detection-sse)
> **Date:** 2026-09-08
> **Objective:** Refactorizar el consumer del worker de geo-service para que persista posiciones a través de un `PositionProcessor` (servicio intermedio) en vez de llamar `PositionRepository` directamente, envolviendo la persistencia en un Circuit Breaker (`sony/gobreaker`) con retry queue vía DLX mutuo para que fallos temporales de PostgreSQL no causen pérdida de datos ni tight loops.

## Scope

**In:**

- Nuevo port `PositionProcessor` en `internal/core/ports/services/position_processor.go` con método `Process(ctx, domain.Position) error` + anotación mockery.
- Nuevo error de dominio `ErrCircuitBreakerOpen` en `internal/core/domain/errors.go`.
- Implementación `Processor` en `internal/core/position/processor.go`: envuelve `PositionRepository.Save` con `gobreaker.CircuitBreaker`. `Process` ejecuta `cb.Execute(func() error { return repo.Save(ctx, pos) })`. Si CB open → retorna `ErrCircuitBreakerOpen`.
- Circuit Breaker config: threshold 5 fallos consecutivos, timeout 30s. Configurable via `CB_FAILURE_THRESHOLD` (default 5) y `CB_TIMEOUT` (default 30s).
- Retry queue topology: cola `gps_positions_retry` con `x-message-ttl` configurable (`RETRY_QUEUE_TTL`, default 10s) y DLX mutuo entre `gps_positions` y `gps_positions_retry`.
  - `gps_positions` declara `x-dead-letter-exchange: ""` y `x-dead-letter-routing-key: gps_positions_retry`.
  - `gps_positions_retry` declara `x-message-ttl: 10000`, `x-dead-letter-exchange: ""` y `x-dead-letter-routing-key: gps_positions`.
- Refactor de `internal/infrastructure/rabbitmq/consumer.go`: reemplazar `repositories.PositionRepository` con `services.PositionProcessor`. `processMessage` deserializa y delega a `processor.Process`. Error → `Nack(false, false)` (mensaje va a DLX → `gps_positions_retry`). Nil → `Ack(false)`.
- Fix del bug existente en `consumer.go`: importación de `encoding/json` faltante (usa `json.Unmarshal` sin importar el paquete).
- Actualización de `DeclareTopology` en `client.go`: declarar `gps_positions` con DLX args y `gps_positions_retry` con TTL + DLX args.
- Wire FX: `rabbitmq/module.go` y `cmd/worker/module.go` actualizados para proveer `Processor` como `services.PositionProcessor` al consumer (en vez de `PositionRepository` directa).
- Env vars nuevas: `CB_FAILURE_THRESHOLD`, `CB_TIMEOUT`, `RETRY_QUEUE_TTL` en `Configuration`.
- Dependencia nueva: `github.com/sony/gobreaker` en `go.mod`.
- Tests unitarios del `Processor` con mocks (CB cerrado → save exitoso, CB open → ErrCircuitBreakerOpen sin llamar repo, CB half-open → un intento).
- Tests unitarios del `Consumer` refactorizado con mock de `PositionProcessor` (ack en éxito, nack en error, nack en CB open).

**Out of scope (for future specs):**

- Métricas/observabilidad del Circuit Breaker (Prometheus, tracing, eventos de CB).
- Backoff exponencial en la retry queue (TTL fijo de 10s).
- DLQ terminal para mensajes que excedan un máximo de reintentos.
- Circuit Breaker en el cliente HTTP de vehicle-service (`VehicleClient`).
- Circuit Breaker en el publisher de RabbitMQ.
- Modificaciones al alert-service (no se toca — el fanout exchange ya entrega a `alert_positions` independientemente del retry flow).
- Validación de rangos lat/lng en el consumer (la posición ya fue validada en el API-side `Service.Ingest` antes de publicar).
- Deduplicación en el consumer (ya manejada por el cache Redis en el API-side).

## Data model

### Nuevo port — PositionProcessor

```go
// internal/core/ports/services/position_processor.go
type PositionProcessor interface {
    Process(ctx context.Context, pos domain.Position) error
}
```

### Nuevo error de dominio

```go
// internal/core/domain/errors.go (adición)
var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
```

### Configuración (env — adición a `Configuration`)

```go
CBFailureThreshold int    `env:"CB_FAILURE_THRESHOLD" envDefault:"5"`
CBTimeout          string `env:"CB_TIMEOUT" envDefault:"30s"`
RetryQueueTTL      string `env:"RETRY_QUEUE_TTL" envDefault:"10s"`
```

### Topología RabbitMQ (modificación)

```
Cola: gps_positions (existente, modificada)
  durable: true
  x-dead-letter-exchange:    ""    (default exchange)
  x-dead-letter-routing-key: gps_positions_retry

Cola: gps_positions_retry (nueva)
  durable: true
  x-message-ttl:             10000  (RETRY_QUEUE_TTL, default 10s)
  x-dead-letter-exchange:    ""    (default exchange)
  x-dead-letter-routing-key: gps_positions
```

El fanout exchange `gps_positions_fanout` y el binding de `gps_positions` no cambian. La cola `alert_positions` (alert-service) no se toca.

## Implementation plan

1. **PositionProcessor port + domain error.** Crear `internal/core/ports/services/position_processor.go` con interfaz `PositionProcessor` + `//go:generate mockery`. Añadir `ErrCircuitBreakerOpen` a `internal/core/domain/errors.go` (crear el archivo si no existe). Generar mocks con `go generate ./...`. Commit `feat(geo-service): add position processor port and circuit breaker error`. Manual: `go build ./...` pasa.

2. **Processor implementation con Circuit Breaker.** Añadir `github.com/sony/gobreaker` a `go.mod`. Crear `internal/core/position/processor.go`: struct `Processor` con `repo PositionRepository`, `cb *gobreaker.CircuitBreaker`, `logger`. `NewProcessor` construye `gobreaker.CircuitBreaker` con `Settings{ReadyToTrip: contador de fallos consecutivos == threshold}`. `Process`: `cb.Execute(func() error { return repo.Save(ctx, pos) })`; si el error es `gobreaker.ErrOpenState` → retornar `domain.ErrCircuitBreakerOpen`. Crear `processor_test.go` con mocks: (a) CB cerrado → save exitoso → Process retorna nil, (b) CB open → Process retorna `ErrCircuitBreakerOpen` sin llamar `repo.Save`, (c) save falla 5 veces → CB abre en la 6ª. Commit `feat(geo-service): add position processor with circuit breaker`. Manual: `go test ./internal/core/position/... -count=1` pasa.

3. **Retry queue topology.** Modificar `DeclareTopology` en `client.go`: declarar `gps_positions` con DLX args (`x-dead-letter-exchange: ""`, `x-dead-letter-routing-key: gps_positions_retry`). Declarar `gps_positions_retry` con TTL + DLX args (`x-message-ttl` desde `RETRY_QUEUE_TTL`, `x-dead-letter-exchange: ""`, `x-dead-letter-routing-key: gps_positions`). Mantener el exchange fanout y el binding existentes. Commit `feat(geo-service): add retry queue with dlx for position persistence`. Manual: `go build ./...` pasa.

4. **Refactor del Consumer.** Modificar `consumer.go`: (a) añadir import de `encoding/json` (fix del bug existente), (b) reemplazar campo `repo repositories.PositionRepository` con `processor services.PositionProcessor`, (c) `processMessage` llama `c.processor.Process(ctx, pos)` en vez de `c.repo.Save(ctx, pos)`, (d) `NewConsumer` recibe `services.PositionProcessor` en vez de `repositories.PositionRepository`. Commit `refactor(geo-service): delegate consumer processing to position processor`. Manual: `go build ./cmd/worker` pasa.

5. **Wire FX + config.** Añadir `CBFailureThreshold`, `CBTimeout`, `RetryQueueTTL` a `Configuration` en `env.go`. Actualizar `cmd/worker/module.go`: proveer `position.NewProcessor` (anotado como `services.PositionProcessor`) en vez de pasar `PositionRepository` directa al consumer. `NewConsumer` ahora recibe `services.PositionProcessor`. Actualizar `rabbitmq/module.go` si es necesario. Commit `feat(geo-service): wire processor and retry config in worker`. Manual: `go build ./cmd/worker` pasa.

6. **Consumer tests.** Crear `consumer_test.go` con mock de `services.PositionProcessor`: (a) Process retorna nil → `Ack(false)`, (b) Process retorna error → `Nack(false, false)`, (c) Process retorna `ErrCircuitBreakerOpen` → `Nack(false, false)`, (d) unmarshal falla → `Nack(false, false)`. Ejecutar `go test -race ./... -count=1` 3 veces consecutivas. Commit `test(geo-service): add consumer tests with mocked processor`. Manual: tests pasan 3 veces sin flakiness.

7. **Verificación final.** `go build ./cmd/worker` y `go build ./cmd/api` compilan. `go vet ./...` limpio. `go test -race ./... -count=1` pasa 3 veces. `golangci-lint run ./...` limpio. Test de integración manual: levantar RabbitMQ + PostgreSQL, publicar posición, verificar persistencia; parar PostgreSQL, publicar → verificar que CB abre tras 5 fallos y mensajes van a `gps_positions_retry`; restaurar PostgreSQL → verificar que tras 30s el CB pasa a half-open y los mensajes se reintentan. Commit `test(geo-service): verify circuit breaker and retry queue integration`.

## Acceptance criteria

- [ ] `PositionProcessor` interface existe en `internal/core/ports/services/position_processor.go` con `Process(ctx, domain.Position) error` + mockery.
- [ ] `Processor` en `internal/core/position/processor.go` envuelve `PositionRepository.Save` con `gobreaker.CircuitBreaker`.
- [ ] `Processor.Process` retorna `domain.ErrCircuitBreakerOpen` cuando el CB está open, sin llamar `repo.Save`.
- [ ] CB abre tras 5 fallos consecutivos (`CB_FAILURE_THRESHOLD`, configurable via env, default 5).
- [ ] Tras 30s (`CB_TIMEOUT`, configurable via env, default 30s), CB pasa a half-open y permite un intento de prueba.
- [ ] `Consumer` usa `services.PositionProcessor` (no `repositories.PositionRepository` directamente).
- [ ] `consumer.go` importa `encoding/json` (bug fix).
- [ ] Consumer ackea cuando `Process` retorna nil.
- [ ] Consumer nackea con `requeue=false` cuando `Process` retorna error (incluyendo `ErrCircuitBreakerOpen`).
- [ ] `gps_positions` cola declara `x-dead-letter-routing-key=gps_positions_retry`.
- [ ] `gps_positions_retry` cola declara `x-message-ttl` desde `RETRY_QUEUE_TTL` (default 10000ms) y `x-dead-letter-routing-key=gps_positions`.
- [ ] El fanout exchange `gps_positions_fanout` y el binding de `gps_positions` no cambian.
- [ ] `cmd/worker` compila con `go build ./cmd/worker`.
- [ ] `go test -race ./... -count=1` pasa en geo-service, 3 veces consecutivas sin flakiness.
- [ ] `golangci-lint run` pasa sin errores en geo-service.
- [ ] `internal/core` no importa `internal/infrastructure` (regla de Arquitectura Limpia respetada).
- [ ] `CB_FAILURE_THRESHOLD`, `CB_TIMEOUT`, `RETRY_QUEUE_TTL` son configurables por env con defaults documentados.
- [ ] Todo el código nuevo (identificadores, comentarios, nombres de tests) en inglés.

## Decisions

- **Sí:** `sony/gobreaker` para el Circuit Breaker. Standard de facto en Go, sin dependencias pesadas, API simple. **No:** implementación custom — más código que mantener sin valor agregado. **No:** `failsafe-go` — más potente pero overkill para un solo CB.
- **Sí:** `PositionProcessor` como servicio intermedio entre consumer y repositorio. El consumer (infraestructura) no debe llamar el repositorio directamente — Clean Architecture exige que la lógica de negocio viva en `core` y que la infraestructura dependa de interfaces de `ports`. **No:** consumer → repo directo (estado actual) — viola la regla de dependencias.
- **Sí:** Retry queue con DLX mutuo. `Nack(requeue=false)` + `gps_positions_retry` con TTL 10s. Evita tight loop (a diferencia de `requeue=true`) y evita pérdida de datos (a diferencia de drop sin DLX). **No:** `Nack(requeue=true)` — tight loop cuando BD está caída. **No:** drop sin DLX — pierde datos.
- **Sí:** El retry flow no re-publica al fanout exchange. El DLX opera a nivel de cola: `gps_positions` → `gps_positions_retry` → `gps_positions`. Alert-service no recibe el mensaje de nuevo durante reintentos — ya lo procesó cuando se publicó originalmente. **No:** re-publicar al exchange durante retry — duplicaría la detección de anomalías en alert-service.
- **Sí:** 5 fallos / 30s timeout (defaults). Balance entre sensibilidad y tolerancia a blips transitorios. Mismo default que la spec 04 original (obsolete). Configurable via env.
- **Sí:** `RETRY_QUEUE_TTL` default 10s. Retraso corto para reintentos — si la BD se recupera rápido, los mensajes se procesan sin mucha demora. **No:** backoff exponencial — diferido a otra spec (complejidad adicional sin beneficio inmediato).
- **Sí:** Fix del bug de `json.Unmarshal` sin import en `consumer.go`. Es un bug pre-existente que impide que `cmd/worker` compile. Al refactorizar el consumer, se corrige naturalmente. **No:** ignorar el bug — el refactoring no compilaría sin el fix.
- **Sí:** Tests del Processor con CB (cerrado, open, half-open) y del Consumer con mock de Processor. Cubren la lógica nueva end-to-end. **No:** tests de integración con RabbitMQ real — diferidos (setup pesado, poco ROI en esta fase).
- **Sí:** Solo geo-service se modifica. Alert-service no se toca — el fanout exchange ya entrega a `alert_positions` independientemente del retry flow. **No:** modificar alert-service — fuera de scope.
- **Sí:** Sin validación de rangos ni deduplicación en el consumer. La posición ya fue validada y deduplicada en el API-side `Service.Ingest` antes de publicar al exchange. El consumer solo persiste. **No:** re-validar en el consumer — duplica lógica sin beneficio.

## Risks

| Risk | Mitigación |
| --- | --- |
| `gps_positions` ya existe en RabbitMQ sin DLX args; redeclarar con args diferentes falla con `PRECONDITION_FAILED` | En desarrollo, borrar la cola antes de recrear: `rabbitmqctl delete_queue gps_positions`. Documentar en README. En producción, migrar con cuidado (declarar cola nueva, mover mensajes, renombrar). |
| Mensaje envenenado (unmarshal siempre falla) entra en loop retry infinito | Aceptado. DLQ terminal con contador de reintentos va en otra spec. El TTL de 10s limita la velocidad del loop. |
| CB en half-open permite 1 petición; si falla, vuelve a open | Comportamiento esperado de `gobreaker`. Documentado en tests. |
| `RETRY_QUEUE_TTL` fijo (sin backoff) puede causar retry storm si la BD tarda mucho en recuperarse | Aceptado para esta spec. Backoff exponencial diferido a otra spec. Monitorear profundidad de `gps_positions_retry`. |
| `sony/gobreaker` cuenta fallos consecutivos, no fallos totales — un éxito resetea el contador | Comportamiento esperado y deseado: si la BD se recupera brevemente y luego vuelve a fallar, el CB necesita 5 fallos nuevos para abrir. |

## What is **not** in this spec

- Métricas y observabilidad del Circuit Breaker (Prometheus, tracing).
- Backoff exponencial en la retry queue.
- DLQ terminal con máximo de reintentos.
- Circuit Breaker en el cliente HTTP de vehicle-service.
- Circuit Breaker en el publisher de RabbitMQ.
- Modificaciones al alert-service.
- Validación de rangos y deduplicación en el consumer.
- Tests de integración con RabbitMQ real.

Cada uno de esos, si llega, va en su propio spec.
