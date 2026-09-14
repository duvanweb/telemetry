import type { ListAlertsResult } from "../domain/pagination";

// AlertRepository is the port for alert data access.
export interface AlertRepository {
  list(params: { limit: number; offset: number; signal?: AbortSignal }): Promise<ListAlertsResult>;
}
