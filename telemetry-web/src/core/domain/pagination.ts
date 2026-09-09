import type { Vehicle } from "./vehicle";

// ListVehiclesResult — paginated result from GET /api/vehicles.
export interface ListVehiclesResult {
  data: Vehicle[];
  limit: number;
  offset: number;
  total: number;
}
