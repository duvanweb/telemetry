// PositionRepository is the port for sending vehicle positions to geo-service.
// The application layer depends on this interface, never on a concrete infrastructure implementation.
export interface PositionRepository {
  // sendPosition posts a valid position report. Resolves on 202 Accepted.
  // Throws ApiError(409) on duplicate, ApiError(400) on invalid, ApiError(404) if vehicle not found.
  sendPosition(
    vehicleId: number,
    lat: number,
    lng: number,
    timestamp: string,
    signal?: AbortSignal,
  ): Promise<void>;

  // sendMalformed posts an intentionally invalid JSON body. Throws ApiError(400).
  sendMalformed(vehicleId: number, signal?: AbortSignal): Promise<void>;
}
