// Base URL of the vehicle-service backend.
// Expo public env vars must be prefixed with EXPO_PUBLIC_ to be available client-side.
export const VEHICLE_SERVICE_URL =
  process.env.EXPO_PUBLIC_VEHICLE_SERVICE_URL ?? "http://localhost:8080";

// Base URL of the geo-service backend (SPEC 02-geo position ingestion endpoint).
export const GEO_SERVICE_URL =
  process.env.EXPO_PUBLIC_GEO_SERVICE_URL ?? "http://localhost:8081";

// Base URL of the alert-service backend (alerts list + SSE stream).
export const ALERT_SERVICE_URL =
  process.env.EXPO_PUBLIC_ALERT_SERVICE_URL ?? "http://localhost:8082";
