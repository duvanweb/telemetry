import { ApiError } from "./errors";

// HttpClient — minimal fetch wrapper that prefixes a base URL and throws ApiError on non-2xx.
export class HttpClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async get<T>(path: string, signal?: AbortSignal): Promise<T> {
    let response: Response;
    try {
      response = await fetch(`${this.baseUrl}${path}`, {
        signal,
        headers: { Accept: "application/json" },
      });
    } catch (cause) {
      // Re-throw abort errors so callers (e.g. TanStack Query) can handle cancellation.
      if (cause instanceof DOMException && cause.name === "AbortError") {
        throw cause;
      }
      throw new ApiError(0, "Network error, try again");
    }

    if (!response.ok) {
      throw new ApiError(response.status, `HTTP ${response.status}`);
    }

    return (await response.json()) as T;
  }

  // post — sends a POST request. If body is a string it is sent as-is (for raw/malformed payloads);
  // otherwise it is JSON.stringify'd. Throws ApiError on non-2xx, re-throws AbortError.
  async post<T>(path: string, body?: unknown, signal?: AbortSignal): Promise<T> {
    let response: Response;
    try {
      response = await fetch(`${this.baseUrl}${path}`, {
        method: "POST",
        signal,
        headers: { Accept: "application/json", "Content-Type": "application/json" },
        body: typeof body === "string" ? body : body !== undefined ? JSON.stringify(body) : undefined,
      });
    } catch (cause) {
      if (cause instanceof DOMException && cause.name === "AbortError") {
        throw cause;
      }
      throw new ApiError(0, "Network error, try again");
    }

    if (!response.ok) {
      throw new ApiError(response.status, `HTTP ${response.status}`);
    }

    const text = await response.text();
    return (text ? JSON.parse(text) : null) as T;
  }
}
