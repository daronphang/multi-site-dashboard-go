-- name: GetMachineResourceUsage :many
SELECT * FROM machine_resource_usage
WHERE machine = $1;

-- name: CreateMachineResourceUsage :one
INSERT INTO machine_resource_usage (
  machine, metric1, metric2, metric3 
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- -- name: CreateMachineResourceUsage :exec
-- INSERT INTO machine_resource_usage (
--   machine, metric1, metric2, metric3 
-- ) VALUES (
--     unnest(@machines::VARCHAR[]),
--     unnest(@metrics1::INTEGER[]),
--     unnest(@metrics2::INTEGER[]),
--     unnest(@metrics3::INTEGER[]),
-- )