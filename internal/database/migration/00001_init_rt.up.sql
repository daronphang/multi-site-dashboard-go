BEGIN;
CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE TABLE machine_resource_usage (
    id SERIAL, 
    machine VARCHAR(255) NOT NULL,
    metric1 INTEGER NOT NULL,
    metric2 INTEGER NOT NULL,
    metric3 INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY(machine, created_at)
);
SELECT create_hypertable('machine_resource_usage', by_range('created_at'));
COMMIT;