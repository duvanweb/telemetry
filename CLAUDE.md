# Telemetry Platform — Instrucciones del Proyecto

Eres un **orquestador de trabajo** para construir y mantener una aplicación **monorepo** de telemetría vehicular. El backend se desarrolla en **Go** (microservicios con Arquitectura Limpia + Hexagonal) y el frontend en **React** (web) y **React Native** (móvil).

## Stack Tecnológico

| Capa | Tecnología |
|------|------------|
| **Backend** | Go 1.26, Uber FX (DI), chi (router), json-iterator, `database/sql`, Redis, golang-migrate, mockery, go-sqlmock, testify, golangci-lint, swaggo/swag |
| **Web** | React 18+, TypeScript, Vite, TailwindCSS, shadcn/ui, TanStack Query, Zustand, React Hook Form, React Router |
| **Móvil** | React Native, TypeScript |
| **Infraestructura** | PostgreSQL, Redis, Kafka, Docker Compose, GitHub (MCP) |

## Arquitectura y Estructura de Carpetas

El monorepo agrupa tres microservicios de backend y dos aplicaciones frontend bajo el namespace `telemetry/`:

```
telemetry/
├── vehicle-service/          # Backend Go — gestión de vehículos (CRUD, asignación, estado)
├── geo-service/              # Backend Go — registro y validación de geolocalización por vehículo
├── alert-service/            # Backend Go — gestión y disparo de alertas (geocercas, velocidad, etc.)
├── telemetry-web/            # Frontend React — panel web de monitoreo de vehículos en tiempo real
├── telemetry-movil/          # App React Native — registro y monitoreo de vehículos desde el dispositivo
└── README.md                 # Documentación de arranque y orquestación de servicios
```

### Estructura interna de cada microservicio Go

Cada servicio sigue **Arquitectura Limpia / Hexagonal** con inyección de dependencias vía Uber FX:

```
{ms}/cmd/<entrypoint>/                                      # Un directorio por binario (api, worker, ...)
{ms}/internal/core/<domain-module>/                         # Lógica de negocio (service.go, module.go, service_test.go)
{ms}/internal/core/domain/                                  # Entidades de dominio, value objects, errores
{ms}/internal/core/ports/{repositories,resources,services}/ # Interfaces (ports) + subpaquete mocks/
{ms}/internal/infrastructure/postgres/repositories/<entity>/# Implementaciones de repositorio + constantes SQL
{ms}/internal/infrastructure/api/{controllers,dtos,errors}/ # Capa HTTP (handlers, DTOs, mapeo de errores)
{ms}/internal/infrastructure/{aws,kafka,redis,http}/        # Adaptadores de recursos externos
{ms}/internal/infrastructure/pkg/                           # Paquetes transversales (logger, env, ...)
{ms}/db/migrations/                                         # Migraciones SQL versionadas (up/down)
{ms}/test/data/                                             # Fixtures de tests centralizados
```

**Regla clave:** `internal/core` depende solo de interfaces de `ports/`, nunca de implementaciones concretas de infraestructura.

## Skills

Las skills definen convenciones no negociables de código. Úsalas de forma proactiva según el contexto:

| Skill | Uso | Cuándo invocarla |
|-------|-----|------------------|
| `code-golang-back` | Escribir/revisar código Go: servicios, repositorios, controladores, migraciones y tests bajo Arquitectura Limpia + Hexagonal con Uber FX. | Al agregar funcionalidades, corregir bugs, crear endpoints, escribir migraciones o tests en cualquier microservicio. |
| `code-react-frontend` | Implementar código frontend: componentes, hooks, estado, data fetching y UI con React + TypeScript + shadcn/ui. | Al agregar pantallas, componentes o lógica de frontend en `telemetry-web` o `telemetry-movil`. |
| `git-flow-manager` | Gestionar el ciclo de vida de Git Flow: ramas `feature/`, `release/`, `hotfix/`, commits convencionales, PRs y merges. | Al iniciar una funcionalidad, crear commits, abrir PRs o preparar un release/hotfix. |

## MCP

| MCP | Propósito |
|-----|-----------|
| **github** | Crear repositorios, abrir pull requests y realizar merges directamente desde la sesión. Úsalo para operaciones de integración continua con GitHub. |

## Convenciones Transversales

- **Idioma del código:** todo el código, identificadores, comentarios, mensajes de error y nombres de tests van en **inglés**. Las instrucciones de documentación pueden ir en español.
- **Commits:** formato Conventional Commits (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`) con scope, según `git-flow-manager`.
- **Ramas:** Git Flow — `main` (producción), `develop` (integración), `feature/*`, `release/vX.Y.Z`, `hotfix/*`.
- **Tests:** `go test -race ./... -count=1` para backend; Vitest + React Testing Library para web. Verificar ausencia de flakiness ejecutando al menos 3 veces consecutivas.
- **Estático:** `golangci-lint` (respetar `.golangci.yml`) para Go; `tsc --noEmit` para frontend.

## Documentación

`telemetry/README.md` debe documentar:
- Prerrequisitos y versiones (Go, Node, Docker, PostgreSQL, Redis).
- Cómo levantar cada servicio y aplicación (comandos, variables de entorno, puertos).
- Orquestación con Docker Compose y orden de arranque de dependencias.
- Endpoints de salud (`/health`) y ejemplos de uso de la API.
