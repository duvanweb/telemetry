import { createContext, useContext } from "react";
import type { ReactNode } from "react";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import type { AlertRepository } from "@/core/ports/alert-repository";

// Repositories holds all repository interfaces provided via React Context.
interface Repositories {
  vehicle: VehicleRepository;
  alert: AlertRepository;
}

const RepositoryContext = createContext<Repositories | null>(null);

// RepositoryProvider injects VehicleRepository and AlertRepository via React Context.
export function RepositoryProvider({
  vehicleRepository,
  alertRepository,
  children,
}: {
  vehicleRepository: VehicleRepository;
  alertRepository: AlertRepository;
  children: ReactNode;
}) {
  return (
    <RepositoryContext.Provider value={{ vehicle: vehicleRepository, alert: alertRepository }}>
      {children}
    </RepositoryContext.Provider>
  );
}

// useVehicleRepository returns the injected VehicleRepository. Throws if no provider is present.
export function useVehicleRepository(): VehicleRepository {
  const repositories = useContext(RepositoryContext);
  if (repositories === null) {
    throw new Error("useVehicleRepository must be used within a RepositoryProvider");
  }
  return repositories.vehicle;
}

// useAlertRepository returns the injected AlertRepository. Throws if no provider is present.
export function useAlertRepository(): AlertRepository {
  const repositories = useContext(RepositoryContext);
  if (repositories === null) {
    throw new Error("useAlertRepository must be used within a RepositoryProvider");
  }
  return repositories.alert;
}
