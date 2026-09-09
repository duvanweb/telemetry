// Vehicle mirrors the VehicleResponse DTO from vehicle-service (SPEC 02).
// Dates arrive as ISO 8601 strings; id is a number (Go int64 within JS safe range).
export interface Vehicle {
  id: number;
  plate: string;
  createdAt: string;
  updatedAt: string;
}
