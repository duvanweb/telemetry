import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { HttpClient } from "@/infrastructure/http/client";
import { ApiError } from "@/infrastructure/http/errors";
import { HttpAlertRepository } from "@/infrastructure/alert/http-alert-repository";

const ALERT = {
  id: 1,
  vehicleId: 1,
  type: "VEHICLE_STOPPED",
  latitude: 4.71,
  longitude: -74.07,
  detectedAt: "2026-09-09T10:00:00Z",
  createdAt: "2026-09-09T10:00:05Z",
};

function mockResponse(body: unknown, ok = true, status = 200): Response {
  return { ok, status, json: () => Promise.resolve(body) } as unknown as Response;
}

describe("HttpAlertRepository", () => {
  let repository: HttpAlertRepository;

  beforeEach(() => {
    repository = new HttpAlertRepository(new HttpClient("http://localhost:8082"));
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("list", () => {
    it("returns paginated alerts on success", async () => {
      const payload = { data: [ALERT], limit: 20, offset: 0, total: 1 };
      vi.stubGlobal("fetch", vi.fn().mockResolvedValue(mockResponse(payload)));

      const result = await repository.list({ limit: 20, offset: 0 });

      expect(result).toEqual(payload);
      expect(fetch).toHaveBeenCalledWith(
        "http://localhost:8082/api/alerts?limit=20&offset=0",
        expect.objectContaining({ headers: { Accept: "application/json" } }),
      );
    });

    it("throws ApiError on non-2xx response", async () => {
      vi.stubGlobal("fetch", vi.fn().mockResolvedValue(mockResponse({}, false, 500)));
      try {
        await repository.list({ limit: 20, offset: 0 });
        expect.fail("Should have thrown");
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        expect((error as ApiError).status).toBe(500);
      }
    });

    it("throws ApiError with status 0 on network failure", async () => {
      vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new TypeError("Failed to fetch")));
      try {
        await repository.list({ limit: 20, offset: 0 });
        expect.fail("Should have thrown");
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        expect((error as ApiError).status).toBe(0);
        expect((error as ApiError).message).toBe("Network error, try again");
      }
    });
  });
});
