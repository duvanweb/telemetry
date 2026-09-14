CREATE TABLE gps_positions (
    id          BIGSERIAL PRIMARY KEY,
    vehicle_id  BIGINT NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gps_positions_vehicle_id  ON gps_positions (vehicle_id);
CREATE INDEX idx_gps_positions_recorded_at ON gps_positions (recorded_at);
