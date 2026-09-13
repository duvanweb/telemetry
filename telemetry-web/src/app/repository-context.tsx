import { createContext, useContext } from "react";
import type { ReactNode } from "react";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import type { AlertRepository } from "@/core/ports/alert-repository";
import type { PositionRepository } from "@/core/ports/position-repository";

// Repositories holds all repository interfaces provided via React Context.
interface Repositories {
  vehicle: VehicleRepository;
  alert: AlertRepository;
  position: PositionRepository;
}

const RepositoryContext = createContext<Repositories | null>(null);

// RepositoryProvider injects VehicleRepository, AlertRepository and PositionRepository via React Context.
export function RepositoryProvider({
  vehicleRepository,
  alertRepository,
  positionRepository,
  children,
}: {
  vehicleRepository: VehicleRepository;
  alertRepository: AlertRepository;
  positionRepository: PositionRepository;
  children: ReactNode;
}) {
  return (
    <RepositoryContext.Provider value={{ vehicle: vehicleRepository, alert: alertRepository, position: positionRepository }}>
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

// usePositionRepository returns the injected PositionRepository. Throws if no provider is present.
export function usePositionRepository(): PositionRepository {
  const repositories = useContext(RepositoryContext);
  if (repositories === null) {
    throw new Error("usePositionRepository must be used within a RepositoryProvider");
  }
  return repositories.position;
}
