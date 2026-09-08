CREATE TABLE vehicles (
    id         BIGSERIAL PRIMARY KEY,
    plate      VARCHAR(7) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_plate_format CHECK (plate ~ '^[A-Z]{3}-[0-9]{3}$')
);
