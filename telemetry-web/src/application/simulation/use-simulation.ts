import { useEffect, useRef, useState } from "react";
import type { VehicleRepository } from "@/core/ports/vehicle-repository";
import type { PositionRepository } from "@/core/ports/position-repository";
import type { SimulationStatus, SimulationStats } from "@/core/domain/simulation";
import { idleSimulationStatus } from "@/core/domain/simulation";
import { SIMULATION_VEHICLES, type Waypoint } from "@/core/domain/simulation-data";
import { useVehicleRepository, usePositionRepository } from "@/app/repository-context";

// ─── Simulation Engine ─────────────────────────────────────────────────────────
// Manages 5 vehicles sending positions to geo-service every 2–5s with chaos injection.
// Testable independently of React.

interface VehicleSim {
  id: number;
  plate: string;
  route: Waypoint[];
  idx: number;
  progress: number;
  stopped: boolean;
  lastLat: number;
  lastLng: number;
}

const PROGRESS_STEP = 0.4;
const DUPLICATE_RATE = 0.1;
const MALFORMED_RATE = 0.05;

function randomInterval(): number {
  return 2000 + Math.random() * 3000; // 2–5 seconds
}

function errorStatus(error: unknown): number {
  if (error && typeof error === "object" && "status" in error) {
    return (error as { status: number }).status;
  }
  return 0;
}

function freshStats(): SimulationStats {
  return { totalSent: 0, success: 0, duplicates: 0, malformed: 0, errors: 0 };
}

class SimulationEngine {
  private vehicleRepo: VehicleRepository;
  private positionRepo: PositionRepository;
  private vehicles = new Map<number, VehicleSim>();
  private timers = new Map<number, ReturnType<typeof setTimeout>>();
  private stats: SimulationStats = freshStats();
  private running = false;

  constructor(vehicleRepo: VehicleRepository, positionRepo: PositionRepository) {
    this.vehicleRepo = vehicleRepo;
    this.positionRepo = positionRepo;
  }

  async start(): Promise<void> {
    if (this.running) return;

    this.stats = freshStats();
    this.vehicles.clear();

    for (const def of SIMULATION_VEHICLES) {
      let vehicle: { id: number };
      try {
        vehicle = await this.vehicleRepo.findByPlate(def.plate);
      } catch (error) {
        if (errorStatus(error) !== 404) throw error;
        vehicle = await this.vehicleRepo.create(def.plate);
      }

      const first = def.route[0];
      this.vehicles.set(vehicle.id, {
        id: vehicle.id,
        plate: def.plate,
        route: def.route,
        idx: 0,
        progress: 0,
        stopped: false,
        lastLat: first.lat,
        lastLng: first.lng,
      });
      this.scheduleTick(vehicle.id);
    }

    this.running = true;
  }

  stop(): void {
    for (const timer of this.timers.values()) clearTimeout(timer);
    this.timers.clear();
    this.vehicles.clear();
    this.running = false;
  }

  startVehicle(id: number): void {
    const vs = this.vehicles.get(id);
    if (vs) vs.stopped = false;
  }

  stopVehicle(id: number): void {
    const vs = this.vehicles.get(id);
    if (vs) vs.stopped = true;
  }

  getStatus(): SimulationStatus {
    return {
      running: this.running,
      vehicles: Array.from(this.vehicles.values()).map((v) => ({
        id: v.id,
        plate: v.plate,
        stopped: v.stopped,
      })),
      stats: { ...this.stats },
    };
  }

  private scheduleTick(id: number): void {
    const timer = setTimeout(() => {
      void this.tick(id);
      if (this.vehicles.has(id)) this.scheduleTick(id);
    }, randomInterval());
    this.timers.set(id, timer);
  }

  private async tick(id: number): Promise<void> {
    const vs = this.vehicles.get(id);
    if (!vs) return;

    let lat: number;
    let lng: number;

    if (!vs.stopped) {
      const from = vs.route[vs.idx];
      const to = vs.route[(vs.idx + 1) % vs.route.length];
      lat = from.lat + (to.lat - from.lat) * vs.progress;
      lng = from.lng + (to.lng - from.lng) * vs.progress;

      vs.progress += PROGRESS_STEP;
      if (vs.progress >= 1) {
        vs.idx = (vs.idx + 1) % vs.route.length;
        vs.progress -= 1;
      }
      vs.lastLat = lat;
      vs.lastLng = lng;
    } else {
      lat = vs.lastLat;
      lng = vs.lastLng;
    }

    // Normal send
    await this.sendAndCategorize(id, lat, lng);

    // Chaos: duplicate (10%) — resend identical coords, expects 409 from geo-service dedup
    if (Math.random() < DUPLICATE_RATE) {
      await this.sendAndCategorize(id, lat, lng);
    }

    // Chaos: malformed (5%) — send invalid JSON, expects 400
    if (Math.random() < MALFORMED_RATE) {
      await this.sendMalformedAndCategorize(id);
    }
  }

  private async sendAndCategorize(id: number, lat: number, lng: number): Promise<void> {
    this.stats.totalSent++;
    try {
      await this.positionRepo.sendPosition(id, lat, lng, new Date().toISOString());
      this.stats.success++;
    } catch (error) {
      const status = errorStatus(error);
      if (status === 409) this.stats.duplicates++;
      else if (status === 400) this.stats.malformed++;
      else this.stats.errors++;
    }
  }

  private async sendMalformedAndCategorize(id: number): Promise<void> {
    this.stats.totalSent++;
    try {
      await this.positionRepo.sendMalformed(id);
      this.stats.success++;
    } catch (error) {
      const status = errorStatus(error);
      if (status === 400) this.stats.malformed++;
      else this.stats.errors++;
    }
  }
}

// ─── Hook ──────────────────────────────────────────────────────────────────────

export function useSimulation() {
  const vehicleRepo = useVehicleRepository();
  const positionRepo = usePositionRepository();
  const engineRef = useRef<SimulationEngine | null>(null);
  const [status, setStatus] = useState<SimulationStatus>(idleSimulationStatus);
  const [starting, setStarting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!engineRef.current) {
    engineRef.current = new SimulationEngine(vehicleRepo, positionRepo);
  }

  // Poll engine state every 1s for UI updates
  useEffect(() => {
    const interval = setInterval(() => {
      if (engineRef.current) setStatus(engineRef.current.getStatus());
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  const start = async () => {
    if (!engineRef.current) return;
    setStarting(true);
    setError(null);
    try {
      await engineRef.current.start();
      setStatus(engineRef.current.getStatus());
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setStarting(false);
    }
  };

  const stop = () => {
    if (!engineRef.current) return;
    engineRef.current.stop();
    setStatus(engineRef.current.getStatus());
  };

  const startVehicle = (id: number) => {
    engineRef.current?.startVehicle(id);
    if (engineRef.current) setStatus(engineRef.current.getStatus());
  };

  const stopVehicle = (id: number) => {
    engineRef.current?.stopVehicle(id);
    if (engineRef.current) setStatus(engineRef.current.getStatus());
  };

  return { status, starting, error, start, stop, startVehicle, stopVehicle };
}

export type SimulationControls = ReturnType<typeof useSimulation>;
