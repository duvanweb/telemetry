// Base URL of the vehicle-service backend.
// Expo public env vars must be prefixed with EXPO_PUBLIC_ to be available client-side.
export const VEHICLE_SERVICE_URL =
  process.env.EXPO_PUBLIC_VEHICLE_SERVICE_URL ?? "http://localhost:8080";
