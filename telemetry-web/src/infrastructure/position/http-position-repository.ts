import type { PositionRepository } from "@/core/ports/position-repository";
import type { HttpClient } from "../http/client";

// HttpPositionRepository — implements PositionRepository against the geo-service HTTP API.
export class HttpPositionRepository implements PositionRepository {
  private http: HttpClient;

  constructor(http: HttpClient) {
    this.http = http;
  }

  async sendPosition(
    vehicleId: number,
    lat: number,
    lng: number,
    timestamp: string,
    signal?: AbortSignal,
  ): Promise<void> {
    await this.http.post(
      `/api/vehicles/${vehicleId}/positions`,
      { lat, lng, timestamp },
      signal,
    );
  }

  async sendMalformed(vehicleId: number, signal?: AbortSignal): Promise<void> {
    // Send intentionally invalid JSON to exercise geo-service error handling (expects 400).
    await this.http.post(`/api/vehicles/${vehicleId}/positions`, '{"lat":}', signal);
  }
}
