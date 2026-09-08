import { router } from 'expo-router';
import { Pressable, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { useVehicleStore } from '@/store/vehicle-store';

export default function Home() {
  const activeVehicle = useVehicleStore((s) => s.activeVehicle);

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 px-6 py-8">
        <Text className="text-2xl font-bold text-slate-900">Vehicles</Text>

        {activeVehicle ? (
          <View className="mt-6 rounded-2xl border border-slate-200 bg-slate-50 p-4">
            <Text className="text-sm font-medium text-slate-500">Active vehicle</Text>
            <Text className="mt-1 text-xl font-bold text-slate-900">
              {activeVehicle.plate}
            </Text>
            <Text className="mt-1 text-sm text-slate-500">ID: {activeVehicle.id}</Text>
          </View>
        ) : (
          <View className="mt-6 rounded-2xl border border-dashed border-slate-200 p-6">
            <Text className="text-center text-slate-500">No active vehicle yet.</Text>
            <Text className="mt-1 text-center text-sm text-slate-400">
              Register a new one or search by plate.
            </Text>
          </View>
        )}

        <View className="mt-auto gap-3">
          <Pressable
            onPress={() => router.push('/register')}
            className="items-center rounded-xl bg-indigo-600 px-4 py-4 active:opacity-70"
          >
            <Text className="text-base font-semibold text-white">Register vehicle</Text>
          </Pressable>
          <Pressable
            onPress={() => router.push('/search')}
            className="items-center rounded-xl border border-slate-300 px-4 py-4 active:opacity-70"
          >
            <Text className="text-base font-semibold text-slate-700">Search by plate</Text>
          </Pressable>
        </View>
      </View>
    </SafeAreaView>
  );
}
