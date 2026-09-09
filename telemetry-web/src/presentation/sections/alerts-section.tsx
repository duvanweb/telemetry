import type { Alert } from "@/core/domain/alert";
import { SectionCard } from "@/presentation/components/section-card";
import { Badge } from "@/presentation/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/presentation/ui/table";

// Static example alerts for the mockup. Not wired to alert-service.
const alerts: Alert[] = [
  {
    id: 1,
    type: "SPEED",
    vehiclePlate: "ABC-123",
    message: "Speed limit exceeded (80 km/h)",
    timestamp: "2026-09-08T10:30:00Z",
  },
  {
    id: 2,
    type: "GEOFENCE",
    vehiclePlate: "DEF-456",
    message: "Exited geofence 'Zone A'",
    timestamp: "2026-09-08T11:15:00Z",
  },
  {
    id: 3,
    type: "SPEED",
    vehiclePlate: "GHI-789",
    message: "Speed limit exceeded (60 km/h)",
    timestamp: "2026-09-08T12:00:00Z",
  },
  {
    id: 4,
    type: "GEOFENCE",
    vehiclePlate: "ABC-123",
    message: "Entered geofence 'Zone B'",
    timestamp: "2026-09-08T12:45:00Z",
  },
];

function formatTimestamp(iso: string): string {
  return new Date(iso).toLocaleString();
}

// AlertsSection — mockup table of alerts with static example data and a "Demo" badge.
export function AlertsSection() {
  return (
    <SectionCard title="Alerts" demo>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Vehicle</TableHead>
            <TableHead>Message</TableHead>
            <TableHead>Timestamp</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {alerts.map((alert) => (
            <TableRow key={alert.id}>
              <TableCell>{alert.id}</TableCell>
              <TableCell>
                <Badge variant={alert.type === "SPEED" ? "destructive" : "secondary"}>
                  {alert.type}
                </Badge>
              </TableCell>
              <TableCell className="font-medium">{alert.vehiclePlate}</TableCell>
              <TableCell>{alert.message}</TableCell>
              <TableCell>{formatTimestamp(alert.timestamp)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </SectionCard>
  );
}
