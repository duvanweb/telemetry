import AsyncStorage from "@react-native-async-storage/async-storage";
import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

import type { GpsPosition, ReportState, TrackingMode } from "@/api/types";

// Soft cap on the offline queue to avoid unbounded growth when the driver
// forgets to disable offline mode. The oldest positions are dropped first.
const MAX_QUEUE_SIZE = 1000;

interface MonitoringStore {
  // reportState drives the status card badge (idle/reporting/stopped/error).
  reportState: ReportState;
  // mode selects the position source: real GPS or simulated route.
  mode: TrackingMode;
  // offline is the manual toggle; when true positions are queued, not sent.
  offline: boolean;
  // vehicleStopped is set by "Parar vehículo": the reporter resends lastPosition.
  vehicleStopped: boolean;
  // lastPosition is the most recent point shown in the card and resent when stopped.
  lastPosition: GpsPosition | null;
  // queue holds positions waiting to sync to geo-service (offline-first).
  queue: GpsPosition[];

  setReportState: (s: ReportState) => void;
  setMode: (m: TrackingMode) => void;
  setOffline: (o: boolean) => void;
  setVehicleStopped: (s: boolean) => void;
  setLastPosition: (p: GpsPosition) => void;
  enqueue: (p: GpsPosition) => void;
  dequeue: () => GpsPosition | undefined;
  clearQueue: () => void;
  // reset restores all monitoring state to defaults. Called when the user
  // exits (clears the active vehicle) to prevent cross-vehicle data leakage.
  reset: () => void;
}

// useMonitoringStore persists the monitoring state (including the offline queue)
// to AsyncStorage so pending positions survive app close/reopen. The key is
// versioned to allow future schema migrations.
export const useMonitoringStore = create<MonitoringStore>()(
  persist(
    (set, get) => ({
      reportState: "idle",
      mode: "simulation",
      offline: false,
      vehicleStopped: false,
      lastPosition: null,
      queue: [],

      setReportState: (s) => set({ reportState: s }),
      setMode: (m) => set({ mode: m }),
      setOffline: (o) => set({ offline: o }),
      setVehicleStopped: (s) => set({ vehicleStopped: s }),
      setLastPosition: (p) => set({ lastPosition: p }),
      enqueue: (p) =>
        set((state) => {
          const queue = [...state.queue, p];
          // drop oldest when over the soft cap
          if (queue.length > MAX_QUEUE_SIZE) {
            return { queue: queue.slice(queue.length - MAX_QUEUE_SIZE) };
          }
          return { queue };
        }),
      dequeue: () => {
        const [next, ...rest] = get().queue;
        set({ queue: rest });
        return next;
      },
      clearQueue: () => set({ queue: [] }),
      reset: () =>
        set({
          reportState: "idle",
          mode: "simulation",
          offline: false,
          vehicleStopped: false,
          lastPosition: null,
          queue: [],
        }),
    }),
    {
      name: "telemetry-movil:monitoring:v1",
      storage: createJSONStorage(() => AsyncStorage),
    },
  ),
);
