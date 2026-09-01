-- Postgres will not let a newly added enum value be used in the transaction that adds it, and
-- migrations run one per transaction, so this file must contain nothing else.
ALTER TYPE broker ADD VALUE IF NOT EXISTS 'avance';
