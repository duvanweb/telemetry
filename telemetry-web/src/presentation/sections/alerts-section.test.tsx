import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider, type UseQueryResult } from "@tanstack/react-query";
import type { ListAlertsResult } from "@/core/domain/pagination";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import type { AlertRepository } from "@/core/ports/alert-repository";
import { RepositoryProvider } from "@/app/repository-context";
import { ApiError } from "@/infrastructure/http/errors";

vi.mock("@/application/alerts/use-alerts");

import { useAlerts } from "@/application/alerts/use-alerts";
import { AlertsSection } from "@/presentation/sections/alerts-section";

const ALERT = {
  id: 1,
  vehicleId: 1,
  type: "VEHICLE_STOPPED",
  latitude: 4.71,
  longitude: -74.07,
  detectedAt: "2026-09-09T10:00:00Z",
  createdAt: "2026-09-09T10:00:05Z",
};

// Cast a partial object as a UseQueryResult (test-only helper).
function asResult<T>(partial: Record<string, unknown>): UseQueryResult<T> {
  return partial as unknown as UseQueryResult<T>;
}

function renderSection() {
  const mockVehicleRepo: VehicleRepository = {
    list: vi.fn(),
    findByPlate: vi.fn(),
  };
  const mockAlertRepo: AlertRepository = {
    list: vi.fn(),
  };
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <RepositoryProvider vehicleRepository={mockVehicleRepo} alertRepository={mockAlertRepo}>
        <AlertsSection />
      </RepositoryProvider>
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  vi.mocked(useAlerts).mockReturnValue(
    asResult<ListAlertsResult>({
      isLoading: false,
      isError: false,
      data: undefined,
      refetch: vi.fn(),
    }),
  );
});

describe("AlertsSection", () => {
  it("renders the section title while loading", () => {
    vi.mocked(useAlerts).mockReturnValue(
      asResult<ListAlertsResult>({ isLoading: true, data: undefined }),
    );
    renderSection();
    expect(screen.getByText("Alerts")).toBeInTheDocument();
  });

  it("renders error message with retry button on failure", () => {
    vi.mocked(useAlerts).mockReturnValue(
      asResult<ListAlertsResult>({
        isLoading: false,
        isError: true,
        error: new ApiError(0, "Network error, try again"),
        refetch: vi.fn(),
      }),
    );
    renderSection();
    expect(screen.getByText("Network error, try again")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Retry" })).toBeInTheDocument();
  });

  it("renders empty state when no alerts exist", () => {
    vi.mocked(useAlerts).mockReturnValue(
      asResult<ListAlertsResult>({
        isLoading: false,
        isError: false,
        data: { data: [], limit: 20, offset: 0, total: 0 },
      }),
    );
    renderSection();
    expect(screen.getByText("No alerts found")).toBeInTheDocument();
  });

  it("renders alert data in the table", () => {
    vi.mocked(useAlerts).mockReturnValue(
      asResult<ListAlertsResult>({
        isLoading: false,
        isError: false,
        data: { data: [ALERT], limit: 20, offset: 0, total: 1 },
      }),
    );
    renderSection();
    expect(screen.getByText("VEHICLE_STOPPED")).toBeInTheDocument();
    expect(screen.getByText("4.7100")).toBeInTheDocument();
    expect(screen.getByText("-74.0700")).toBeInTheDocument();
  });
});
