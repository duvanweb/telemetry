// Environment configuration for vehicle-service.
// Docker maps host port 8090 → container port 8080 (see docker-compose.yml).
export const VEHICLE_SERVICE_URL =
  import.meta.env.VITE_VEHICLE_SERVICE_URL ?? "http://localhost:8090";

// Environment configuration for alert-service.
// Docker maps host port 8092 → container port 8082 (see docker-compose.yml).
export const ALERT_SERVICE_URL =
  import.meta.env.VITE_ALERT_SERVICE_URL ?? "http://localhost:8092";

// Environment configuration for geo-service.
// Docker maps host port 8081 → container port 8081 (see docker-compose.yml).
// In Docker, this is empty (relative) so position POSTs go through the nginx proxy.
export const GEO_SERVICE_URL =
  import.meta.env.VITE_GEO_SERVICE_URL ?? "http://localhost:8081";

// Page size for the vehicle list.
export const VEHICLES_PAGE_SIZE = 10;
