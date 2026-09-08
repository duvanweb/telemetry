import { VEHICLE_SERVICE_URL } from "@/config/env";

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

// request is a thin fetch wrapper that prefixes the vehicle-service base URL,
// parses the Go backend's {"message": string} error body on non-2xx, and types the success body.
export async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(`${VEHICLE_SERVICE_URL}${path}`, {
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
