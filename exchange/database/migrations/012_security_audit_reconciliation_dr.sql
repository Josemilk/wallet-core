-- 012: authentication/RBAC, immutable audit trail, reconciliation and DR metadata.
-- No private keys, bearer tokens, JWTs, seeds or raw transaction secrets are stored here.

CREATE TABLE IF NOT EXISTS app_roles (
    role_name TEXT PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id TEXT NOT NULL,
    role_name TEXT NOT NULL REFERENCES app_roles(role_name),
    granted_by TEXT NOT NULL,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_name)
);

INSERT INTO app_roles(role_name, description) VALUES
 ('user','normal exchange user'),
 ('support','read-only operational support'),
 ('compliance','compliance and review operations'),
 ('treasury','treasury operations'),
 ('admin','administrative operations'),
 ('security','security/audit operations')
ON CONFLICT (role_name) DO NOTHING;

CREATE TABLE IF NOT EXISTS audit_log (
    event_id UUID PRIMARY KEY,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    actor_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    request_id TEXT NOT NULL,
    ip_hash TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    prev_hash BYTEA,
    event_hash BYTEA NOT NULL,
    CHECK (octet_length(event_hash) = 32)
);
CREATE INDEX IF NOT EXISTS audit_log_time_idx ON audit_log(occurred_at);
CREATE INDEX IF NOT EXISTS audit_log_actor_idx ON audit_log(actor_id);

CREATE TABLE IF NOT EXISTS reconciliation_runs (
    run_id UUID PRIMARY KEY,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ,
    status TEXT NOT NULL CHECK (status IN ('running','passed','failed')),
    discrepancy_count BIGINT NOT NULL DEFAULT 0,
    report JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS dr_backup_verifications (
    verification_id UUID PRIMARY KEY,
    backup_reference TEXT NOT NULL,
    verified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    restore_tested BOOLEAN NOT NULL,
    checksum TEXT,
    notes TEXT
);
