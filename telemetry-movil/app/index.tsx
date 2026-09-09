import { useEffect, useRef, useState } from "react";
import { router } from "expo-router";
import { Pressable, ScrollView, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import type { GpsPosition } from "@/api/types";
import { useMonitoringStore } from "@/store/monitoring-store";
import { useVehicleStore } from "@/store/vehicle-store";
import {
  startBackgroundTracking,
  stopBackgroundTracking,
} from "@/tracking/background-task";
import {
  clearLastLocation,
  getLastLocation,
  locationToGpsPosition,
  startForegroundTracking,
  stopForegroundTracking,
} from "@/tracking/gps-tracker";
import {
  getAdaptiveInterval,
  subscribeBatteryLevel,
} from "@/tracking/adaptive-sampler";
import { report, syncQueue } from "@/tracking/reporter";
import { RouteSimulator } from "@/tracking/simulator";

const REPORT_STATE_LABEL: Record<string, string> = {
  idle: "Idle",
  reporting: "Reporting",
  stopped: "Stopped",
  error: "Error",
};

export default function Home() {
  const activeVehicle = useVehicleStore((s) => s.activeVehicle);
  const reportState = useMonitoringStore((s) => s.reportState);
  const mode = useMonitoringStore((s) => s.mode);
  const offline = useMonitoringStore((s) => s.offline);
  const vehicleStopped = useMonitoringStore((s) => s.vehicleStopped);
  const lastPosition = useMonitoringStore((s) => s.lastPosition);
  const setReportState = useMonitoringStore((s) => s.setReportState);
  const setMode = useMonitoringStore((s) => s.setMode);
  const setOffline = useMonitoringStore((s) => s.setOffline);
  const setVehicleStopped = useMonitoringStore((s) => s.setVehicleStopped);

  const [batteryLevel, setBatteryLevel] = useState(1);
  const simulatorRef = useRef(new RouteSimulator());

  // Subscribe to battery level changes for the adaptive sampling interval.
  useEffect(() => {
    let unsub = () => {};
    subscribeBatteryLevel(setBatteryLevel).then((fn) => {
      unsub = fn;
    });
    return unsub;
  }, []);

  // Tracking loop: ticks at the adaptive interval while reporting.
  useEffect(() => {
    if (reportState !== "reporting") return;
    const interval = getAdaptiveInterval(batteryLevel);

    const tick = () => {
      const store = useMonitoringStore.getState();
      // When the vehicle is stopped, the reporter resends lastPosition and we
      // do not advance the simulator or read new GPS samples.
      if (store.vehicleStopped) {
        if (!store.lastPosition) return;
        void report(store.lastPosition);
        return;
      }
      let pos: GpsPosition;
      if (store.mode === "simulation") {
        const p = simulatorRef.current.next();
        pos = { lat: p.lat, lng: p.lng, timestamp: new Date().toISOString() };
      } else {
        const loc = getLastLocation();
        if (!loc) return;
        pos = locationToGpsPosition(loc);
      }
      void report(pos);
    };

    tick();
    const id = setInterval(tick, interval);
    return () => clearInterval(id);
  }, [reportState, batteryLevel]);

  const handleStart = async () => {
    if (mode === "gps") {
      const ok = await startForegroundTracking();
      if (!ok) {
        setReportState("error");
        return;
      }
    }
    clearLastLocation();
    if (mode === "simulation") simulatorRef.current.reset();
    setVehicleStopped(false);
    setReportState("reporting");
    // Start the background task so tracking continues when the app is backgrounded.
    void startBackgroundTracking();
  };

  const handleStop = () => {
    setReportState("idle");
    stopForegroundTracking();
    void stopBackgroundTracking();
  };

  const handleModeChange = async (m: "gps" | "simulation") => {
    setMode(m);
    if (useMonitoringStore.getState().reportState === "reporting") {
      if (m === "gps") {
        await startForegroundTracking();
      } else {
        stopForegroundTracking();
      }
    }
  };

  const handleOfflineToggle = (value: boolean) => {
    setOffline(value);
    if (!value) void syncQueue();
  };

  // Empty state: no active vehicle yet.
  if (!activeVehicle) {
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

  const isReporting = reportState === "reporting";

  return (
    <SafeAreaView className="flex-1 bg-white">
      <ScrollView className="flex-1" contentContainerClassName="px-6 py-8 gap-5">
        <Text className="text-2xl font-bold text-slate-900">Monitoring</Text>

        {/* Status card (display only) */}
        <View className="rounded-2xl border border-slate-200 bg-slate-50 p-4 gap-2">
          <View className="flex-row items-center justify-between">
            <Text className="text-sm font-medium text-slate-500">
              Active vehicle
            </Text>
            <View
              className={`rounded-full px-2.5 py-1 ${
                reportState === "reporting"
                  ? "bg-emerald-100"
                  : reportState === "error"
                    ? "bg-red-100"
                    : "bg-slate-200"
              }`}
            >
              <Text
                className={`text-xs font-semibold ${
                  reportState === "reporting"
                    ? "text-emerald-700"
                    : reportState === "error"
                      ? "text-red-700"
                      : "text-slate-600"
                }`}
              >
                {REPORT_STATE_LABEL[reportState] ?? reportState}
              </Text>
            </View>
          </View>
          <Text className="text-xl font-bold text-slate-900">
            {activeVehicle.plate}
          </Text>
          <Text className="text-sm text-slate-500">ID: {activeVehicle.id}</Text>
          <View className="mt-2 flex-row gap-4">
            <Text className="text-sm text-slate-600">
              Lat: {lastPosition ? lastPosition.lat.toFixed(4) : "—"}
            </Text>
            <Text className="text-sm text-slate-600">
              Lng: {lastPosition ? lastPosition.lng.toFixed(4) : "—"}
            </Text>
          </View>
        </View>

        {/* Mode toggle */}
        <View className="gap-2">
          <Text className="text-sm font-medium text-slate-700">
            Position source
          </Text>
          <View className="flex-row gap-3">
            {(["simulation", "gps"] as const).map((m) => (
              <Pressable
                key={m}
                onPress={() => handleModeChange(m)}
                disabled={isReporting}
                className={`flex-1 items-center rounded-xl border px-4 py-3 ${
                  mode === m
                    ? "border-indigo-600 bg-indigo-50"
                    : "border-slate-300"
                } ${isReporting ? "opacity-50" : ""}`}
              >
                <Text
                  className={`text-sm font-semibold ${
                    mode === m ? "text-indigo-600" : "text-slate-700"
                  }`}
                >
                  {m === "simulation" ? "Simulation" : "GPS real"}
                </Text>
              </Pressable>
            ))}
          </View>
        </View>

        {/* Start / Stop */}
        <Pressable
          onPress={isReporting ? handleStop : handleStart}
          className={`items-center rounded-xl px-4 py-4 active:opacity-70 ${
            isReporting ? "bg-red-600" : "bg-indigo-600"
          }`}
        >
          <Text className="text-base font-semibold text-white">
            {isReporting ? "Stop reporting" : "Start"}
          </Text>
        </Pressable>

        {/* Parar vehículo */}
        <Pressable
          onPress={() => setVehicleStopped(!vehicleStopped)}
          disabled={!isReporting}
          className={`items-center rounded-xl border px-4 py-4 active:opacity-70 ${
            vehicleStopped
              ? "border-amber-500 bg-amber-50"
              : "border-slate-300"
          } ${!isReporting ? "opacity-50" : ""}`}
        >
          <Text
            className={`text-base font-semibold ${
              vehicleStopped ? "text-amber-600" : "text-slate-700"
            }`}
          >
            {vehicleStopped ? "Resume vehicle" : "Stop vehicle"}
          </Text>
        </Pressable>

        {/* Offline toggle */}
        <Pressable
          onPress={() => handleOfflineToggle(!offline)}
          className={`flex-row items-center justify-between rounded-xl border px-4 py-4 active:opacity-70 ${
            offline ? "border-indigo-600 bg-indigo-50" : "border-slate-300"
          }`}
        >
          <Text
            className={`text-base font-semibold ${
              offline ? "text-indigo-600" : "text-slate-700"
            }`}
          >
            Offline mode
          </Text>
          <Text className="text-sm text-slate-500">
            {offline ? "ON (queueing)" : "OFF"}
          </Text>
        </Pressable>

        {/* Panic button (no-op placeholder) */}
        <Pressable
          onPress={() => {
            /* no-op: panic button is a placeholder, implementation in a future spec */
          }}
          className="items-center rounded-xl border border-red-300 px-4 py-4 active:opacity-70"
        >
          <Text className="text-base font-semibold text-red-600">Panic</Text>
        </Pressable>
      </ScrollView>
    </SafeAreaView>
  );
}
