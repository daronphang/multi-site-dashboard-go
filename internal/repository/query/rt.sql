-- name: GetAggregatedMachineResourceUsage :many
SELECT 
time_bucket(CAST(@time_bucket::text AS INTERVAL), created_at)::timestamp AS bucket,
AVG(metric1) AS metric1,
AVG(metric2) AS metric2,
AVG(metric3) AS metric3
FROM machine_resource_usage
WHERE machine = $1
AND created_at > NOW() - CAST(@look_back_period::text AS INTERVAL)
GROUP BY bucket
ORDER BY bucket ASC;

-- name: CreateMachineResourceUsage :one
INSERT INTO machine_resource_usage (
  machine, metric1, metric2, metric3 
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateMachineResourceUsage :exec
UPDATE machine_resource_usage SET metric1 = $1 WHERE machine = $2;

-- -- name: CreateMachineResourceUsage :exec
-- INSERT INTO machine_resource_usage (
--   machine, metric1, metric2, metric3 
-- ) VALUES (
--     unnest(@machines::VARCHAR[]),
--     unnest(@metrics1::INTEGER[]),
--     unnest(@metrics2::INTEGER[]),
--     unnest(@metrics3::INTEGER[]),
-- )