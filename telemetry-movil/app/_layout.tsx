import "../global.css";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Stack } from "expo-router";
import { SafeAreaProvider } from "react-native-safe-area-context";

// Importing the background-task module registers the TaskManager.defineTask
// handler at top-level scope, which is required before the app can start
// background location updates (see src/tracking/background-task.ts).
import "@/tracking/background-task";

const queryClient = new QueryClient();

// RootLayout sets up the providers and the Stack navigator. The Stack contains
// the non-tab screens (index redirect, setup, register, search) and the (tabs)
// route group which provides the bottom tab navigator for Home, Alerts, Exit.
// All headers are hidden — each screen manages its own SafeAreaView layout.
export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <SafeAreaProvider>
        <Stack screenOptions={{ headerShown: false }}>
          <Stack.Screen name="index" />
          <Stack.Screen name="setup" />
          <Stack.Screen name="register" />
          <Stack.Screen name="search" />
          <Stack.Screen name="(tabs)" />
        </Stack>
      </SafeAreaProvider>
    </QueryClientProvider>
  );
}
