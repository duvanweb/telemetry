import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { HttpClient } from "@/infrastructure/http/client";
import { ApiError } from "@/infrastructure/http/errors";
import { HttpVehicleRepository } from "@/infrastructure/vehicle/http-vehicle-repository";

const VEHICLE = {
  id: 1,
  plate: "ABC-123",
  createdAt: "2026-09-08T10:00:00Z",
  updatedAt: "2026-09-08T10:00:00Z",
};

function mockResponse(body: unknown, ok = true, status = 200): Response {
  return { ok, status, json: () => Promise.resolve(body) } as unknown as Response;
}

describe("HttpVehicleRepository", () => {
  let repository: HttpVehicleRepository;

  beforeEach(() => {
    repository = new HttpVehicleRepository(new HttpClient("http://localhost:8080"));
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("list", () => {
    it("returns paginated vehicles on success", async () => {
      const payload = { data: [VEHICLE], limit: 10, offset: 0, total: 1 };
      vi.stubGlobal("fetch", vi.fn().mockResolvedValue(mockResponse(payload)));

      const result = await repository.list({ limit: 10, offset: 0 });

      expect(result).toEqual(payload);
      expect(fetch).toHaveBeenCalledWith(
        "http://localhost:8080/api/vehicles?limit=10&offset=0",
        expect.objectContaining({ headers: { Accept: "application/json" } }),
      );
    });

    it("throws ApiError on non-2xx response", async () => {
      vi.stubGlobal("fetch", vi.fn().mockResolvedValue(mockResponse({}, false, 500)));

      try {
        await repository.list({ limit: 10, offset: 0 });
        expect.fail("Should have thrown");
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        expect((error as ApiError).status).toBe(500);
      }
    });
  });

  describe("findByPlate", () => {
    it("returns vehicle on success", async () => {
      vi.stubGlobal("fetch", vi.fn().mockResolvedValue(mockResponse(VEHICLE)));

      const result = await repository.findByPlate("ABC-123");

      expect(result).toEqual(VEHICLE);
      expect(fetch).toHaveBeenCalledWith(
        "http://localhost:8080/api/vehicles/plates/ABC-123",
        expect.objectContaining({ headers: { Accept: "application/json" } }),
      );
    });

    it("throws ApiError with status 404 when not found", async () => {
      vi.stubGlobal("fetch", vi.fn().mockResolvedValue(mockResponse({}, false, 404)));

      try {
        await repository.findByPlate("XYZ-999");
        expect.fail("Should have thrown");
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        expect((error as ApiError).status).toBe(404);
      }
    });

    it("throws ApiError with status 0 on network failure", async () => {
      vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new TypeError("Failed to fetch")));

      try {
        await repository.findByPlate("ABC-123");
        expect.fail("Should have thrown");
      } catch (error) {
        expect(error).toBeInstanceOf(ApiError);
        expect((error as ApiError).status).toBe(0);
        expect((error as ApiError).message).toBe("Network error, try again");
      }
    });
  });
});
