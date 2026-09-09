// Environment configuration for vehicle-service.
export const VEHICLE_SERVICE_URL =
  import.meta.env.VITE_VEHICLE_SERVICE_URL ?? "http://localhost:8080";

// Environment configuration for alert-service.
export const ALERT_SERVICE_URL =
  import.meta.env.VITE_ALERT_SERVICE_URL ?? "http://localhost:8082";

// Page size for the vehicle list.
export const VEHICLES_PAGE_SIZE = 10;
