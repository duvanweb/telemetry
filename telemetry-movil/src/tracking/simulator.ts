import { BOGOTA_ROUTE } from "./route";

// RouteSimulator traverses BOGOTA_ROUTE cyclically, emitting one waypoint per
// call to next(). Used by the monitoring screen in simulation mode to emulate
// a vehicle driving through central Bogota without real GPS hardware.
export class RouteSimulator {
  private index = 0;

  // next returns the current waypoint and advances the internal index,
  // wrapping around to 0 after the last point (closed loop).
  next(): { lat: number; lng: number } {
    const point = BOGOTA_ROUTE[this.index];
    this.index = (this.index + 1) % BOGOTA_ROUTE.length;
    return point;
  }

  // peek returns the current waypoint without advancing the index.
  peek(): { lat: number; lng: number } {
    return BOGOTA_ROUTE[this.index];
  }

  // reset returns the simulator to the first waypoint.
  reset(): void {
    this.index = 0;
  }

  // get currentIndex exposes the internal position (useful for tests/debug).
  get currentIndex(): number {
    return this.index;
  }
}
