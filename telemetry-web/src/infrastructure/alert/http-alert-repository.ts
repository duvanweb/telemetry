import type { ListAlertsResult } from "@/core/domain/pagination";
import type { AlertRepository } from "@/core/ports/alert-repository";
import type { HttpClient } from "../http/client";

// HttpAlertRepository — implements AlertRepository against the alert-service HTTP API.
export class HttpAlertRepository implements AlertRepository {
  private http: HttpClient;

  constructor(http: HttpClient) {
    this.http = http;
  }

  async list(params: { limit: number; offset: number; signal?: AbortSignal }): Promise<ListAlertsResult> {
    const query = new URLSearchParams({
      limit: String(params.limit),
      offset: String(params.offset),
    });
    return this.http.get<ListAlertsResult>(`/api/alerts?${query.toString()}`, params.signal);
  }
}
