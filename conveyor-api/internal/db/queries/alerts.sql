-- name: CreateAlert :one
INSERT INTO alerts (
    type, trigger_value, threshold
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAlertByID :one
SELECT * FROM alerts WHERE id = $1;

-- name: ListAlerts :many
SELECT * FROM alerts ORDER BY timestamp DESC LIMIT $1 OFFSET $2;

-- name: ListActiveAlerts :many
SELECT * FROM alerts WHERE active = true ORDER BY timestamp DESC;

-- name: ResolveAlert :one
UPDATE alerts SET active = false, resolved_at = NOW() WHERE id = $1 RETURNING *;
