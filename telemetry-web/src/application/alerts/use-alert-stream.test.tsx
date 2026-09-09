import { act } from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactNode } from "react";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import type { AlertRepository } from "@/core/ports/alert-repository";
import { RepositoryProvider } from "@/app/repository-context";
import { useAlertStream } from "@/application/alerts/use-alert-stream";

// MockEventSource — test double for the browser EventSource.
class MockEventSource {
  static instances: MockEventSource[] = [];
  listeners: Record<string, ((event: { data: string }) => void)[]> = {};
  closed = false;
  url: string;

  constructor(url: string) {
    this.url = url;
    MockEventSource.instances.push(this);
  }

  addEventListener(type: string, listener: (event: { data: string }) => void) {
    if (!this.listeners[type]) this.listeners[type] = [];
    this.listeners[type].push(listener);
  }

  close() {
    this.closed = true;
  }

  // Test helper: simulate receiving an event.
  emit(type: string, data: unknown) {
    const event = { data: JSON.stringify(data) };
    (this.listeners[type] ?? []).forEach((fn) => fn(event));
  }
}

const mockVehicleRepo: VehicleRepository = {
  list: vi.fn(),
  findByPlate: vi.fn(),
};

const mockAlertRepo: AlertRepository = {
  list: vi.fn(),
};

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return function Wrapper({ children }: { children: ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <RepositoryProvider vehicleRepository={mockVehicleRepo} alertRepository={mockAlertRepo}>
          {children}
        </RepositoryProvider>
      </QueryClientProvider>
    );
  };
}

describe("useAlertStream", () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.stubGlobal("EventSource", MockEventSource);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("updates positions on position event", async () => {
    const { result } = renderHook(() => useAlertStream([1]), {
      wrapper: createWrapper(),
    });

    const es = MockEventSource.instances[0];
    act(() => {
      es.emit("position", { vehicleId: 1, latitude: 4.71, longitude: -74.07, recordedAt: "2026-09-09T10:00:00Z" });
    });

    await waitFor(() => expect(result.current[1]).toBeDefined());
    expect(result.current[1]).toEqual({
      vehicleId: 1,
      latitude: 4.71,
      longitude: -74.07,
      recordedAt: "2026-09-09T10:00:00Z",
    });
  });

  it("closes EventSource on unmount", () => {
    const { unmount } = renderHook(() => useAlertStream([]), {
      wrapper: createWrapper(),
    });

    const es = MockEventSource.instances[0];
    expect(es.closed).toBe(false);

    unmount();

    expect(es.closed).toBe(true);
  });
});
