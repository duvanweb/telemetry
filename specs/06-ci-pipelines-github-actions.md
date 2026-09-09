# SPEC 06 — Pipelines básicos de GitHub Actions por servicio

> **Status:** Approved
> **Depends on:** SPEC 01 (telemetry-repo-bootstrap)
> **Date:** 2026-09-08
> **Objective:** Crear pipelines básicos de GitHub Actions —un workflow por servicio— que ejecuten lint, test y build para los tres microservicios Go (vehicle-service, geo-service, alert-service) y los dos frontends (telemetry-web, telemetry-movil) en cada push a `main`/`develop` y en cada PR a `develop`.

## Scope

**In:**

- Directorio `.github/workflows/` con **5 workflows**, uno por servicio:
  - `vehicle-service.yml` — Go: `go vet`, `golangci-lint`, `go test -race ./... -count=1`, `go build ./...`.
  - `geo-service.yml` — Go: idéntico al de vehicle-service.
  - `alert-service.yml` — Go: idéntico al patrón; `go test` pasa sin test files.
  - `telemetry-web.yml` — Node: `npm ci`, `tsc --noEmit` (typecheck), `npm run lint` (oxlint), `npm run build` (vite build). Sin `test` (no hay script todavía — step condicional).
  - `telemetry-movil.yml` — Node: `npm ci`, `npm run typecheck` (tsc --noEmit), `npm test` (jest). Sin `lint` (no hay script). Sin build (EAS Build no es "básico").
- Trigger uniforme: `push` en `main` y `develop`, `pull_request` en `develop`. Sin path filters (todos los workflows corren en cada push/PR).
- Versiones pinadas: Go 1.26 (`actions/setup-go@v5`), Node 20 LTS (`actions/setup-node@v4`).
- Caching: `setup-go` con cache de módulos Go; `setup-node` con `cache: 'npm'` apuntando al `package-lock.json` de cada frontend.
- `golangci-lint` vía `golangci/golangci-lint-action@v6` con config default (no existe `.golangci.yml` en los servicios hoy; se respeta si se añade después).
- `working-directory` seteado por workflow para que cada job corra dentro del directorio de su servicio.
- Checkout con `actions/checkout@v4`.

**Out of scope (for future specs):**

- Deploy / CD (publicar binaries, Docker images, EAS Build, Vercel/Netlify).
- Service containers en CI (PostgreSQL, Redis, Kafka) para tests de integración — los tests actuales usan go-sqlmock y jest mocks.
- Path filters por servicio (optimización de minutos de CI; diferida para mantener simplicidad).
- Matrix builds (múltiples versiones de Go/Node, múltiples OS).
- Coverage reports, SonarQube, CodeClimate, u otras integraciones de calidad.
- Required status checks / branch protection rules (configuración de GitHub, no de código).
- GitLab CI (el repo vive en GitHub; se reconsideraría si se migra o se añade un mirror).
- Lint para telemetry-movil (no tiene script de lint configurado hoy).
- Tests E2E en CI (Playwright/Cypress/Maestro).
- Notifications (Slack, email) on failure.
- Reusable workflows / composite actions para DRY (los 3 Go workflows se repiten; extraer un reusable workflow es una optimización futura).

## Data model

This feature introduces no new data structures. It adds CI configuration files (YAML) under `.github/workflows/` that orchestrate existing build, lint, and test commands.

## Implementation plan

1. **Workflow de vehicle-service.** Crear `.github/workflows/vehicle-service.yml`: trigger `push` en `main`/`develop` + `pull_request` en `develop`; job `build` en `ubuntu-latest` con `actions/checkout@v4`, `actions/setup-go@v5` (`go-version: '1.26'`, `cache: true`), `working-directory: vehicle-service`; steps `go vet ./...`, `golangci/golangci-lint-action@v6` (`working-directory: vehicle-service`), `go test -race ./... -count=1`, `go build ./...`. Commit `ci(vehicle-service): add github actions workflow`. Manual: el YAML es válido (lint con `actionlint` o revisión visual).

2. **Workflow de geo-service.** Crear `.github/workflows/geo-service.yml` con el mismo patrón que vehicle-service cambiando `working-directory: geo-service`. Commit `ci(geo-service): add github actions workflow`. Manual: YAML válido.

3. **Workflow de alert-service.** Crear `.github/workflows/alert-service.yml` con el mismo patrón; `go test -race ./... -count=1` pasa sin test files (Go reporta `no test files` y el step termina 0). Commit `ci(alert-service): add github actions workflow`. Manual: YAML válido.

4. **Workflow de telemetry-web.** Crear `.github/workflows/telemetry-web.yml`: `actions/setup-node@v4` (`node-version: '20'`, `cache: 'npm'`, `cache-dependency-path: telemetry-web/package-lock.json`); `working-directory: telemetry-web`; steps `npm ci`, `npx tsc --noEmit` (typecheck), `npm run lint` (oxlint), `npm run build` (tsc -b && vite build); step de test condicional (`if: npm test script exists` — se omite si no hay script `test`). Commit `ci(telemetry-web): add github actions workflow`. Manual: YAML válido.

5. **Workflow de telemetry-movil.** Crear `.github/workflows/telemetry-movil.yml`: `actions/setup-node@v4` (`node-version: '20'`, `cache: 'npm'`, `cache-dependency-path: telemetry-movil/package-lock.json`); `working-directory: telemetry-movil`; steps `npm ci`, `npm run typecheck` (tsc --noEmit), `npm test` (jest). Sin lint (no hay script). Sin build (EAS no es básico). Commit `ci(telemetry-movil): add github actions workflow`. Manual: YAML válido.

6. **Verificación en GitHub.** Push de la rama a origin, abrir PR a `develop`; verificar en la tab "Actions" que los 5 workflows se disparan y terminan en verde (o se omite el step de test de web si no hay script). Commit `ci: verify all workflows trigger on pr`.

## Acceptance criteria

- [ ] `.github/workflows/vehicle-service.yml` existe y define trigger `push` en `main`/`develop` + `pull_request` en `develop`.
- [ ] El workflow de vehicle-service ejecuta `go vet ./...`, `golangci-lint`, `go test -race ./... -count=1` y `go build ./...` con Go 1.26 en `ubuntu-latest`.
- [ ] `.github/workflows/geo-service.yml` existe con el mismo patrón que vehicle-service apuntando a `geo-service/`.
- [ ] `.github/workflows/alert-service.yml` existe con el mismo patrón; `go test` no falla por ausencia de tests.
- [ ] `.github/workflows/telemetry-web.yml` existe y ejecuta `npm ci`, `tsc --noEmit`, `npm run lint` (oxlint) y `npm run build` (vite build) con Node 20.
- [ ] El step de test de telemetry-web es condicional y se omite sin error si no existe script `test`.
- [ ] `.github/workflows/telemetry-movil.yml` existe y ejecuta `npm ci`, `npm run typecheck` (tsc --noEmit) y `npm test` (jest) con Node 20.
- [ ] Todos los workflows usan `actions/checkout@v4`.
- [ ] Los workflows Go usan `actions/setup-go@v5` con `go-version: '1.26'` y cache de módulos.
- [ ] Los workflows frontend usan `actions/setup-node@v4` con `node-version: '20'` y `cache: 'npm'`.
- [ ] Al abrir un PR a `develop`, los 5 workflows se disparan automáticamente en GitHub Actions.
- [ ] Los workflows que tienen gates existentes (vehicle-service, geo-service, telemetry-movil) terminan en verde.
- [ ] Todo el contenido nuevo (nombres de jobs, steps, comentarios) en inglés.

## Decisions

- **Sí:** GitHub Actions. El repo vive en `github.com/duvanweb/telemetry`; es nativo, sin configuración extra. **No:** GitLab CI — requeriría un runner self-hosted o un mirror; se mencionó pero el repo está en GitHub.
- **Sí:** Un workflow por servicio (5 archivos). Falla uno sin bloquear los demás; paralelismo total; claro ownership. **No:** un workflow por lenguaje (acopla servicios del mismo lenguaje) ni monolítico (acopla todo).
- **Sí:** Trigger `push` en `main`/`develop` + `pull_request` en `develop`. Casa con Git Flow (CLAUDE.md): `develop` es integración, `main` es producción. **No:** trigger en todas las ramas — gastaría minutos en ramas efímeras.
- **Sí:** Sin path filters. Todos los workflows corren en cada push/PR. Simplicidad: un cambio en cualquier lugar valida todo. **No:** path filters por servicio — optimización de minutos diferida; añade complejidad y el gotcha de required checks que no se disparan.
- **Sí:** Go 1.26 pinado (CLAUDE.md). Node 20 LTS pinado (CLAUDE.md). **No:** matrix de versiones — overkill para pipeline básico.
- **Sí:** `golangci-lint` con config default (sin `.golangci.yml`). CLAUDE.md dice "respetar .golangci.yml" — si se añade después, la action la detecta automáticamente. **No:** añadir `.golangci.yml` ahora — no existe en ningún servicio y crear una config custom es decisión de cada servicio.
- **Sí:** `go test -race ./... -count=1` (CLAUDE.md). El flag `-count=1` desactiva cache de tests; `-race` activa el race detector. **No:** sin `-race` — CLAUDE.md lo exige.
- **Sí:** `go build ./...` para validar compilación de todos los paquetes. **No:** `go build ./cmd/api` — algunos servicios podrían no tener `cmd/api` o tener múltiples entrypoints; `./...` es más robusto.
- **Sí:** Caching de módulos Go (`setup-go` cache) y npm (`setup-node` cache). Acelera CI sin coste de configuración. **No:** sin cache — más lento sin beneficio.
- **Sí:** Step de test condicional en telemetry-web. No hay script `test` todavía; el step se omite sin error en lugar de fallar. **No:** omitir el step completamente — cuando se añadan tests (SPEC 05), el step ya existe y solo se activa.
- **Sí:** Sin build para telemetry-movil. Expo/EAS Build es pesado y no es "básico"; el typecheck + jest cubre el gate. **No:** `eas build` — requiere cuenta EAS, créditos y configuración de credentials.
- **Sí:** Sin deploy/CD. Esta spec es solo CI (validación en push/PR). **No:** publicar binaries o imágenes — va en otra spec cuando haya staging/prod.
- **Sí:** Sin service containers (Postgres/Redis/Kafka). Los tests actuales usan go-sqlmock y jest mocks; no necesitan dependencias reales. **No:** `services:` en el workflow — añade complejidad y tiempo; se reconsidera cuando haya tests de integración.
- **Sí:** Los 3 workflows Go se repiten (DRY violation intencional). Extraer un reusable workflow es una optimización futura; por ahora la copia es explícita y fácil de leer. **No:** reusable workflow composite — añade indirección sin mucho beneficio con 3 servicios.

## Risks

| Risk | Mitigation |
| --- | --- |
| `golangci-lint` falla en código existente (errcheck, unused, etc. sobre código que ignora errores) | Si falla, el CI lo expone — es el punto. Arreglar los issues o añadir `.golangci.yml` con los linters deseados. El pipeline es nuevo; si rompe, se ajusta la config. |
| `npm ci` falla si no existe `package-lock.json` | telemetry-movil tiene `package-lock.json` (instalado con npm). telemetry-web debe tenerlo; si no, usar `npm install` como fallback. |
| `tsc --noEmit` en telemetry-web puede diferir de `tsc -b` (project references) | El step de build (`npm run build` = `tsc -b && vite build`) también corre y valida el build completo. Si `tsc --noEmit` falla por project references, se ajusta a `tsc -b --noEmit`. |
| alert-service no tiene `cmd/api` o código compilable | `go build ./...` compila todos los paquetes; si no hay código Go más allá de `go.mod`, el step pasa sin hacer nada. `go test` reporta "no test files" y termina 0. |
| Los workflows corren en cada push a `develop` gastando minutos | Sin path filters, 5 workflows corren en cada push. Para un repo de este tamaño el coste es marginal; si escala, se añaden path filters. |

## What is **not** in this spec

- Deploy / CD (binaries, Docker images, EAS Build, Vercel/Netlify).
- Service containers en CI (PostgreSQL, Redis, Kafka).
- Path filters por servicio.
- Matrix builds (múltiples versiones de Go/Node, múltiples OS).
- Coverage reports e integraciones de calidad (SonarQube, CodeClimate).
- Required status checks y branch protection rules.
- GitLab CI.
- Lint para telemetry-movil.
- Tests E2E en CI.
- Notifications on failure.
- Reusable workflows / composite actions para DRY.

Cada uno de esos, si llega, va en su propio spec.
