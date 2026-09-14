import L from "leaflet";
import type { LatLngExpression } from "leaflet";
import "leaflet/dist/leaflet.css";
import { MapContainer, Marker, Popup, TileLayer } from "react-leaflet";
import type { GpsPosition } from "@/core/domain/gps-position";
import { SectionCard } from "@/presentation/components/section-card";

const BOGOTA: LatLngExpression = [4.711, -74.0721];

// Custom div icon — avoids the leaflet default marker PNG bundling issue.
const pinIcon = L.divIcon({
  className: "",
  html: '<div style="font-size: 28px; line-height: 1;">📍</div>',
  iconSize: [28, 28],
  iconAnchor: [14, 28],
});

// MapSection — real-time map with OpenStreetMap tiles and markers for each vehicle position.
export function MapSection({ positions }: { positions: Record<number, GpsPosition> }) {
  const positionList = Object.values(positions);

  return (
    <SectionCard title="Vehicle tracking">
      <MapContainer
        center={BOGOTA}
        zoom={12}
        style={{ height: "400px", width: "100%" }}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        {positionList.map((pos) => (
          <Marker
            key={pos.vehicleId}
            position={[pos.latitude, pos.longitude]}
            icon={pinIcon}
          >
            <Popup>
              Vehicle {pos.vehicleId}
              <br />
              {pos.latitude.toFixed(4)}, {pos.longitude.toFixed(4)}
            </Popup>
          </Marker>
        ))}
      </MapContainer>
    </SectionCard>
  );
}
