import { Redirect } from "expo-router";

import { useHasHydrated } from "@/store/use-has-hydrated";
import { useVehicleStore } from "@/store/vehicle-store";

// Index is the app entry point. It redirects to the setup screen when no
// active vehicle is in the store, or to the home tab when a vehicle exists.
// The hydration check prevents a flash of /setup before AsyncStorage loads
// the persisted vehicle.
export default function Index() {
  const hydrated = useHasHydrated();
  if (!hydrated) return null;

  const activeVehicle = useVehicleStore((s) => s.activeVehicle);

  if (!activeVehicle) return <Redirect href="/setup" />;
  return <Redirect href="/home" />;
}
