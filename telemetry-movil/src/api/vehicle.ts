import { request } from "./client";
import type { Vehicle } from "./types";

// createVehicle calls POST /api/vehicles with the given plate.
export function createVehicle(plate: string): Promise<Vehicle> {
  return request<Vehicle>("/api/vehicles", {
    method: "POST",
    body: JSON.stringify({ plate }),
  });
}

// findByPlate calls GET /api/vehicles/plates/{plate}.
export function findByPlate(plate: string): Promise<Vehicle> {
  return request<Vehicle>(`/api/vehicles/plates/${encodeURIComponent(plate)}`, {
    method: "GET",
  });
}
