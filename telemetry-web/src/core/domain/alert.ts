// Alert entity — represents an alert from alert-service.
export interface Alert {
  id: number;
  vehicleId: number;
  type: string;
  latitude: number;
  longitude: number;
  detectedAt: string; // ISO 8601
  createdAt: string;  // ISO 8601
}
