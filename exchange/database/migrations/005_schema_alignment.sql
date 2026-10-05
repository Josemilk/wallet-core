-- Reconcile the original exchange schema with the later service contracts.
ALTER TABLE orders ADD COLUMN IF NOT EXISTS account_id UUID REFERENCES accounts(id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS market TEXT;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS type TEXT;
CREATE INDEX IF NOT EXISTS orders_account_id_idx ON orders(account_id);
CREATE INDEX IF NOT EXISTS orders_market_status_idx ON orders(market, status);

ALTER TABLE deposits ADD COLUMN IF NOT EXISTS account_id UUID REFERENCES accounts(id);
ALTER TABLE deposits ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE deposits ADD COLUMN IF NOT EXISTS ledger_transaction_id UUID;
CREATE UNIQUE INDEX IF NOT EXISTS deposits_ledger_transaction_unique ON deposits(ledger_transaction_id) WHERE ledger_transaction_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS deposits_account_status_idx ON deposits(account_id, status);

-- The UI/backend contract uses explicit state names; existing rows are not rewritten here.
-- Backfills must be performed by a controlled migration job once production account mapping exists.
