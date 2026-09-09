import type { ReactNode } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RepositoryProvider } from "@/app/repository-context";
import { HttpClient } from "@/infrastructure/http/client";
import { HttpVehicleRepository } from "@/infrastructure/vehicle/http-vehicle-repository";
import { HttpAlertRepository } from "@/infrastructure/alert/http-alert-repository";
import { VEHICLE_SERVICE_URL, ALERT_SERVICE_URL } from "@/infrastructure/config/env";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

// Composition root — wires concrete repository implementations into the RepositoryProvider.
const vehicleRepository = new HttpVehicleRepository(
  new HttpClient(VEHICLE_SERVICE_URL),
);

const alertRepository = new HttpAlertRepository(
  new HttpClient(ALERT_SERVICE_URL),
);

// Providers — wraps the app with QueryClientProvider and RepositoryProvider.
export function Providers({ children }: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <RepositoryProvider vehicleRepository={vehicleRepository} alertRepository={alertRepository}>
        {children}
      </RepositoryProvider>
    </QueryClientProvider>
  );
}
