-- migrations/000010_create_audit_logs.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS audit_logs ( 
    id              UUID              PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type     VARCHAR(50)       NOT NULL,
    entity_id       UUID              NOT NULL,
    action          VARCHAR(50)       NOT NULL,
    actor_id        UUID              NOT NULL,
    actor_type      VARCHAR(20)       NOT NULL DEFAULT 'system',
    old_value       JSONB,
    new_value       JSONB,
    ip_address      VARCHAR(45),
    created_at      TIMESTAMPTZ      NOT NULL DEFAULT NOW()

    CONSTRAINT chk_actor_type 
        CHECK (actor_type IN ('customer','admin','system'))
);

-- Fast Lookup: All audit logs for a specific entity id 
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_id
    ON audit_logs(entity_id);

-- Fast Lookup: All audit logs for a specific actor id 
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_id
    ON audit_logs(actor_id);

-- Fast Lookup: All audit logs for a specific time 
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at
    ON audit_logs(created_at);

COMMIT;