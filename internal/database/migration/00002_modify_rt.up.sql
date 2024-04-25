BEGIN;
ALTER TABLE machine_resource_usage ADD CONSTRAINT CHK_METRIC
CHECK (metric1 >= 0 AND metric2 >=0 AND metric3 >= 0);
COMMIT;