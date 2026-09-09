// Vehicle entity — mirrors VehicleResponse from vehicle-service (SPEC 02).
export interface Vehicle {
  id: number;
  plate: string;
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
}
