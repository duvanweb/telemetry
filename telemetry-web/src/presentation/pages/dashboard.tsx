import { VehiclesSection } from "@/presentation/sections/vehicles-section";
import { AlertsSection } from "@/presentation/sections/alerts-section";
import { MapSection } from "@/presentation/sections/map-section";

// Dashboard — single page with three stacked sections.
export function Dashboard() {
  return (
    <div className="min-h-screen bg-background p-6">
      <header className="mb-6">
        <h1 className="text-2xl font-bold text-foreground">Telemetry Dashboard</h1>
      </header>
      <main className="mx-auto max-w-5xl space-y-6">
        <VehiclesSection />
        <AlertsSection />
        <MapSection />
      </main>
    </div>
  );
}
