import AsyncStorage from '@react-native-async-storage/async-storage';

import type { Vehicle } from '@/api/types';
import { useVehicleStore } from './vehicle-store';

jest.mock('@react-native-async-storage/async-storage', () =>
  require('@react-native-async-storage/async-storage/jest/async-storage-mock'),
);

const STORAGE_KEY = 'telemetry-movil:active-vehicle:v1';

const vehicle: Vehicle = {
  id: 1,
  plate: 'ABC-123',
  createdAt: '2026-09-08T00:00:00.000Z',
  updatedAt: '2026-09-08T00:00:00.000Z',
};

// Flush pending microtasks so the persist middleware's async AsyncStorage write completes.
const flushPromises = () => new Promise((resolve) => setTimeout(resolve, 0));

beforeEach(async () => {
  await AsyncStorage.clear();
  useVehicleStore.setState({ activeVehicle: null });
});

describe('vehicle-store', () => {
  it('persists the active vehicle on setActiveVehicle', async () => {
    useVehicleStore.getState().setActiveVehicle(vehicle);
    await flushPromises();

    const raw = await AsyncStorage.getItem(STORAGE_KEY);
    expect(raw).not.toBeNull();
    expect(JSON.parse(raw!).state.activeVehicle).toEqual(vehicle);
  });

  it('restores the active vehicle on hydration', async () => {
    await AsyncStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({ state: { activeVehicle: vehicle }, version: 0 }),
    );

    await useVehicleStore.persist.rehydrate();

    expect(useVehicleStore.getState().activeVehicle).toEqual(vehicle);
  });

  it('clears the active vehicle on clearActiveVehicle', async () => {
    useVehicleStore.getState().setActiveVehicle(vehicle);
    await flushPromises();

    useVehicleStore.getState().clearActiveVehicle();
    await flushPromises();

    expect(useVehicleStore.getState().activeVehicle).toBeNull();
    const raw = await AsyncStorage.getItem(STORAGE_KEY);
    expect(JSON.parse(raw!).state.activeVehicle).toBeNull();
  });
});
