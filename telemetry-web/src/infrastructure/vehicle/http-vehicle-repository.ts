import type { Vehicle } from "@/core/domain/vehicle";
import type { ListVehiclesResult } from "@/core/domain/pagination";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import type { HttpClient } from "../http/client";

// HttpVehicleRepository — implements VehicleRepository against the vehicle-service HTTP API.
export class HttpVehicleRepository implements VehicleRepository {
  private http: HttpClient;

  constructor(http: HttpClient) {
    this.http = http;
  }

  async list(params: { limit: number; offset: number; signal?: AbortSignal }): Promise<ListVehiclesResult> {
    const query = new URLSearchParams({
      limit: String(params.limit),
      offset: String(params.offset),
    });
    return this.http.get<ListVehiclesResult>(`/api/vehicles?${query.toString()}`, params.signal);
  }

  async findByPlate(plate: string, signal?: AbortSignal): Promise<Vehicle> {
    return this.http.get<Vehicle>(
      `/api/vehicles/plates/${encodeURIComponent(plate)}`,
      signal,
    );
  }
}
