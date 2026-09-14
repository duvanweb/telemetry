import { geoRequest } from "./client";
import type { GpsPosition } from "./types";

// CreatePositionResponse mirrors the 202 Accepted body from geo-service.
interface CreatePositionResponse {
  status: string;
  vehicle_id: string;
}

// reportPosition sends a single GPS position to geo-service
// POST /api/vehicles/{vehicle_id}/positions (SPEC 02-geo).
// Returns the accepted response; throws ApiError on non-2xx (400/404/409/503).
export function reportPosition(
  vehicleId: number,
  position: GpsPosition,
): Promise<CreatePositionResponse> {
  return geoRequest<CreatePositionResponse>(
    `/api/vehicles/${vehicleId}/positions`,
    {
      method: "POST",
      body: JSON.stringify(position),
    },
  );
}
