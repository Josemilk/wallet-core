CREATE TABLE IF NOT EXISTS settlement_records (
    trade_id UUID PRIMARY KEY,
    status TEXT NOT NULL CHECK (status IN ('pending','settled','failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    settled_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS withdrawal_state (
    withdrawal_id UUID PRIMARY KEY,
    status TEXT NOT NULL CHECK (status IN ('requested','risk_approved','signing','broadcast','confirmed','failed','cancelled')),
    tx_hash TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_withdrawal_state_status ON withdrawal_state(status);
