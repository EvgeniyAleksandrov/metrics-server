CREATE TABLE IF NOT EXISTS metrics (
   id    TEXT PRIMARY KEY,
   type  TEXT NOT NULL,
   delta BIGINT,
   value DOUBLE PRECISION

   CONSTRAINT metrics_type_value_check
   CHECK (
        (type = 'counter' AND delta IS NOT NULL AND value IS NULL)
        OR
        (type = 'gauge'   AND value IS NOT NULL AND delta IS NULL)
    )
);