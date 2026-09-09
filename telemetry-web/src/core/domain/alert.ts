// Alert entity — used by the alerts mockup section. Not wired to a backend service yet.
export interface Alert {
  id: number;
  type: string; // e.g. "SPEED", "GEOFENCE"
  vehiclePlate: string;
  message: string;
  timestamp: string; // ISO 8601
}
