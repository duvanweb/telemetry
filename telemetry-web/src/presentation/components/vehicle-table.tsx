import type { Vehicle } from "@/core/domain/vehicle";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/presentation/ui/table";
import { Button } from "@/presentation/ui/button";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString();
}

export interface SimulationTableProps {
  globalRunning: boolean;
  vehicleStopped: Record<number, boolean>;
  onToggleVehicle: (id: number) => void;
}

// VehicleTable — table of vehicles with columns ID, Plate, Created At, and optional Actions.
export function VehicleTable({
  vehicles,
  simulation,
}: {
  vehicles: Vehicle[];
  simulation?: SimulationTableProps;
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>ID</TableHead>
          <TableHead>Plate</TableHead>
          <TableHead>Created At</TableHead>
          {simulation && <TableHead className="text-right">Actions</TableHead>}
        </TableRow>
      </TableHeader>
      <TableBody>
        {vehicles.map((vehicle) => {
          const isSimVehicle = simulation && vehicle.id in simulation.vehicleStopped;
          const stopped = isSimVehicle ? simulation.vehicleStopped[vehicle.id] : false;

          return (
            <TableRow key={vehicle.id}>
              <TableCell>{vehicle.id}</TableCell>
              <TableCell className="font-medium">{vehicle.plate}</TableCell>
              <TableCell>{formatDate(vehicle.createdAt)}</TableCell>
              {simulation && (
                <TableCell className="text-right">
                  {isSimVehicle && (
                    <Button
                      size="sm"
                      variant={stopped ? "outline" : "destructive"}
                      disabled={!simulation.globalRunning}
                      onClick={() => simulation.onToggleVehicle(vehicle.id)}
                    >
                      {stopped ? "Start" : "Stop"}
                    </Button>
                  )}
                </TableCell>
              )}
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
