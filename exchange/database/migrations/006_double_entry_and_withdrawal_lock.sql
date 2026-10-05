ALTER TABLE accounts ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS accounts_system_owner_check;
ALTER TABLE accounts ADD CONSTRAINT accounts_system_owner_check CHECK ((kind = 'SYSTEM' AND user_id IS NULL) OR (kind <> 'SYSTEM' AND user_id IS NOT NULL));
CREATE UNIQUE INDEX IF NOT EXISTS accounts_system_asset_unique ON accounts(asset) WHERE kind='SYSTEM';

CREATE OR REPLACE FUNCTION enforce_ledger_balance() RETURNS trigger AS $$
DECLARE txid UUID; total NUMERIC(78,0);
BEGIN
  txid := COALESCE(NEW.transaction_id, OLD.transaction_id);
  SELECT COALESCE(SUM(amount),0) INTO total FROM ledger_entries WHERE transaction_id = txid;
  IF total <> 0 THEN RAISE EXCEPTION 'ledger transaction % is unbalanced', txid; END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS withdrawal_ledger_links (
  withdrawal_id UUID PRIMARY KEY REFERENCES withdrawals(id),
  ledger_transaction_id UUID NOT NULL UNIQUE REFERENCES ledger_transactions(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS withdrawal_queue_ledger_links (
  withdrawal_queue_id UUID PRIMARY KEY REFERENCES withdrawal_queue(id),
  ledger_transaction_id UUID NOT NULL UNIQUE REFERENCES ledger_transactions(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS withdrawal_queue_claim_idx ON withdrawal_queue(status, available_after, created_at, id);
