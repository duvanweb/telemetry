import type { Vehicle } from "./vehicle";
import type { Alert } from "./alert";

// ListVehiclesResult — paginated result from GET /api/vehicles.
export interface ListVehiclesResult {
  data: Vehicle[];
  limit: number;
  offset: number;
  total: number;
}

// ListAlertsResult — paginated result from GET /api/alerts.
export interface ListAlertsResult {
  data: Alert[];
  limit: number;
  offset: number;
  total: number;
}
