// Environment configuration for vehicle-service.
// Docker maps host port 8090 → container port 8080 (see docker-compose.yml).
export const VEHICLE_SERVICE_URL =
  import.meta.env.VITE_VEHICLE_SERVICE_URL ?? "http://localhost:8090";

// Environment configuration for alert-service.
// Docker maps host port 8092 → container port 8082 (see docker-compose.yml).
export const ALERT_SERVICE_URL =
  import.meta.env.VITE_ALERT_SERVICE_URL ?? "http://localhost:8092";

// Page size for the vehicle list.
export const VEHICLES_PAGE_SIZE = 10;
