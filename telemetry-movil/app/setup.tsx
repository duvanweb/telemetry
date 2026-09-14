import { Redirect, router } from "expo-router";
import { Pressable, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { useHasHydrated } from "@/store/use-has-hydrated";
import { useVehicleStore } from "@/store/vehicle-store";

// SetupScreen is the initial screen shown when no active vehicle is in the
// store. It offers the user two entry points: register a new vehicle or
// search for an existing one by plate. If a vehicle is active (e.g. the user
// navigated back to setup on Android), it redirects to the home tab.
export default function SetupScreen() {
  const hydrated = useHasHydrated();
  if (!hydrated) return null;

  const activeVehicle = useVehicleStore((s) => s.activeVehicle);
  if (activeVehicle) return <Redirect href="/home" />;

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 px-6 py-8">
        <Text className="text-2xl font-bold text-slate-900">Vehicles</Text>
        <View className="mt-6 rounded-2xl border border-dashed border-slate-200 p-6">
          <Text className="text-center text-slate-500">
            No active vehicle yet.
          </Text>
          <Text className="mt-1 text-center text-sm text-slate-400">
            Register a new one or search by plate to start monitoring.
          </Text>
        </View>
        <View className="mt-auto gap-3">
          <Pressable
            onPress={() => router.push("/register")}
            className="items-center rounded-xl bg-indigo-600 px-4 py-4 active:opacity-70"
          >
            <Text className="text-base font-semibold text-white">
              Register vehicle
            </Text>
          </Pressable>
          <Pressable
            onPress={() => router.push("/search")}
            className="items-center rounded-xl border border-slate-300 px-4 py-4 active:opacity-70"
          >
            <Text className="text-base font-semibold text-slate-700">
              Search by plate
            </Text>
          </Pressable>
        </View>
      </View>
    </SafeAreaView>
  );
}
