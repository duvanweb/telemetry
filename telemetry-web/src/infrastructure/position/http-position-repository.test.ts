import { describe, expect, it, vi } from "vitest";
import { HttpPositionRepository } from "@/infrastructure/position/http-position-repository";
import type { HttpClient } from "@/infrastructure/http/client";

describe("HttpPositionRepository", () => {
  function createMockClient(): { client: HttpClient; post: ReturnType<typeof vi.fn> } {
    const post = vi.fn().mockResolvedValue(null);
    const client = { get: vi.fn(), post } as unknown as HttpClient;
    return { client, post };
  }

  describe("sendPosition", () => {
    it("posts to the correct geo-service endpoint with lat/lng/timestamp", async () => {
      const { client, post } = createMockClient();
      const repo = new HttpPositionRepository(client);

      await repo.sendPosition(42, 4.71, -74.07, "2026-01-01T12:00:00Z");

      expect(post).toHaveBeenCalledWith(
        "/api/vehicles/42/positions",
        { lat: 4.71, lng: -74.07, timestamp: "2026-01-01T12:00:00Z" },
        undefined,
      );
    });

    it("forwards the abort signal", async () => {
      const { client, post } = createMockClient();
      const repo = new HttpPositionRepository(client);
      const controller = new AbortController();

      await repo.sendPosition(1, 0, 0, "2026-01-01T00:00:00Z", controller.signal);

      expect(post).toHaveBeenCalledWith(
        expect.any(String),
        expect.any(Object),
        controller.signal,
      );
    });
  });

  describe("sendMalformed", () => {
    it("posts invalid JSON body to trigger 400", async () => {
      const { client, post } = createMockClient();
      const repo = new HttpPositionRepository(client);

      await repo.sendMalformed(7);

      expect(post).toHaveBeenCalledWith("/api/vehicles/7/positions", '{"lat":}', undefined);
    });
  });
});
