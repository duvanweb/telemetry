import { beforeEach, describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { QueryClient, QueryClientProvider, type UseQueryResult } from "@tanstack/react-query";
import type { ListVehiclesResult } from "@/core/domain/pagination";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import { RepositoryProvider } from "@/app/repository-context";
import { ApiError } from "@/infrastructure/http/errors";

vi.mock("@/application/vehicles/use-vehicles");
vi.mock("@/application/vehicles/use-vehicle-by-plate");

import { useVehicles } from "@/application/vehicles/use-vehicles";
import { useVehicleByPlate } from "@/application/vehicles/use-vehicle-by-plate";
import { VehiclesSection } from "@/presentation/sections/vehicles-section";

const VEHICLE = {
  id: 1,
  plate: "ABC-123",
  createdAt: "2026-09-08T10:00:00Z",
  updatedAt: "2026-09-08T10:00:00Z",
};

// Cast a partial object as a UseQueryResult (test-only helper).
function asResult<T>(partial: Record<string, unknown>): UseQueryResult<T> {
  return partial as unknown as UseQueryResult<T>;
}

function renderSection(initialPath = "/") {
  const mockRepo: VehicleRepository = {
    list: vi.fn(),
    findByPlate: vi.fn(),
  };
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <MemoryRouter initialEntries={[initialPath]}>
      <QueryClientProvider client={queryClient}>
        <RepositoryProvider repository={mockRepo}>
          <VehiclesSection />
        </RepositoryProvider>
      </QueryClientProvider>
    </MemoryRouter>,
  );
}

beforeEach(() => {
  vi.mocked(useVehicles).mockReturnValue(
    asResult<ListVehiclesResult>({
      isLoading: false,
      isError: false,
      data: undefined,
      refetch: vi.fn(),
    }),
  );
  vi.mocked(useVehicleByPlate).mockReturnValue(
    asResult({
      isLoading: false,
      isError: false,
      data: undefined,
      refetch: vi.fn(),
    }),
  );
});

describe("VehiclesSection", () => {
  it("renders the section title while loading", () => {
    vi.mocked(useVehicles).mockReturnValue(
      asResult<ListVehiclesResult>({ isLoading: true, data: undefined }),
    );
    renderSection();
    expect(screen.getByText("Vehicles")).toBeInTheDocument();
  });

  it("renders error message with retry button on failure", () => {
    vi.mocked(useVehicles).mockReturnValue(
      asResult<ListVehiclesResult>({
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

  it("renders empty state when no vehicles exist", () => {
    vi.mocked(useVehicles).mockReturnValue(
      asResult<ListVehiclesResult>({
        isLoading: false,
        isError: false,
        data: { data: [], limit: 10, offset: 0, total: 0 },
      }),
    );
    renderSection();
    expect(screen.getByText("No vehicles found")).toBeInTheDocument();
  });

  it("renders vehicle data and pagination in the table", () => {
    vi.mocked(useVehicles).mockReturnValue(
      asResult<ListVehiclesResult>({
        isLoading: false,
        isError: false,
        data: { data: [VEHICLE], limit: 10, offset: 0, total: 1 },
      }),
    );
    renderSection();
    expect(screen.getByText("ABC-123")).toBeInTheDocument();
    expect(screen.getByText(/Page 1 of 1/)).toBeInTheDocument();
  });

  it("shows inline error for invalid plate format without calling the API", () => {
    renderSection();

    const input = screen.getByPlaceholderText("Search by plate (ABC-123)");
    fireEvent.change(input, { target: { value: "abc-123" } });
    fireEvent.click(screen.getByRole("button", { name: "Search" }));

    expect(
      screen.getByText("Invalid plate format, expected ABC-123"),
    ).toBeInTheDocument();
  });
});
