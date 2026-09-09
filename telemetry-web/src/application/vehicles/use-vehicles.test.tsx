import { describe, expect, it, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import { RepositoryProvider } from "@/app/repository-context";
import { useVehicles } from "@/application/vehicles/use-vehicles";
import { useVehicleByPlate } from "@/application/vehicles/use-vehicle-by-plate";

const VEHICLE = {
  id: 1,
  plate: "ABC-123",
  createdAt: "2026-09-08T10:00:00Z",
  updatedAt: "2026-09-08T10:00:00Z",
};

function createWrapper(repository: VehicleRepository) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return function Wrapper({ children }: { children: ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <RepositoryProvider repository={repository}>
          {children}
        </RepositoryProvider>
      </QueryClientProvider>
    );
  };
}

describe("useVehicles", () => {
  it("returns paginated vehicles on success", async () => {
    const payload = { data: [VEHICLE], limit: 10, offset: 0, total: 1 };
    const mockRepo: VehicleRepository = {
      list: vi.fn().mockResolvedValue(payload),
      findByPlate: vi.fn(),
    };

    const { result } = renderHook(() => useVehicles({ limit: 10, offset: 0 }), {
      wrapper: createWrapper(mockRepo),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(payload);
    expect(mockRepo.list).toHaveBeenCalledWith(
      expect.objectContaining({ limit: 10, offset: 0 }),
    );
  });

  it("exposes error when repository fails", async () => {
    const mockRepo: VehicleRepository = {
      list: vi.fn().mockRejectedValue(new Error("fail")),
      findByPlate: vi.fn(),
    };

    const { result } = renderHook(() => useVehicles({ limit: 10, offset: 0 }), {
      wrapper: createWrapper(mockRepo),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useVehicleByPlate", () => {
  it("returns vehicle on success", async () => {
    const mockRepo: VehicleRepository = {
      list: vi.fn(),
      findByPlate: vi.fn().mockResolvedValue(VEHICLE),
    };

    const { result } = renderHook(() => useVehicleByPlate("ABC-123"), {
      wrapper: createWrapper(mockRepo),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(VEHICLE);
  });

  it("does not fetch when plate is empty", () => {
    const mockRepo: VehicleRepository = {
      list: vi.fn(),
      findByPlate: vi.fn(),
    };

    const { result } = renderHook(() => useVehicleByPlate(""), {
      wrapper: createWrapper(mockRepo),
    });

    expect(result.current.fetchStatus).toBe("idle");
    expect(mockRepo.findByPlate).not.toHaveBeenCalled();
  });
});
