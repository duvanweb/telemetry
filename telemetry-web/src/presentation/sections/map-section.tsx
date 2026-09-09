import L from "leaflet";
import type { LatLngExpression } from "leaflet";
import "leaflet/dist/leaflet.css";
import { MapContainer, Marker, TileLayer } from "react-leaflet";
import { SectionCard } from "@/presentation/components/section-card";
import { Badge } from "@/presentation/ui/badge";

const BOGOTA: LatLngExpression = [4.711, -74.0721];

// Mock vehicle positions around Bogotá.
const mockPositions: LatLngExpression[] = [
  [4.711, -74.0721],
  [4.6524, -74.0651],
  [4.6824, -74.0479],
];

// Custom div icon — avoids the leaflet default marker PNG bundling issue.
const pinIcon = L.divIcon({
  className: "",
  html: '<div style="font-size: 28px; line-height: 1;">📍</div>',
  iconSize: [28, 28],
  iconAnchor: [14, 28],
});

// MapSection — mockup map with OpenStreetMap tiles and 3 example markers.
// A "Demo" badge overlays the map. Not wired to geo-service.
export function MapSection() {
  return (
    <SectionCard title="Vehicle tracking">
      <div className="relative">
        <div className="absolute right-3 top-3 z-[1000]">
          <Badge variant="secondary">Demo</Badge>
        </div>
        <MapContainer
          center={BOGOTA}
          zoom={12}
          style={{ height: "400px", width: "100%" }}
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          {mockPositions.map((position, i) => (
            <Marker key={i} position={position} icon={pinIcon} />
          ))}
        </MapContainer>
      </div>
    </SectionCard>
  );
}
