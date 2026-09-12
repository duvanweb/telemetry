import { ALERT_SERVICE_URL, GEO_SERVICE_URL, VEHICLE_SERVICE_URL } from "@/config/env";

// ApiError carries the HTTP status alongside the backend message.
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// fetchJson is the shared HTTP wrapper: prefixes the base URL, parses the Go
// backend's {"message": string} error body on non-2xx, and types the success body.
async function fetchJson<T>(
  baseUrl: string,
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const res = await fetch(`${baseUrl}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...(init.headers ?? {}) },
  });

  if (!res.ok) {
    let message = `HTTP ${res.status}`;
    try {
      const body = await res.json();
      if (body?.message) message = body.message;
    } catch {
      // non-JSON body; keep the default status message
    }
    throw new ApiError(res.status, message);
  }

  if (res.status === 204) return undefined as T;

  return (await res.json()) as T;
}

// request calls vehicle-service (kept for backward compatibility with SPEC 03).
export async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  return fetchJson<T>(VEHICLE_SERVICE_URL, path, init);
}

// geoRequest calls geo-service (SPEC 02-geo position ingestion endpoint).
export async function geoRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  return fetchJson<T>(GEO_SERVICE_URL, path, init);
}

// alertRequest calls alert-service (alerts list + SSE stream).
export async function alertRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  return fetchJson<T>(ALERT_SERVICE_URL, path, init);
}
