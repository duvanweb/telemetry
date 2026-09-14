import { Button } from "@/presentation/ui/button";
import { Badge } from "@/presentation/ui/badge";
import { SectionCard } from "@/presentation/components/section-card";
import type { SimulationControls } from "@/application/simulation/use-simulation";

// SimulationSection — global simulation controls (start/stop) + live stats.
export function SimulationSection({ simulation }: { simulation: SimulationControls }) {
  const { status, starting, error, start, stop } = simulation;
  const { stats } = status;

  return (
    <SectionCard title="Simulation" demo>
      <div className="space-y-4">
        <div className="flex items-center gap-3">
          {status.running ? (
            <Button variant="destructive" onClick={stop}>
              Stop Simulation
            </Button>
          ) : (
            <Button onClick={() => void start()} disabled={starting}>
              {starting ? "Starting…" : "Start Simulation"}
            </Button>
          )}
          {error && <span className="text-sm text-destructive">{error}</span>}
        </div>

        <div className="flex flex-wrap gap-4 text-sm">
          <StatItem label="Total" value={stats.totalSent} />
          <StatItem label="Success" value={stats.success} variant="secondary" />
          <StatItem label="Duplicates" value={stats.duplicates} variant="outline" />
          <StatItem label="Malformed" value={stats.malformed} variant="outline" />
          <StatItem label="Errors" value={stats.errors} variant="destructive" />
        </div>
      </div>
    </SectionCard>
  );
}

function StatItem({
  label,
  value,
  variant = "secondary",
}: {
  label: string;
  value: number;
  variant?: "secondary" | "outline" | "destructive";
}) {
  return (
    <div className="flex items-center gap-1.5">
      <span className="text-muted-foreground">{label}</span>
      <Badge variant={variant}>{value}</Badge>
    </div>
  );
}
