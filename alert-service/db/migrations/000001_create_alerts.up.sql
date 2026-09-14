CREATE TABLE IF NOT EXISTS alerts (
    id          BIGSERIAL PRIMARY KEY,
    vehicle_id  BIGINT NOT NULL,
    type        TEXT NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    detected_at TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alerts_vehicle_id  ON alerts (vehicle_id);
CREATE INDEX IF NOT EXISTS idx_alerts_detected_at ON alerts (detected_at);
