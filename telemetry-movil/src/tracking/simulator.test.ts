import { BOGOTA_ROUTE } from "./route";
import { RouteSimulator } from "./simulator";

describe("RouteSimulator", () => {
  it("traverses the waypoints in order and wraps around (closed loop)", () => {
    const sim = new RouteSimulator();
    const first = sim.next();
    expect(first).toEqual(BOGOTA_ROUTE[0]);

    const second = sim.next();
    expect(second).toEqual(BOGOTA_ROUTE[1]);

    // Advance to the end of the route and verify it wraps to the first point.
    for (let i = 2; i < BOGOTA_ROUTE.length; i++) {
      sim.next();
    }
    const wrapped = sim.next();
    expect(wrapped).toEqual(BOGOTA_ROUTE[0]);
  });

  it("returns all 15 distinct waypoints before repeating", () => {
    const sim = new RouteSimulator();
    const visited: string[] = [];
    for (let i = 0; i < BOGOTA_ROUTE.length; i++) {
      const p = sim.next();
      visited.push(`${p.lat},${p.lng}`);
    }
    expect(new Set(visited).size).toBe(BOGOTA_ROUTE.length);
  });

  it("reset returns the simulator to the first waypoint", () => {
    const sim = new RouteSimulator();
    sim.next();
    sim.next();
    sim.reset();
    expect(sim.currentIndex).toBe(0);
    expect(sim.next()).toEqual(BOGOTA_ROUTE[0]);
  });

  it("peek returns the current waypoint without advancing", () => {
    const sim = new RouteSimulator();
    expect(sim.peek()).toEqual(BOGOTA_ROUTE[0]);
    expect(sim.currentIndex).toBe(0);
    expect(sim.peek()).toEqual(BOGOTA_ROUTE[0]);
  });
});
