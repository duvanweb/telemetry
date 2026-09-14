import { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { VehiclesSection } from "@/presentation/sections/vehicles-section";
import { AlertsSection } from "@/presentation/sections/alerts-section";
import { MapSection } from "@/presentation/sections/map-section";
import { SimulationSection } from "@/presentation/sections/simulation-section";
import { useVehicles } from "@/application/vehicles/use-vehicles";
import { useAlertStream } from "@/application/alerts/use-alert-stream";
import { useSimulation } from "@/application/simulation/use-simulation";

// Dashboard — single page with simulation controls, vehicles, alerts, and map.
export function Dashboard() {
  const vehiclesQuery = useVehicles({ limit: 100, offset: 0 });
  const vehicleIds = vehiclesQuery.data?.data.map((v) => v.id) ?? [];
  const positions = useAlertStream(vehicleIds);
  const simulation = useSimulation();
  const queryClient = useQueryClient();

  // Invalidate vehicles list when the simulation starts (it creates 5 DB rows).
  useEffect(() => {
    if (simulation.status.running) {
      queryClient.invalidateQueries({ queryKey: ["vehicles"] });
    }
  }, [simulation.status.running, queryClient]);

  return (
    <div className="min-h-screen bg-background p-6">
      <header className="mb-6">
        <h1 className="text-2xl font-bold text-foreground">Telemetry Dashboard</h1>
      </header>
      <main className="mx-auto max-w-5xl space-y-6">
        <SimulationSection simulation={simulation} />
        <VehiclesSection simulation={simulation} />
        <AlertsSection />
        <MapSection positions={positions} />
      </main>
    </div>
  );
}
