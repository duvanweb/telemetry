// Vehicle mirrors the VehicleResponse DTO from vehicle-service (SPEC 02).
// Dates arrive as ISO 8601 strings; id is a number (Go int64 within JS safe range).
export interface Vehicle {
  id: number;
  plate: string;
  createdAt: string;
  updatedAt: string;
}

// GpsPosition mirrors the CreatePositionRequest body of geo-service
// POST /api/vehicles/{vehicle_id}/positions (SPEC 02-geo).
export interface GpsPosition {
  lat: number;
  lng: number;
  timestamp: string; // RFC3339
}

// ReportState tracks the telemetry reporting lifecycle shown in the status card.
export type ReportState = "idle" | "reporting" | "stopped" | "error";

// TrackingMode selects the position source: real device GPS or simulated route.
export type TrackingMode = "gps" | "simulation";

// Alert mirrors the AlertResponse DTO from alert-service.
// detectedAt/createdAt arrive as ISO 8601 strings; id/vehicleId are numbers.
export interface Alert {
  id: number;
  vehicleId: number;
  type: string;
  latitude: number;
  longitude: number;
  detectedAt: string;
  createdAt: string;
}

// ListAlertsResponse mirrors the paginated response from GET /api/alerts.
export interface ListAlertsResponse {
  data: Alert[];
  limit: number;
  offset: number;
  total: number;
}
