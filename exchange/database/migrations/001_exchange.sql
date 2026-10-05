CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','SUSPENDED','CLOSED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  asset TEXT NOT NULL,
  kind TEXT NOT NULL CHECK (kind IN ('AVAILABLE','LOCKED','FEE','SYSTEM')),
  UNIQUE(user_id, asset, kind)
);

CREATE TABLE IF NOT EXISTS ledger_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  idempotency_key TEXT UNIQUE NOT NULL,
  type TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ledger_entries (
  id BIGSERIAL PRIMARY KEY,
  transaction_id UUID NOT NULL REFERENCES ledger_transactions(id),
  account_id UUID NOT NULL REFERENCES accounts(id),
  asset TEXT NOT NULL,
  amount NUMERIC(78,0) NOT NULL CHECK (amount <> 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS ledger_entries_tx_idx ON ledger_entries(transaction_id);
CREATE INDEX IF NOT EXISTS ledger_entries_account_asset_idx ON ledger_entries(account_id, asset);

CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  symbol TEXT NOT NULL,
  side TEXT NOT NULL CHECK (side IN ('BUY','SELL')),
  order_type TEXT NOT NULL CHECK (order_type IN ('LIMIT','MARKET')),
  price NUMERIC(78,0),
  quantity NUMERIC(78,0) NOT NULL CHECK (quantity > 0),
  remaining NUMERIC(78,0) NOT NULL CHECK (remaining >= 0),
  status TEXT NOT NULL,
  client_order_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(user_id, client_order_id)
);

CREATE TABLE IF NOT EXISTS trades (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  symbol TEXT NOT NULL,
  maker_order_id UUID NOT NULL REFERENCES orders(id),
  taker_order_id UUID NOT NULL REFERENCES orders(id),
  price NUMERIC(78,0) NOT NULL,
  quantity NUMERIC(78,0) NOT NULL CHECK (quantity > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deposits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  asset TEXT NOT NULL,
  network TEXT NOT NULL,
  tx_hash TEXT NOT NULL,
  amount NUMERIC(78,0) NOT NULL,
  confirmations BIGINT NOT NULL DEFAULT 0,
  status TEXT NOT NULL,
  UNIQUE(network, tx_hash)
);

CREATE TABLE IF NOT EXISTS withdrawals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  asset TEXT NOT NULL,
  network TEXT NOT NULL,
  destination TEXT NOT NULL,
  amount NUMERIC(78,0) NOT NULL,
  status TEXT NOT NULL,
  idempotency_key TEXT UNIQUE NOT NULL,
  tx_hash TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS risk_events (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID,
  category TEXT NOT NULL,
  decision TEXT NOT NULL,
  reason TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION enforce_ledger_balance() RETURNS trigger AS $$
DECLARE total NUMERIC(78,0);
BEGIN
  SELECT COALESCE(SUM(amount),0) INTO total FROM ledger_entries WHERE transaction_id = NEW.transaction_id;
  IF total <> 0 THEN RAISE EXCEPTION 'ledger transaction % is unbalanced', NEW.transaction_id; END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS ledger_balance_deferred ON ledger_entries;
CREATE CONSTRAINT TRIGGER ledger_balance_deferred
AFTER INSERT OR UPDATE OR DELETE ON ledger_entries
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW EXECUTE FUNCTION enforce_ledger_balance();
