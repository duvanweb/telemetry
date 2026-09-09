import { useEffect, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { GpsPosition } from "@/core/domain/gps-position";
import { ALERT_SERVICE_URL } from "@/infrastructure/config/env";

// useAlertStream — subscribes to the alert-service SSE stream and returns real-time
// vehicle positions. On position events, updates the positions map. On alert events,
// invalidates the alerts query. If a position arrives for an unknown vehicle ID,
// invalidates the vehicles query so the new vehicle appears in the list.
export function useAlertStream(vehicleIds: number[]): Record<number, GpsPosition> {
  const queryClient = useQueryClient();
  const [positions, setPositions] = useState<Record<number, GpsPosition>>({});

  useEffect(() => {
    const es = new EventSource(`${ALERT_SERVICE_URL}/api/alerts/stream`);

    es.addEventListener("position", (event) => {
      const pos = JSON.parse(event.data) as GpsPosition;

      setPositions((prev) => ({ ...prev, [pos.vehicleId]: pos }));

      if (!vehicleIds.includes(pos.vehicleId)) {
        queryClient.invalidateQueries({ queryKey: ["vehicles"] });
      }
    });

    es.addEventListener("alert", () => {
      queryClient.invalidateQueries({ queryKey: ["alerts"] });
    });

    return () => {
      es.close();
    };
    // vehicleIds is passed as a dependency — the hook re-subscribes when it changes.
    // We use a stringified key to avoid re-subscribing on array reference changes.
  }, [queryClient, vehicleIds.join(",")]);

  return positions;
}
