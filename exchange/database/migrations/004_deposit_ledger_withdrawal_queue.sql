ALTER TABLE deposits ADD COLUMN IF NOT EXISTS ledger_transaction_id UUID;
CREATE UNIQUE INDEX IF NOT EXISTS deposits_ledger_tx_unique ON deposits(ledger_transaction_id) WHERE ledger_transaction_id IS NOT NULL;
CREATE TABLE IF NOT EXISTS withdrawal_queue (
 id UUID PRIMARY KEY,
 user_id UUID NOT NULL REFERENCES users(id),
 account_id UUID NOT NULL REFERENCES accounts(id),
 asset TEXT NOT NULL,
 amount NUMERIC(78,0) NOT NULL CHECK (amount > 0),
 destination TEXT NOT NULL,
 network TEXT NOT NULL,
 key_ref TEXT NOT NULL,
 status TEXT NOT NULL CHECK (status IN ('queued','risk_approved','signing','broadcast','confirmed','failed','cancelled')),
 tx_hash TEXT,
 attempt_count INTEGER NOT NULL DEFAULT 0,
 available_after TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS withdrawal_queue_ready ON withdrawal_queue(status, available_after, created_at);
