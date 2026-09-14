import AsyncStorage from "@react-native-async-storage/async-storage";
import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

import type { Vehicle } from "@/api/types";

interface VehicleStore {
  // The single active vehicle the user is working with (null when none).
  activeVehicle: Vehicle | null;
  setActiveVehicle: (vehicle: Vehicle) => void;
  clearActiveVehicle: () => void;
}

// useVehicleStore persists the active vehicle to AsyncStorage so it survives
// app close/reopen. The key is versioned to allow future schema migrations.
export const useVehicleStore = create<VehicleStore>()(
  persist(
    (set) => ({
      activeVehicle: null,
      setActiveVehicle: (vehicle) => set({ activeVehicle: vehicle }),
      clearActiveVehicle: () => set({ activeVehicle: null }),
    }),
    {
      name: "telemetry-movil:active-vehicle:v1",
      storage: createJSONStorage(() => AsyncStorage),
    },
  ),
);
