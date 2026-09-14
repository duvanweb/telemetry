import {
  DEFAULT_INTERVAL_MS,
  LOW_BATTERY_INTERVAL_MS,
  MIN_DISTANCE_M,
  getAdaptiveInterval,
  haversineDistance,
  shouldSkipPosition,
} from "./adaptive-sampler";

describe("getAdaptiveInterval", () => {
  it("returns the default interval for a healthy battery", () => {
    expect(getAdaptiveInterval(0.9)).toBe(DEFAULT_INTERVAL_MS);
    expect(getAdaptiveInterval(0.5)).toBe(DEFAULT_INTERVAL_MS);
    expect(getAdaptiveInterval(0.2)).toBe(DEFAULT_INTERVAL_MS);
  });

  it("returns the low-battery interval when battery is below the threshold", () => {
    expect(getAdaptiveInterval(0.19)).toBe(LOW_BATTERY_INTERVAL_MS);
    expect(getAdaptiveInterval(0.05)).toBe(LOW_BATTERY_INTERVAL_MS);
    expect(getAdaptiveInterval(0)).toBe(LOW_BATTERY_INTERVAL_MS);
  });

  it("falls back to the default interval when the level is unknown (-1)", () => {
    expect(getAdaptiveInterval(-1)).toBe(DEFAULT_INTERVAL_MS);
  });
});

describe("haversineDistance", () => {
  it("returns 0 for the same point", () => {
    expect(haversineDistance({ lat: 4.71, lng: -74.07 }, { lat: 4.71, lng: -74.07 })).toBe(0);
  });

  it("returns a positive distance for distinct points", () => {
    const d = haversineDistance({ lat: 4.71, lng: -74.07 }, { lat: 4.72, lng: -74.07 });
    expect(d).toBeGreaterThan(0);
    // ~0.01 deg latitude ≈ 1.1 km
    expect(d).toBeGreaterThan(1000);
    expect(d).toBeLessThan(1200);
  });
});

describe("shouldSkipPosition", () => {
  it("never skips when there is no previous position", () => {
    expect(shouldSkipPosition(null, { lat: 4.71, lng: -74.07 })).toBe(false);
  });

  it("skips when the movement is below the minimum distance", () => {
    // ~0.00001 deg latitude ≈ 1.1 m, well under MIN_DISTANCE_M (10 m).
    const prev = { lat: 4.71, lng: -74.07 };
    const next = { lat: 4.71001, lng: -74.07 };
    expect(shouldSkipPosition(prev, next)).toBe(true);
  });

  it("does not skip when the movement meets the minimum distance", () => {
    // ~0.001 deg latitude ≈ 111 m, above MIN_DISTANCE_M (10 m).
    const prev = { lat: 4.71, lng: -74.07 };
    const next = { lat: 4.711, lng: -74.07 };
    expect(shouldSkipPosition(prev, next)).toBe(false);
  });

  it("uses MIN_DISTANCE_M as the threshold", () => {
    expect(MIN_DISTANCE_M).toBe(10);
  });
});
