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
