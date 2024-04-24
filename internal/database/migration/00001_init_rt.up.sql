BEGIN;
CREATE TABLE machine_resource_usage (
    id SERIAL PRIMARY KEY, 
    machine VARCHAR(255) NOT NULL,
    metric1 INTEGER NOT NULL,
    metric2 INTEGER NOT NULL,
    metric3 INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
COMMIT;