import { Ionicons } from "@expo/vector-icons";
import { Tabs, router } from "expo-router";
import { Pressable, Text } from "react-native";

import { useMonitoringStore } from "@/store/monitoring-store";
import { useVehicleStore } from "@/store/vehicle-store";
import { stopBackgroundTracking } from "@/tracking/background-task";
import {
  clearLastLocation,
  stopForegroundTracking,
} from "@/tracking/gps-tracker";

// Ensure the home tab shows first when entering the tabs group.
export const unstable_settings = {
  initialRouteName: "home",
};

// TabsLayout defines the bottom tab navigator shown when a vehicle is active.
// Home shows the monitoring/simulation screen, Alerts shows the alert list,
// and Exit clears the vehicle from the store and returns to the setup screen.
export default function TabsLayout() {
  const clearActiveVehicle = useVehicleStore((s) => s.clearActiveVehicle);

  const handleExit = () => {
    // Stop tracking first so no new positions are generated while tearing down.
    stopForegroundTracking();
    clearLastLocation();
    void stopBackgroundTracking();
    // Reset all monitoring state to defaults (prevents cross-vehicle data leakage).
    useMonitoringStore.getState().reset();
    // Clear the active vehicle and navigate back to the setup screen.
    clearActiveVehicle();
    router.replace("/setup");
  };

  return (
    <Tabs screenOptions={{ headerShown: false }}>
      <Tabs.Screen
        name="home"
        options={{
          title: "Home",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="home" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="alerts"
        options={{
          title: "Alerts",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="notifications" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="exit"
        options={{
          title: "Exit",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="log-out" size={size} color={color} />
          ),
          tabBarButton: (props) => (
            <Pressable
              onPress={handleExit}
              style={props.style}
              accessibilityRole="button"
              accessibilityLabel="Exit"
              className="flex-1 items-center justify-center py-3 active:opacity-70"
            >
              <Ionicons name="log-out" size={24} color="#dc2626" />
              <Text className="text-xs font-medium text-red-600">Exit</Text>
            </Pressable>
          ),
        }}
      />
    </Tabs>
  );
}
