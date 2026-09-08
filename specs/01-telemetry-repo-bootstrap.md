# SPEC 01 — Bootstrap del monorepo telemetry con tres microservicios Go y /health

> **Status:** Aprobado
> **Depends on:** —
> **Date:** 2026-09-08
> **Objective:** Inicializar el repo `telemetry-platform/telemetry` con la estructura monorepo y el scaffold de vehicle-service, geo-service y alert-service, cada uno como módulo Go independiente con endpoint `/health` de liveness y recursos totalmente aislados, dejando el flujo Git Flow (`main`, `develop`, `feature/repo-init`) creado en local y empujado a GitHub.

## Scope

**In:**

- Repo git local inicializado en `G:\Cursos\tecnical-test` con `.gitignore` raíz.
- Namespace `telemetry/` con `README.md` de arranque (prerrequisitos y tabla de servicios/puertos).
- Tres microservicios Go (`vehicle-service`, `geo-service`, `alert-service`), cada uno con su propio `go.mod` (module path `github.com/telemetry-platform/<service>`), sin `go.work`.
- Estructura interna mínima por servicio bajo Arquitectura Limpia/Hexagonal + Uber FX: `cmd/api/`, `internal/core/`, `internal/infrastructure/api/controllers/health/`, `internal/infrastructure/pkg/env/`, `internal/infrastructure/pkg/logger/`.
- Endpoint `GET /health` por servicio que responde `200 {"status":"ok","service":"<service>"}` (liveness simple).
- Configuración de puerto por env (`HTTP_PORT`) con defaults: vehicle `8080`, geo `8081`, alert `8082`.
- Git Flow local: ramas `main` y `develop`, rama `feature/repo-init` con todo el scaffolding.
- Repo GitHub `telemetry-platform/telemetry` (public) creado vía MCP github, con `main` y `develop` empujadas y PR `feature/repo-init` → `develop` abierto.

**Out of scope (for future specs):**

- Frontends `telemetry-web` y `telemetry-movil`.
- Infraestructura: PostgreSQL, Redis, Kafka, Docker Compose y migraciones SQL.
- Endpoints de negocio (CRUD de vehículos, geolocalización, alertas) y lógica de dominio.
- Readiness checks reales (`/health/ready` con BD/Redis).
- `go.work` (Go workspace) en la raíz.
- CI/CD, linting automático, swagger y tests de integración.

## Data model

Este spec no introduce entidades de dominio ni persistencia. La única estructura de datos es el DTO de respuesta del health check, idéntico en los tres servicios:

```go
// health response DTO
type HealthResponse struct {
    Status  string `json:"status"`  // "ok"
    Service string `json:"service"` // "vehicle-service" | "geo-service" | "alert-service"
}
```

## Implementation plan

1. **Init repo + Git Flow base.** `git init` en `G:\Cursos\tecnical-test`; crear `.gitignore` raíz para Go (`bin/`, `*.exe`, `vendor/`, `*.log`, `*.env`, `*.local`). Commit inicial en `main` (`chore: init repo`). Crear `develop` desde `main`, crear `feature/repo-init` desde `develop` y posicionarse en ella. Manual: `git branch` muestra las tres ramas.
2. **Namespace telemetry/.** Crear `telemetry/README.md` con prerrequisitos (Go 1.26) y tabla de servicios/puertos. Commit `docs: add telemetry readme`.
3. **Scaffold vehicle-service.** Crear `telemetry/vehicle-service/` con `go.mod` (module `github.com/telemetry-platform/vehicle-service`, go 1.26), `cmd/api/main.go` (app Uber FX), `internal/infrastructure/pkg/env/env.go` (lee `HTTP_PORT` con default `8080`), `internal/infrastructure/pkg/logger/logger.go` (`log/slog`), `internal/infrastructure/api/controllers/health/health.go` (handler `GET /health`), router chi y módulo FX que ensambla todo. Commit `feat(vehicle-service): scaffold api with health check`. Manual: `cd telemetry/vehicle-service && go run ./cmd/api` → `GET http://localhost:8080/health` responde `200 {"status":"ok","service":"vehicle-service"}`.
4. **Scaffold geo-service.** Mismo patrón del paso 3 con module `github.com/telemetry-platform/geo-service`, default port `8081`, service name `geo-service`. Commit `feat(geo-service): scaffold api with health check`. Manual: `/health` responde en `8081`.
5. **Scaffold alert-service.** Mismo patrón del paso 3 con module `github.com/telemetry-platform/alert-service`, default port `8082`, service name `alert-service`. Commit `feat(alert-service): scaffold api with health check`. Manual: `/health` responde en `8082`.
6. **Verificar build y tests.** En cada servicio: `go build ./...` y `go test ./...` pasan sin errores; `go vet ./...` limpio. Sin tests nuevos obligatorios en esta fase, pero los comandos deben exit.
7. **Repo GitHub + PR.** Vía MCP github crear el repo `telemetry-platform/telemetry` (public). Añadir remoto `origin`, empujar `main`, `develop` y `feature/repo-init`. Abrir PR `feature/repo-init` → `develop` con resumen del scaffolding y los tres endpoints de health.

## Acceptance criteria

- [ ] `git init` ejecutado en `G:\Cursos\tecnical-test`; las ramas `main`, `develop` y `feature/repo-init` existen localmente.
- [ ] `telemetry/README.md` existe y documenta prerrequisitos (Go 1.26) y tabla de servicios/puertos.
- [ ] Cada servicio tiene su propio `go.mod` con module path `github.com/telemetry-platform/<service>`; no existe `go.work` en la raíz.
- [ ] `cd telemetry/vehicle-service && go build ./...` termina sin errores (ídem `geo-service` y `alert-service`).
- [ ] `GET http://localhost:8080/health` responde `200 {"status":"ok","service":"vehicle-service"}`.
- [ ] `GET http://localhost:8081/health` responde `200 {"status":"ok","service":"geo-service"}`.
- [ ] `GET http://localhost:8082/health` responde `200 {"status":"ok","service":"alert-service"}`.
- [ ] Cada servicio arranca con override de `HTTP_PORT` (p.ej. `HTTP_PORT=9090` → escucha en `9090`).
- [ ] Ningún servicio importa paquetes de otro servicio (aislamiento verificado: `go list -deps ./...` de cada servicio no referencia los otros dos módulos).
- [ ] Repo GitHub `telemetry-platform/telemetry` (public) existe y contiene las ramas `main` y `develop`.
- [ ] PR `feature/repo-init` → `develop` abierto en GitHub.
- [ ] Todos los commits siguen Conventional Commits (`chore`, `docs`, `feat`).

## Decisions

- **Sí:** `go.mod` independiente por servicio. Cumple el requisito de no compartir recursos; cada servicio builda y corre de forma aislada.
- **No:** `go.work` (workspace) en la raíz. El usuario lo descartó; cada servicio se construye desde su propio directorio.
- **Sí:** Liveness simple `GET /health` con JSON `{"status":"ok","service":"..."}`. No hay dependencias externas que chequear en esta fase.
- **No:** Endpoints `/health/live` y `/health/ready` separados. Readiness no tiene nada que verificar hasta que exista BD/Redis.
- **Sí:** Defaults de puerto por servicio (8080/8081/8082) + override vía `HTTP_PORT`. Permite correr los tres en paralelo sin colisión.
- **No:** Frontends `telemetry-web`/`telemetry-movil` en este spec. Se difirieron a specs posteriores.
- **No:** Infraestructura (PostgreSQL, Redis, Kafka, Docker Compose). Diferida; los servicios arrancan sin dependencias externas.
- **Sí:** Repo GitHub `telemetry-platform/telemetry` público. El usuario lo eligió; el module path usa `github.com/telemetry-platform/<service>` como namespace lógico del monorepo.
- **Sí:** `log/slog` (stdlib Go 1.26) para logging. Sin dependencia externa de logging en esta fase.
- **Sí:** Router `chi` y DI vía Uber FX. Son el estándar del proyecto definido en `CLAUDE.md`.

## Risks

| Risk                                              | Mititación                                                                                          |
| ------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| El owner `telemetry-platform` no existe en GitHub o no hay permisos para crear el repo | Verificar con `gh api orgs/telemetry-platform` (o `gh api users/telemetry-platform`) antes de crear; si falla, pedir al usuario el owner correcto. |
| Colisión de puertos si se corren varios servicios sin override | Defaults distintos (8080/8081/8082) y `HTTP_PORT` documentado en `telemetry/README.md`.               |
| Uber FX añade complejidad al scaffold inicial     | FX es el estándar del proyecto (`CLAUDE.md`); se usa el módulo FX mínimo necesario para arrancar.    |

## What is **not** in this spec

- Frontends `telemetry-web` y `telemetry-movil`.
- Infraestructura (PostgreSQL, Redis, Kafka, Docker Compose, migraciones SQL).
- Endpoints de negocio y lógica de dominio.
- Readiness checks con dependencias reales.
- `go.work` / Go workspace.
- CI/CD, swagger, linting automático y tests de integración.

Cada uno de esos, si llega, va en su propio spec.
