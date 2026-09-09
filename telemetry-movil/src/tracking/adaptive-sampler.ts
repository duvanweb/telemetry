import * as Battery from "expo-battery";
import type { Subscription } from "expo-battery";

// Adaptive sampling constants (documented in README).
// Default interval is 5s; when battery drops below 20% the interval widens to
// 15s to reduce drain. The distance filter skips positions that moved less than
// 10m vs the last sent point, avoiding redundant reports in slow/stopped traffic.
export const DEFAULT_INTERVAL_MS = 5000;
export const LOW_BATTERY_INTERVAL_MS = 15000;
export const LOW_BATTERY_THRESHOLD = 0.2;
export const MIN_DISTANCE_M = 10;

// Throttle for offline queue sync: at most 2 concurrent POSTs to geo-service.
export const SYNC_MAX_CONCURRENT = 2;

// Exponential backoff bounds for sync retries (ms).
export const SYNC_BACKOFF_BASE_MS = 500;
export const SYNC_BACKOFF_MAX_MS = 30000;

// getAdaptiveInterval returns the sampling interval for the given battery level
// (0..1, or -1 when unknown). Unknown levels fall back to the default interval.
export function getAdaptiveInterval(batteryLevel: number): number {
  if (batteryLevel < 0) return DEFAULT_INTERVAL_MS;
  if (batteryLevel < LOW_BATTERY_THRESHOLD) return LOW_BATTERY_INTERVAL_MS;
  return DEFAULT_INTERVAL_MS;
}

// haversineDistance returns the great-circle distance in meters between two
// lat/lng points. Used by the distance filter to skip negligible movements.
export function haversineDistance(
  a: { lat: number; lng: number },
  b: { lat: number; lng: number },
): number {
  const R = 6371000; // Earth radius in meters
  const toRad = (deg: number) => (deg * Math.PI) / 180;
  const dLat = toRad(b.lat - a.lat);
  const dLng = toRad(b.lng - a.lng);
  const lat1 = toRad(a.lat);
  const lat2 = toRad(b.lat);
  const h =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(lat1) * Math.cos(lat2) * Math.sin(dLng / 2) ** 2;
  return 2 * R * Math.asin(Math.sqrt(h));
}

// shouldSkipPosition returns true when the movement between prev and next is
// below MIN_DISTANCE_M (redundant point). When prev is null (first point) it
// never skips.
export function shouldSkipPosition(
  prev: { lat: number; lng: number } | null,
  next: { lat: number; lng: number },
): boolean {
  if (!prev) return false;
  return haversineDistance(prev, next) < MIN_DISTANCE_M;
}

// subscribeBatteryLevel calls the listener immediately with the current level
// and then on every significant change. Returns an unsubscribe function.
export async function subscribeBatteryLevel(
  listener: (level: number) => void,
): Promise<() => void> {
  const initial = await Battery.getBatteryLevelAsync();
  listener(initial);
  const subscription: Subscription = Battery.addBatteryLevelListener(
    ({ batteryLevel }) => listener(batteryLevel),
  );
  return () => subscription.remove();
}
