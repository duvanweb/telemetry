import type { ReactNode } from "react";
import { useSearchParams } from "react-router-dom";
import { useVehicles } from "@/application/vehicles/use-vehicles";
import { useVehicleByPlate } from "@/application/vehicles/use-vehicle-by-plate";
import { VEHICLES_PAGE_SIZE } from "@/infrastructure/config/env";
import { ApiError } from "@/infrastructure/http/errors";
import { VehicleSearch } from "@/presentation/components/vehicle-search";
import { VehicleTable, type SimulationTableProps } from "@/presentation/components/vehicle-table";
import { PaginationControls } from "@/presentation/components/pagination-controls";
import { SectionCard } from "@/presentation/components/section-card";
import { ErrorState } from "@/presentation/components/error-state";
import { Skeleton } from "@/presentation/ui/skeleton";
import type { SimulationControls } from "@/application/simulation/use-simulation";

function TableSkeleton() {
  return (
    <div className="space-y-2">
      {Array.from({ length: 5 }).map((_, i) => (
        <Skeleton key={i} className="h-10 w-full" />
      ))}
    </div>
  );
}

function networkMessage(error: unknown): string {
  if (error instanceof ApiError && error.status === 0) return "Network error, try again";
  if (error instanceof ApiError) return error.message;
  return "Something went wrong";
}

// VehiclesSection — functional vehicle list with search by plate and pagination.
// State (search query + page) is synced to URL search params.
export function VehiclesSection({ simulation }: { simulation?: SimulationControls }) {
  const [searchParams, setSearchParams] = useSearchParams();
  const q = searchParams.get("q") ?? "";
  const page = Math.max(1, Number(searchParams.get("page") ?? "1"));

  const listQuery = useVehicles({
    limit: VEHICLES_PAGE_SIZE,
    offset: (page - 1) * VEHICLES_PAGE_SIZE,
  });
  const plateQuery = useVehicleByPlate(q);

  // Build simulation table props from the simulation status (if provided).
  const simTableProps: SimulationTableProps | undefined = simulation
    ? {
        globalRunning: simulation.status.running,
        vehicleStopped: Object.fromEntries(
          simulation.status.vehicles.map((v) => [v.id, v.stopped]),
        ),
        onToggleVehicle: (id: number) => {
          const v = simulation.status.vehicles.find((v) => v.id === id);
          if (v?.stopped) simulation.startVehicle(id);
          else simulation.stopVehicle(id);
        },
      }
    : undefined;

  function handleChangePage(newPage: number) {
    const next = new URLSearchParams(searchParams);
    next.set("page", String(newPage));
    setSearchParams(next);
  }

  let content: ReactNode;

  if (q !== "") {
    if (plateQuery.isLoading) {
      content = <TableSkeleton />;
    } else if (plateQuery.isError) {
      const notFound =
        plateQuery.error instanceof ApiError && plateQuery.error.status === 404;
      content = (
        <ErrorState
          message={notFound ? "Vehicle not found" : networkMessage(plateQuery.error)}
          onRetry={() => void plateQuery.refetch()}
        />
      );
    } else if (plateQuery.data) {
      content = <VehicleTable vehicles={[plateQuery.data]} simulation={simTableProps} />;
    } else {
      content = null;
    }
  } else {
    if (listQuery.isLoading) {
      content = <TableSkeleton />;
    } else if (listQuery.isError) {
      content = (
        <ErrorState
          message={networkMessage(listQuery.error)}
          onRetry={() => void listQuery.refetch()}
        />
      );
    } else if (listQuery.data && listQuery.data.data.length > 0) {
      const totalPages = Math.max(
        1,
        Math.ceil(listQuery.data.total / VEHICLES_PAGE_SIZE),
      );
      content = (
        <>
          <VehicleTable vehicles={listQuery.data.data} simulation={simTableProps} />
          <PaginationControls
            page={page}
            totalPages={totalPages}
            total={listQuery.data.total}
            onChangePage={handleChangePage}
          />
        </>
      );
    } else {
      content = <p className="text-sm text-muted-foreground">No vehicles found</p>;
    }
  }

  return (
    <SectionCard title="Vehicles">
      <div className="space-y-4">
        <VehicleSearch />
        {content}
      </div>
    </SectionCard>
  );
}
