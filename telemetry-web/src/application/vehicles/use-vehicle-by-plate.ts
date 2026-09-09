import { useQuery } from "@tanstack/react-query";
import { useVehicleRepository } from "@/app/repository-context";

// useVehicleByPlate — search a single vehicle by plate from vehicle-service.
export function useVehicleByPlate(plate: string) {
  const repository = useVehicleRepository();
  return useQuery({
    queryKey: ["vehicles", "plate", plate],
    queryFn: ({ signal }) => repository.findByPlate(plate, signal),
    enabled: plate !== "",
  });
}
