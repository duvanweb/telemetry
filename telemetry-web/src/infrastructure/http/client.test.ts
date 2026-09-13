import { describe, expect, it, vi } from "vitest";
import { HttpClient } from "@/infrastructure/http/client";
import { ApiError } from "@/infrastructure/http/errors";

function mockResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    text: () => Promise.resolve(body === null ? "" : JSON.stringify(body)),
    json: () => Promise.resolve(body),
  } as unknown as Response;
}

describe("HttpClient.post", () => {
  it("sends POST with JSON body and returns parsed response", async () => {
    const fetchSpy = vi.fn().mockResolvedValue(mockResponse({ id: 1 }, 201));
    vi.stubGlobal("fetch", fetchSpy);

    const client = new HttpClient("http://api");
    const result = await client.post("/test", { name: "x" });

    expect(result).toEqual({ id: 1 });
    expect(fetchSpy).toHaveBeenCalledWith("http://api/test", expect.objectContaining({
      method: "POST",
      body: JSON.stringify({ name: "x" }),
    }));
  });

  it("sends raw string body as-is", async () => {
    const fetchSpy = vi.fn().mockResolvedValue(mockResponse(null, 202));
    vi.stubGlobal("fetch", fetchSpy);

    const client = new HttpClient("http://api");
    const result = await client.post("/test", '{"lat":}');

    expect(result).toBeNull();
    expect(fetchSpy).toHaveBeenCalledWith("http://api/test", expect.objectContaining({
      body: '{"lat":}',
    }));
  });

  it("throws ApiError with status on non-2xx", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(mockResponse({ message: "bad" }, 409)));

    const client = new HttpClient("http://api");
    await expect(client.post("/test", {})).rejects.toThrow(ApiError);
    await expect(client.post("/test", {})).rejects.toHaveProperty("status", 409);
  });

  it("throws ApiError(0) on network failure", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new TypeError("fetch failed")));

    const client = new HttpClient("http://api");
    await expect(client.post("/test", {})).rejects.toThrow(ApiError);
    await expect(client.post("/test", {})).rejects.toHaveProperty("status", 0);
  });

  it("re-throws AbortError", async () => {
    const abortError = new DOMException("aborted", "AbortError");
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(abortError));

    const client = new HttpClient("http://api");
    await expect(client.post("/test", {})).rejects.toBe(abortError);
  });
});
