import type { ReactNode } from "react";
import { useAlerts } from "@/application/alerts/use-alerts";
import { ApiError } from "@/infrastructure/http/errors";
import { SectionCard } from "@/presentation/components/section-card";
import { ErrorState } from "@/presentation/components/error-state";
import { Badge } from "@/presentation/ui/badge";
import { Skeleton } from "@/presentation/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/presentation/ui/table";

function AlertsSkeleton() {
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

function formatTimestamp(iso: string): string {
  return new Date(iso).toLocaleString();
}

// AlertsSection — functional table of alerts from alert-service.
export function AlertsSection() {
  const query = useAlerts({ limit: 20, offset: 0 });

  let content: ReactNode;

  if (query.isLoading) {
    content = <AlertsSkeleton />;
  } else if (query.isError) {
    content = (
      <ErrorState
        message={networkMessage(query.error)}
        onRetry={() => void query.refetch()}
      />
    );
  } else if (!query.data || query.data.data.length === 0) {
    content = <p className="text-sm text-muted-foreground">No alerts found</p>;
  } else {
    content = (
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Vehicle ID</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Latitude</TableHead>
            <TableHead>Longitude</TableHead>
            <TableHead>Detected At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {query.data.data.map((alert) => (
            <TableRow key={alert.id}>
              <TableCell>{alert.id}</TableCell>
              <TableCell className="font-medium">{alert.vehicleId}</TableCell>
              <TableCell>
                <Badge variant={alert.type === "VEHICLE_STOPPED" ? "destructive" : "secondary"}>
                  {alert.type}
                </Badge>
              </TableCell>
              <TableCell>{alert.latitude.toFixed(4)}</TableCell>
              <TableCell>{alert.longitude.toFixed(4)}</TableCell>
              <TableCell>{formatTimestamp(alert.detectedAt)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    );
  }

  return (
    <SectionCard title="Alerts">
      {content}
    </SectionCard>
  );
}
