import { alertRequest } from "./client";
import type { ListAlertsResponse } from "./types";

// listAlerts calls GET /api/alerts with pagination.
// Returns alerts ordered by detectedAt descending; throws ApiError on non-2xx.
export function listAlerts(limit = 20, offset = 0): Promise<ListAlertsResponse> {
  return alertRequest<ListAlertsResponse>(
    `/api/alerts?limit=${limit}&offset=${offset}`,
  );
}
