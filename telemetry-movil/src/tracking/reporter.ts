import { ApiError } from "@/api/client";
import { reportPosition } from "@/api/geo";
import type { GpsPosition } from "@/api/types";
import { useMonitoringStore } from "@/store/monitoring-store";
import { useVehicleStore } from "@/store/vehicle-store";
import {
  SYNC_BACKOFF_BASE_MS,
  SYNC_BACKOFF_MAX_MS,
  SYNC_MAX_CONCURRENT,
  shouldSkipPosition,
} from "./adaptive-sampler";

// report sends a position to geo-service, or enqueues it when offline.
// It applies the "Parar vehículo" resend-last logic and the distance filter,
// then drains the offline queue with throttled concurrency when online.
export async function report(position: GpsPosition): Promise<void> {
  const vehicle = useVehicleStore.getState().activeVehicle;
  if (!vehicle) return;

  const store = useMonitoringStore.getState();

  // When the vehicle is stopped, resend the last known position (with a fresh
  // timestamp) instead of advancing to a new reading.
  let pos: GpsPosition;
  if (store.vehicleStopped) {
    if (!store.lastPosition) return;
    pos = { ...store.lastPosition, timestamp: new Date().toISOString() };
  } else {
    // Distance filter: skip negligible movements to save battery and bandwidth.
    if (shouldSkipPosition(store.lastPosition, position)) return;
    pos = position;
  }

  useMonitoringStore.getState().setLastPosition(pos);

  if (useMonitoringStore.getState().offline) {
    useMonitoringStore.getState().enqueue(pos);
    return;
  }

  // Online: send the current position, then drain the pending queue.
  try {
    await reportPosition(vehicle.id, pos);
  } catch (err) {
    if (err instanceof ApiError && err.status === 409) {
      // Duplicate within TTL — geo-service already has this point; treat as success.
    } else {
      // Network/server error: enqueue to retry later instead of dropping.
      useMonitoringStore.getState().enqueue(pos);
      return;
    }
  }

  void drainQueue(vehicle.id);
}

// syncQueue triggers a drain of the offline queue. Called by the monitoring
// screen when the user toggles offline mode off.
export async function syncQueue(): Promise<void> {
  const vehicle = useVehicleStore.getState().activeVehicle;
  if (!vehicle) return;
  if (useMonitoringStore.getState().offline) return;
  await drainQueue(vehicle.id);
}

// drainQueue sends queued positions to geo-service with at most
// SYNC_MAX_CONCURRENT concurrent requests, retrying with exponential backoff.
// Guarded by a module-level flag so concurrent triggers coalesce into one drain.
let draining = false;

async function drainQueue(vehicleId: number): Promise<void> {
  if (draining) return;
  draining = true;
  let attempt = 0;
  try {
    while (true) {
      const store = useMonitoringStore.getState();
      if (store.offline || store.queue.length === 0) break;

      // Take up to SYNC_MAX_CONCURRENT positions from the front of the queue.
      const batch: GpsPosition[] = [];
      for (let i = 0; i < SYNC_MAX_CONCURRENT; i++) {
        const next = useMonitoringStore.getState().dequeue();
        if (!next) break;
        batch.push(next);
      }
      if (batch.length === 0) break;

      const results = await Promise.allSettled(
        batch.map((pos) => reportPosition(vehicleId, pos)),
      );

      const failed = results
        .map((r, i) => (r.status === "rejected" ? batch[i] : null))
        .filter((p): p is GpsPosition => p !== null);

      if (failed.length > 0) {
        // Re-enqueue failed positions at the back and backoff before retrying.
        failed.forEach((p) => useMonitoringStore.getState().enqueue(p));
        attempt += 1;
        await sleep(Math.min(SYNC_BACKOFF_BASE_MS * 2 ** attempt, SYNC_BACKOFF_MAX_MS));
        continue;
      }

      attempt = 0;
    }
  } finally {
    draining = false;
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
