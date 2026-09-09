import { describe, expect, it, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import type { AlertRepository } from "@/core/ports/alert-repository";
import type { ListAlertsResult } from "@/core/domain/pagination";
import { RepositoryProvider } from "@/app/repository-context";
import { useAlerts } from "@/application/alerts/use-alerts";

const ALERT = {
  id: 1,
  vehicleId: 1,
  type: "VEHICLE_STOPPED",
  latitude: 4.71,
  longitude: -74.07,
  detectedAt: "2026-09-09T10:00:00Z",
  createdAt: "2026-09-09T10:00:05Z",
};

const mockVehicleRepo: VehicleRepository = {
  list: vi.fn(),
  findByPlate: vi.fn(),
};

function createWrapper(alertRepository: AlertRepository) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return function Wrapper({ children }: { children: ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <RepositoryProvider vehicleRepository={mockVehicleRepo} alertRepository={alertRepository}>
          {children}
        </RepositoryProvider>
      </QueryClientProvider>
    );
  };
}

describe("useAlerts", () => {
  it("returns paginated alerts on success", async () => {
    const payload: ListAlertsResult = { data: [ALERT], limit: 20, offset: 0, total: 1 };
    const mockAlertRepo: AlertRepository = {
      list: vi.fn().mockResolvedValue(payload),
    };

    const { result } = renderHook(() => useAlerts({ limit: 20, offset: 0 }), {
      wrapper: createWrapper(mockAlertRepo),
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(payload);
    expect(mockAlertRepo.list).toHaveBeenCalledWith(
      expect.objectContaining({ limit: 20, offset: 0 }),
    );
  });

  it("exposes error when repository fails", async () => {
    const mockAlertRepo: AlertRepository = {
      list: vi.fn().mockRejectedValue(new Error("fail")),
    };

    const { result } = renderHook(() => useAlerts({ limit: 20, offset: 0 }), {
      wrapper: createWrapper(mockAlertRepo),
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});
