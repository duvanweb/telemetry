import * as Location from "expo-location";
import type { LocationObject, LocationSubscription } from "expo-location";

import type { GpsPosition } from "@/api/types";

// Module-level state for the single foreground GPS subscription.
// The monitoring screen owns the tracking loop; this module just keeps the
// most recent reading available so the loop can sample it at the adaptive interval.
let subscription: LocationSubscription | null = null;
let lastLocation: LocationObject | null = null;

// requestForegroundPermissions asks the OS for foreground location access.
// Returns true when granted.
export async function requestForegroundPermissions(): Promise<boolean> {
  const { status } = await Location.requestForegroundPermissionsAsync();
  return status === "granted";
}

// requestBackgroundPermissions asks for background ("always") location access.
// Foreground permissions must be granted first.
export async function requestBackgroundPermissions(): Promise<boolean> {
  const { status } = await Location.requestBackgroundPermissionsAsync();
  return status === "granted";
}

// startForegroundTracking begins watching the device GPS. Updates are stored as
// the last known location (sampled by the monitoring loop at the adaptive interval).
// Returns false (and sets no subscription) when permission is denied.
export async function startForegroundTracking(): Promise<boolean> {
  const granted = await requestForegroundPermissions();
  if (!granted) return false;
  subscription = await Location.watchPositionAsync(
    {
      accuracy: Location.Accuracy.Balanced,
      timeInterval: 5000,
      distanceInterval: 10,
    },
    (location) => {
      lastLocation = location;
    },
  );
  return true;
}

// stopForegroundTracking removes the foreground GPS subscription.
export function stopForegroundTracking(): void {
  if (subscription) {
    subscription.remove();
    subscription = null;
  }
}

// getLastLocation returns the most recent GPS reading (null until the first fix).
export function getLastLocation(): LocationObject | null {
  return lastLocation;
}

// clearLastLocation resets the cached reading (used when stopping/restarting).
export function clearLastLocation(): void {
  lastLocation = null;
}

// locationToGpsPosition converts an expo-location LocationObject into the
// GpsPosition body expected by geo-service (lat, lng, RFC3339 timestamp).
export function locationToGpsPosition(location: LocationObject): GpsPosition {
  return {
    lat: location.coords.latitude,
    lng: location.coords.longitude,
    timestamp: new Date(location.timestamp).toISOString(),
  };
}
