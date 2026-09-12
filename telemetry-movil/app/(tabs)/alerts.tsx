import { useQuery } from "@tanstack/react-query";
import { Pressable, ScrollView, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { ApiError } from "@/api/client";
import { listAlerts } from "@/api/alert";
import type { Alert } from "@/api/types";

// AlertTypeBadge renders a colored badge for the alert type.
function AlertTypeBadge({ type }: { type: string }) {
  const styles =
    type === "VEHICLE_STOPPED"
      ? "bg-amber-100"
      : "bg-slate-200";
  const textStyles =
    type === "VEHICLE_STOPPED"
      ? "text-amber-700"
      : "text-slate-600";

  return (
    <View className={`rounded-full px-2.5 py-1 ${styles}`}>
      <Text className={`text-xs font-semibold ${textStyles}`}>
        {type.replace(/_/g, " ")}
      </Text>
    </View>
  );
}

// AlertCard renders a single alert in the list.
function AlertCard({ alert }: { alert: Alert }) {
  const detected = new Date(alert.detectedAt).toLocaleString();

  return (
    <View className="rounded-2xl border border-slate-200 bg-slate-50 p-4 gap-2">
      <View className="flex-row items-center justify-between">
        <Text className="text-sm font-medium text-slate-500">
          Vehicle #{alert.vehicleId}
        </Text>
        <AlertTypeBadge type={alert.type} />
      </View>
      <View className="mt-1 flex-row gap-4">
        <Text className="text-sm text-slate-600">
          Lat: {alert.latitude.toFixed(4)}
        </Text>
        <Text className="text-sm text-slate-600">
          Lng: {alert.longitude.toFixed(4)}
        </Text>
      </View>
      <Text className="text-xs text-slate-400">{detected}</Text>
    </View>
  );
}

// AlertsScreen shows the alert list in the Alerts tab. The query refetches
// every time the screen mounts because staleTime is 0.
export default function AlertsScreen() {
  const { data, isLoading, isError, error, refetch, isFetching } = useQuery({
    queryKey: ["alerts"],
    queryFn: () => listAlerts(20, 0),
    staleTime: 0,
  });

  return (
    <SafeAreaView className="flex-1 bg-white">
      <ScrollView
        className="flex-1"
        contentContainerClassName="px-6 py-8 gap-4"
      >
        <Text className="text-2xl font-bold text-slate-900">Alerts</Text>

        {isLoading && (
          <View className="items-center py-12">
            <Text className="text-slate-400">Loading alerts…</Text>
          </View>
        )}

        {isError && (
          <View className="rounded-2xl border border-dashed border-red-200 p-6 gap-3">
            <Text className="text-center text-red-500">
              {error instanceof ApiError
                ? `Error ${error.status}: ${error.message}`
                : "Network error, try again"}
            </Text>
            <Pressable
              onPress={() => refetch()}
              className="items-center rounded-xl bg-indigo-600 px-4 py-3 active:opacity-70"
            >
              <Text className="text-sm font-semibold text-white">Retry</Text>
            </Pressable>
          </View>
        )}

        {data && data.data.length === 0 && (
          <View className="rounded-2xl border border-dashed border-slate-200 p-6">
            <Text className="text-center text-slate-500">No alerts yet.</Text>
          </View>
        )}

        {data && data.data.length > 0 && (
          <>
            {isFetching && (
              <Text className="text-sm text-slate-400">Refreshing…</Text>
            )}
            {data.data.map((alert) => (
              <AlertCard key={alert.id} alert={alert} />
            ))}
            <Text className="text-center text-sm text-slate-400">
              {data.data.length} of {data.total} alerts
            </Text>
          </>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}
