/// <reference types="jest" />
import AsyncStorage from "@react-native-async-storage/async-storage";

import type { GpsPosition } from "@/api/types";
import { useMonitoringStore } from "./monitoring-store";

jest.mock("@react-native-async-storage/async-storage", () =>
  require("@react-native-async-storage/async-storage/jest/async-storage-mock"),
);

const STORAGE_KEY = "telemetry-movil:monitoring:v1";

const position = (lat: number): GpsPosition => ({
  lat,
  lng: -74.07,
  timestamp: "2026-09-08T14:30:00.000Z",
});

// Flush pending microtasks so the persist middleware's async AsyncStorage write completes.
const flushPromises = () => new Promise((resolve) => setTimeout(resolve, 0));

beforeEach(async () => {
  await AsyncStorage.clear();
  useMonitoringStore.setState({
    reportState: "idle",
    mode: "simulation",
    offline: false,
    vehicleStopped: false,
    lastPosition: null,
    queue: [],
  });
});

describe("monitoring-store", () => {
  it("enqueue adds a position to the queue and persists it", async () => {
    useMonitoringStore.getState().enqueue(position(4.71));
    await flushPromises();

    const queue = useMonitoringStore.getState().queue;
    expect(queue).toHaveLength(1);
    expect(queue[0].lat).toBe(4.71);

    const raw = await AsyncStorage.getItem(STORAGE_KEY);
    expect(raw).not.toBeNull();
    expect(JSON.parse(raw!).state.queue).toHaveLength(1);
  });

  it("dequeue removes and returns the first position", async () => {
    useMonitoringStore.getState().enqueue(position(4.71));
    useMonitoringStore.getState().enqueue(position(4.72));
    await flushPromises();

    const first = useMonitoringStore.getState().dequeue();
    await flushPromises();

    expect(first?.lat).toBe(4.71);
    expect(useMonitoringStore.getState().queue).toHaveLength(1);
    expect(useMonitoringStore.getState().queue[0].lat).toBe(4.72);
  });

  it("restores the queue, reportState and mode on hydration", async () => {
    await AsyncStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        state: {
          reportState: "reporting",
          mode: "gps",
          offline: true,
          vehicleStopped: false,
          lastPosition: position(4.71),
          queue: [position(4.71), position(4.72)],
        },
        version: 0,
      }),
    );

    await useMonitoringStore.persist.rehydrate();

    const state = useMonitoringStore.getState();
    expect(state.reportState).toBe("reporting");
    expect(state.mode).toBe("gps");
    expect(state.offline).toBe(true);
    expect(state.queue).toHaveLength(2);
  });

  it("clearQueue empties the queue and persists", async () => {
    useMonitoringStore.getState().enqueue(position(4.71));
    useMonitoringStore.getState().enqueue(position(4.72));
    await flushPromises();

    useMonitoringStore.getState().clearQueue();
    await flushPromises();

    expect(useMonitoringStore.getState().queue).toHaveLength(0);
    const raw = await AsyncStorage.getItem(STORAGE_KEY);
    expect(JSON.parse(raw!).state.queue).toHaveLength(0);
  });

  it("drops oldest positions when the queue exceeds the soft cap", async () => {
    // Fill beyond MAX_QUEUE_SIZE (1000) and verify it caps.
    for (let i = 0; i < 1005; i++) {
      useMonitoringStore.getState().enqueue(position(i));
    }
    await flushPromises();

    const queue = useMonitoringStore.getState().queue;
    expect(queue).toHaveLength(1000);
    // The oldest 5 (lat 0..4) should have been dropped.
    expect(queue[0].lat).toBe(5);
  });
});
