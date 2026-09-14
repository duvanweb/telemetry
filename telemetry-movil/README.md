# telemetry-movil

Aplicación móvil (Expo + React Native + TypeScript) para monitoreo vehicular en tiempo real. Permite registrar o buscar un vehículo y, desde la pantalla de monitoreo, reportar su posición GPS a `geo-service` en primer y segundo plano, con simulación de ruta, modo offline-first y estrategia de batería adaptiva.

## Prerrequisitos

- Node.js 20+
- Expo CLI (`npx expo`) — SDK 57
- EAS CLI (`npm i -g eas-cli`) para builds con código nativo
- `vehicle-service` corriendo (default `http://localhost:8080`)
- `geo-service` corriendo (default `http://localhost:8081`)
- Emulador Android/iOS o dispositivo físico

## Instalación y arranque

```bash
npm install
npx expo start
```

> **Importante:** la funcionalidad de **background location** (reporte en segundo plano) requiere un **custom dev client** — Expo Go no soporta `expo-task-manager` en background. Para desarrollo con background real:

```bash
# Crear un development build (requere EAS configurado)
eas build --profile development --platform android
# o local:
npx expo run:android
```

El modo simulación y el reporte en foreground funcionan en Expo Go.

## Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `EXPO_PUBLIC_VEHICLE_SERVICE_URL` | `http://localhost:8080` | URL de `vehicle-service` |
| `EXPO_PUBLIC_GEO_SERVICE_URL` | `http://localhost:8081` | URL de `geo-service` (endpoint de posiciones) |

En emulador Android, usar `http://10.0.2.2:8081` para apuntar al host. En dispositivo físico, la IP LAN del host.

## Permisos de ubicación

La app declara (vía plugin `expo-location` en `app.json`):

- `NSLocationWhenInUseUsageDescription` (iOS) — ubicación en primer plano.
- `NSLocationAlwaysAndWhenInUseUsageDescription` (iOS) — ubicación en background.
- `UIBackgroundModes: ["location"]` (iOS) — background location habilitado.
- `ACCESS_BACKGROUND_LOCATION` + `FOREGROUND_SERVICE_LOCATION` (Android).

Si el usuario deniega el permiso, el modo GPS real no arranca (la tarjeta muestra estado `Error`). El modo simulación no requiere permisos.

## Pantalla de monitoreo (home)

Tras registrar o buscar un vehículo por placa, la app vuelve al home que muestra la pantalla de monitoreo:

- **Tarjeta de estado** (solo lectura): placa, ID del vehículo, lat/lng actual y estado de reporte (`Idle` / `Reporting` / `Stopped` / `Error`).
- **Toggle de fuente** `Simulation` / `GPS real` — selecciona el origen de la posición.
- **Start / Stop** — arranca o detiene el reporte de telemetría.
- **Stop vehicle** — detiene el avance del vehículo y sigue reenviando la última posición conocida cada intervalo (señaliza vehículo quieto). Se reanuda con el mismo botón.
- **Offline mode** — toggle manual; cuando está activo, las posiciones se guardan en una cola local sin enviarse.
- **Panic** — botón placeholder (no-op); su implementación va en un spec futuro.

Si no hay vehículo activo, el home muestra un empty state con botones para registrar o buscar.

## Estrategia offline-first

**Problema:** el conductor puede perder conexión (p.ej. entrar a un túnel por 10 minutos) y las posiciones GPS no deben perderse.

**Solución implementada:**

1. **Cola local en AsyncStorage** (`telemetry-movil:monitoring:v1`): cuando el modo offline está activo, o cuando un envío a geo-service falla por red, la posición se encola en `monitoring-store.queue`. La cola persiste en AsyncStorage, así que sobrevive al cierre de la app.

2. **Sincronización con throttle al volver online:** al desactivar el modo offline (o al completar un envío exitoso), la cola se drena enviando las posiciones pendientes a `POST /api/vehicles/{id}/positions` **una a una** con un máximo de **2 peticiones concurrentes** (`SYNC_MAX_CONCURRENT`). Esto evita saturar el servidor con una ráfaga masiva al reconectar.

3. **Backoff exponencial:** si un envío falla durante el drenado, se reencola y se espera `min(500ms * 2^n, 30s)` antes de reintentar, donde `n` es el número de intentos consecutivos fallidos.

4. **Cota suave:** la cola tiene un máximo de 1000 posiciones; si se excede, se descartan las más antiguas para evitar crecimiento indefinido.

5. **Dedup tolerado:** `geo-service` responde `409 Conflict` para posiciones duplicadas dentro del TTL (60s). El reporter trata el 409 como éxito (el punto ya está registrado).

**Volumen esperado:** 10 minutos offline @ 5s de intervalo = ~120 posiciones. AsyncStorage maneja esto sin problema (≈ decenas de KB).

## Estrategia de batería adaptiva

**Problema:** leer el GPS cada segundo drena la batería rápidamente.

**Solución implementada** (combinación de adaptive sampling + distance filter):

1. **Intervalo adaptivo según batería** (vía `expo-battery`):
   - Batería ≥ 20% → intervalo de **5 segundos** (`DEFAULT_INTERVAL_MS`).
   - Batería < 20% → intervalo de **15 segundos** (`LOW_BATTERY_INTERVAL_MS`).
   - Nivel desconocido (-1) → intervalo por defecto (5s).

2. **Distance filter:** antes de enviar una posición, se calcula la distancia haversine respecto a la última posición enviada. Si el movimiento es **menor a 10 metros** (`MIN_DISTANCE_M`), la posición se omite. Esto evita reportar puntos redundantes en tráfico lento o detenido, reduciendo consumo de GPS, red y CPU.

3. **Background task del SO:** en segundo plano, `expo-task-manager` recibe actualizaciones de ubicación controladas por el sistema operativo, que ya aplica su propia gestión de batería (el SO puede matar la task en condiciones extremas). Al volver a foreground, el loop adaptivo retoma el control.

4. **Configuración de `watchPositionAsync`:** el GPS foreground usa `accuracy: Balanced`, `timeInterval: 5000`, `distanceInterval: 10`, delegando parte del filtrado al SO.

Los parámetros son constantes exportadas en `src/tracking/adaptive-sampler.ts` y pueden ajustarse sin cambiar la lógica.

## Arquitectura

```
telemetry-movil/
├── app/                    # Expo Router (pantallas)
│   ├── _layout.tsx         # Providers + registro de background task
│   ├── index.tsx           # Home = pantalla de monitoreo
│   ├── register.tsx        # Registro de vehículo (SPEC 03)
│   └── search.tsx          # Búsqueda por placa (SPEC 03)
├── src/
│   ├── api/                # Clientes HTTP (vehicle, geo)
│   ├── config/env.ts       # URLs de vehicle-service y geo-service
│   ├── store/              # Zustand: vehicle-store, monitoring-store
│   ├── tracking/           # Lógica de tracking
│   │   ├── route.ts        # Ruta predefinida Bogotá (15 waypoints)
│   │   ├── simulator.ts    # RouteSimulator (loop)
│   │   ├── gps-tracker.ts  # expo-location foreground
│   │   ├── background-task.ts # expo-task-manager background
│   │   ├── adaptive-sampler.ts # Batería adaptiva + distance filter
│   │   └── reporter.ts     # Envío + cola offline + sync throttle
│   └── lib/                # Utilidades (plate-schema)
```

## Tests

```bash
npm test          # jest-expo
npm run typecheck # tsc --noEmit
```

Cubren: `monitoring-store` (cola, hidratación, sync), `RouteSimulator` (loop, reset), `adaptive-sampler` (intervalo adaptivo, distance filter).
