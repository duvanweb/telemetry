import { useQuery } from "@tanstack/react-query";
import { useVehicleRepository } from "@/app/repository-context";

// useVehicles — paginated list of active vehicles from vehicle-service.
export function useVehicles(params: { limit: number; offset: number }) {
  const repository = useVehicleRepository();
  return useQuery({
    queryKey: ["vehicles", params.limit, params.offset],
    queryFn: ({ signal }) => repository.list({ ...params, signal }),
  });
}
