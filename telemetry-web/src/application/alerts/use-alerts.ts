import { useQuery } from "@tanstack/react-query";
import { useAlertRepository } from "@/app/repository-context";

// useAlerts — paginated list of alerts from alert-service.
export function useAlerts(params: { limit: number; offset: number }) {
  const repository = useAlertRepository();
  return useQuery({
    queryKey: ["alerts", params.limit, params.offset],
    queryFn: ({ signal }) => repository.list({ ...params, signal }),
  });
}
