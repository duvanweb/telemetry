import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation } from '@tanstack/react-query';
import { router } from 'expo-router';
import { Controller, useForm } from 'react-hook-form';
import {
  ActivityIndicator,
  Pressable,
  Text,
  TextInput,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { ApiError } from '@/api/client';
import { createVehicle } from '@/api/vehicle';
import type { Vehicle } from '@/api/types';
import { plateFormSchema } from '@/lib/plate-schema';
import { useVehicleStore } from '@/store/vehicle-store';

export default function RegisterScreen() {
  const setActiveVehicle = useVehicleStore((s) => s.setActiveVehicle);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<{ plate: string }>({
    resolver: zodResolver(plateFormSchema),
    defaultValues: { plate: '' },
  });

  const mutation = useMutation({
    mutationFn: createVehicle,
    onSuccess: (vehicle: Vehicle) => {
      setActiveVehicle(vehicle);
      router.replace('/');
    },
  });

  const onSubmit = handleSubmit((data) => mutation.mutate(data.plate));

  const errorMessage = (() => {
    const error = mutation.error;
    if (!error) return null;
    if (error instanceof ApiError && error.status === 409) return 'Vehicle already exists';
    return 'Network error, try again';
  })();

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 px-6 py-8">
        <Text className="text-2xl font-bold text-slate-900">Register vehicle</Text>

        <Controller
          control={control}
          name="plate"
          render={({ field: { onChange, onBlur, value } }) => (
            <View className="mt-6">
              <Text className="mb-2 text-sm font-medium text-slate-700">Plate</Text>
              <TextInput
                value={value}
                onChangeText={(text) => {
                  if (mutation.isError) mutation.reset();
                  onChange(text);
                }}
                onBlur={onBlur}
                placeholder="ABC-123"
                autoCapitalize="characters"
                autoCorrect={false}
                maxLength={7}
                className="rounded-xl border border-slate-300 px-4 py-3 text-base text-slate-900"
              />
              {errors.plate && (
                <Text className="mt-1 text-sm text-red-500">{errors.plate.message}</Text>
              )}
            </View>
          )}
        />

        {errorMessage && <Text className="mt-4 text-sm text-red-500">{errorMessage}</Text>}

        <Pressable
          onPress={onSubmit}
          disabled={mutation.isPending}
          className="mt-6 items-center rounded-xl bg-indigo-600 px-4 py-4 active:opacity-70"
        >
          {mutation.isPending ? (
            <ActivityIndicator color="white" />
          ) : (
            <Text className="text-base font-semibold text-white">Register</Text>
          )}
        </Pressable>
      </View>
    </SafeAreaView>
  );
}
