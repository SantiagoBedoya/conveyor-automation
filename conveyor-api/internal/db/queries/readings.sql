-- name: CreateReading :one
INSERT INTO sensor_readings (
    gas_value, humidity_value, distance_cm, object_count,
    belt_running, fan_on, buzzer_on, door_angle
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetReadingByID :one
SELECT * FROM sensor_readings WHERE id = $1;

-- name: ListReadings :many
SELECT * FROM sensor_readings ORDER BY timestamp DESC LIMIT $1 OFFSET $2;

-- name: ListReadingsSince :many
SELECT * FROM sensor_readings WHERE timestamp > $1 ORDER BY timestamp DESC;

-- name: GetLatestReading :one
SELECT * FROM sensor_readings ORDER BY timestamp DESC LIMIT 1;
