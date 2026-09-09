import type { Vehicle } from "@/core/domain/vehicle";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/presentation/ui/table";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString();
}

// VehicleTable — table of vehicles with columns ID, Plate, Created At.
export function VehicleTable({ vehicles }: { vehicles: Vehicle[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>ID</TableHead>
          <TableHead>Plate</TableHead>
          <TableHead>Created At</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {vehicles.map((vehicle) => (
          <TableRow key={vehicle.id}>
            <TableCell>{vehicle.id}</TableCell>
            <TableCell className="font-medium">{vehicle.plate}</TableCell>
            <TableCell>{formatDate(vehicle.createdAt)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
