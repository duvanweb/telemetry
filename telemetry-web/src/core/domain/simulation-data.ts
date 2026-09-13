// Predefined simulation data — 5 vehicles with routes through Bogotá.
// Plates match the vehicle-service format ^[A-Z]{3}-[0-9]{3}$.
// Each route is a closed polygon of ~12 waypoints; the simulator interpolates
// between consecutive points so every tick yields a unique coordinate.

export interface Waypoint {
  lat: number;
  lng: number;
}

export interface SimulationVehicleDef {
  plate: string;
  route: Waypoint[];
}

export const SIMULATION_VEHICLES: SimulationVehicleDef[] = [
  {
    plate: "SIM-001",
    route: [
      { lat: 4.598, lng: -74.076 },
      { lat: 4.6, lng: -74.074 },
      { lat: 4.602, lng: -74.073 },
      { lat: 4.604, lng: -74.072 },
      { lat: 4.606, lng: -74.071 },
      { lat: 4.608, lng: -74.072 },
      { lat: 4.609, lng: -74.074 },
      { lat: 4.608, lng: -74.076 },
      { lat: 4.606, lng: -74.077 },
      { lat: 4.604, lng: -74.078 },
      { lat: 4.602, lng: -74.077 },
      { lat: 4.6, lng: -74.076 },
    ],
  },
  {
    plate: "SIM-002",
    route: [
      { lat: 4.648, lng: -74.067 },
      { lat: 4.65, lng: -74.065 },
      { lat: 4.652, lng: -74.064 },
      { lat: 4.654, lng: -74.063 },
      { lat: 4.656, lng: -74.064 },
      { lat: 4.657, lng: -74.066 },
      { lat: 4.658, lng: -74.068 },
      { lat: 4.656, lng: -74.07 },
      { lat: 4.654, lng: -74.071 },
      { lat: 4.652, lng: -74.07 },
      { lat: 4.65, lng: -74.068 },
      { lat: 4.649, lng: -74.067 },
    ],
  },
  {
    plate: "SIM-003",
    route: [
      { lat: 4.686, lng: -74.048 },
      { lat: 4.688, lng: -74.046 },
      { lat: 4.69, lng: -74.045 },
      { lat: 4.692, lng: -74.044 },
      { lat: 4.694, lng: -74.045 },
      { lat: 4.695, lng: -74.047 },
      { lat: 4.696, lng: -74.049 },
      { lat: 4.694, lng: -74.051 },
      { lat: 4.692, lng: -74.052 },
      { lat: 4.69, lng: -74.051 },
      { lat: 4.688, lng: -74.049 },
      { lat: 4.687, lng: -74.048 },
    ],
  },
  {
    plate: "SIM-004",
    route: [
      { lat: 4.661, lng: -74.119 },
      { lat: 4.663, lng: -74.117 },
      { lat: 4.665, lng: -74.116 },
      { lat: 4.667, lng: -74.115 },
      { lat: 4.669, lng: -74.116 },
      { lat: 4.67, lng: -74.118 },
      { lat: 4.671, lng: -74.12 },
      { lat: 4.669, lng: -74.122 },
      { lat: 4.667, lng: -74.123 },
      { lat: 4.665, lng: -74.122 },
      { lat: 4.663, lng: -74.12 },
      { lat: 4.662, lng: -74.119 },
    ],
  },
  {
    plate: "SIM-005",
    route: [
      { lat: 4.535, lng: -74.084 },
      { lat: 4.537, lng: -74.082 },
      { lat: 4.539, lng: -74.081 },
      { lat: 4.541, lng: -74.08 },
      { lat: 4.543, lng: -74.081 },
      { lat: 4.544, lng: -74.083 },
      { lat: 4.545, lng: -74.085 },
      { lat: 4.543, lng: -74.087 },
      { lat: 4.541, lng: -74.088 },
      { lat: 4.539, lng: -74.087 },
      { lat: 4.537, lng: -74.085 },
      { lat: 4.536, lng: -74.084 },
    ],
  },
];
