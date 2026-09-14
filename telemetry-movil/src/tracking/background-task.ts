import * as Location from "expo-location";
import * as TaskManager from "expo-task-manager";

import { report } from "@/tracking/reporter";
import { locationToGpsPosition } from "@/tracking/gps-tracker";

// BACKGROUND_LOCATION_TASK is the task name registered with expo-task-manager
// and started via Location.startLocationUpdatesAsync in the app layout.
export const BACKGROUND_LOCATION_TASK = "telemetry-background-location";

// Define the background location task at top-level scope (required by
// TaskManager: background launches run the JS bundle without mounting views).
// On each batch of location updates while the app is backgrounded, it takes the
// most recent reading and hands it to the reporter, which applies the offline
// queue / sync logic exactly like the foreground tracking loop.
TaskManager.defineTask(BACKGROUND_LOCATION_TASK, async ({ data, error }) => {
  if (error) return;
  if (!data) return;

  const { locations } = data as { locations: Location.LocationObject[] };
  const location = locations?.[locations.length - 1];
  if (!location) return;

  await report(locationToGpsPosition(location));
});

// startBackgroundTracking registers and starts the OS background location task.
// Returns false when permissions are not granted.
export async function startBackgroundTracking(): Promise<boolean> {
  const granted = await Location.requestBackgroundPermissionsAsync();
  if (granted.status !== "granted") return false;
  await Location.startLocationUpdatesAsync(BACKGROUND_LOCATION_TASK, {
    accuracy: Location.Accuracy.Balanced,
    timeInterval: 5000,
    distanceInterval: 10,
    showsBackgroundLocationIndicator: true,
  });
  return true;
}

// stopBackgroundTracking stops the OS background location task.
export async function stopBackgroundTracking(): Promise<void> {
  try {
    await Location.stopLocationUpdatesAsync(BACKGROUND_LOCATION_TASK);
  } catch {
    // task was not running; ignore
  }
}
