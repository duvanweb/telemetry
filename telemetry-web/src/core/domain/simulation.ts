// Simulation domain types — shapes for the vehicle telemetry simulator.

export interface SimulationVehicleStatus {
  id: number;
  plate: string;
  stopped: boolean;
}

export interface SimulationStats {
  totalSent: number;
  success: number;
  duplicates: number;
  malformed: number;
  errors: number;
}

export interface SimulationStatus {
  running: boolean;
  vehicles: SimulationVehicleStatus[];
  stats: SimulationStats;
}

export const idleSimulationStatus: SimulationStatus = {
  running: false,
  vehicles: [],
  stats: { totalSent: 0, success: 0, duplicates: 0, malformed: 0, errors: 0 },
};
