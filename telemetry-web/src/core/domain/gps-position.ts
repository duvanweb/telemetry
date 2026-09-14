// GpsPosition — a vehicle GPS position received via SSE from alert-service.
export interface GpsPosition {
  vehicleId: number;
  latitude: number;
  longitude: number;
  recordedAt: string; // ISO 8601
}
