-- name: GetSystemStatus :many
WITH latest_reading AS (
    SELECT * FROM sensor_readings ORDER BY timestamp DESC LIMIT 1
),
active_alerts_count AS (
    SELECT count(*) AS total FROM alerts WHERE active = true
)
SELECT
    l.timestamp AS last_seen,
    l.belt_running,
    l.fan_on,
    l.buzzer_on,
    l.door_angle,
    COALESCE((SELECT sum(object_count) FROM sensor_readings), 0)::bigint AS object_count_total,
    COALESCE(aac.total, 0)::int AS active_alerts
FROM latest_reading l
CROSS JOIN active_alerts_count aac;
