import type { Vehicle } from "../domain/vehicle";
import type { ListVehiclesResult } from "../domain/pagination";

// VehicleRepository is the port for vehicle data access.
// The application layer depends on this interface, never on a concrete infrastructure implementation.
export interface VehicleRepository {
  list(params: { limit: number; offset: number; signal?: AbortSignal }): Promise<ListVehiclesResult>;
  findByPlate(plate: string, signal?: AbortSignal): Promise<Vehicle>;
}
