CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sensor_readings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    gas_value INTEGER NOT NULL,
    humidity_value INTEGER NOT NULL,
    distance_cm DOUBLE PRECISION NOT NULL,
    object_count INTEGER NOT NULL DEFAULT 0,
    belt_running BOOLEAN NOT NULL DEFAULT false,
    fan_on BOOLEAN NOT NULL DEFAULT false,
    buzzer_on BOOLEAN NOT NULL DEFAULT false,
    door_angle INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    type VARCHAR(32) NOT NULL CHECK (type IN ('GAS', 'HUMIDITY')),
    trigger_value INTEGER NOT NULL,
    threshold INTEGER NOT NULL,
    resolved_at TIMESTAMPTZ,
    active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_sensor_readings_timestamp ON sensor_readings (timestamp DESC);
CREATE INDEX idx_alerts_active ON alerts (active) WHERE active = true;
CREATE INDEX idx_alerts_timestamp ON alerts (timestamp DESC);
