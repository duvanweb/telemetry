import { createContext, useContext } from "react";
import type { ReactNode } from "react";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";

const RepositoryContext = createContext<VehicleRepository | null>(null);

// RepositoryProvider injects a VehicleRepository implementation via React Context.
export function RepositoryProvider({
  repository,
  children,
}: {
  repository: VehicleRepository;
  children: ReactNode;
}) {
  return (
    <RepositoryContext.Provider value={repository}>
      {children}
    </RepositoryContext.Provider>
  );
}

// useVehicleRepository returns the injected VehicleRepository. Throws if no provider is present.
export function useVehicleRepository(): VehicleRepository {
  const repository = useContext(RepositoryContext);
  if (repository === null) {
    throw new Error("useVehicleRepository must be used within a RepositoryProvider");
  }
  return repository;
}
