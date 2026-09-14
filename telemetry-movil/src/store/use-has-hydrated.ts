import { useEffect, useState } from "react";

import { useVehicleStore } from "./vehicle-store";

// useHasHydrated returns true once the vehicle store has finished rehydrating
// from AsyncStorage. The plain hasHydrated() function is synchronous and
// non-reactive, so we subscribe to onFinishHydration to trigger a re-render.
export function useHasHydrated(): boolean {
  const [hydrated, setHydrated] = useState(
    useVehicleStore.persist.hasHydrated(),
  );

  useEffect(() => {
    const unsub = useVehicleStore.persist.onFinishHydration(() =>
      setHydrated(true),
    );
    // Handle the edge case where hydration completes between the useState
    // initialization and the effect mount.
    if (useVehicleStore.persist.hasHydrated()) setHydrated(true);
    return unsub;
  }, []);

  return hydrated;
}
