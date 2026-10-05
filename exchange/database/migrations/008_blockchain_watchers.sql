CREATE TABLE IF NOT EXISTS deposit_addresses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  account_id UUID NOT NULL REFERENCES accounts(id),
  asset TEXT NOT NULL,
  network TEXT NOT NULL,
  address TEXT NOT NULL,
  contract_address TEXT,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS deposit_addresses_unique_idx
  ON deposit_addresses(network,address,asset,COALESCE(contract_address,''));

CREATE TABLE IF NOT EXISTS chain_cursors (
  network TEXT PRIMARY KEY,
  last_height BIGINT NOT NULL DEFAULT 0,
  last_hash TEXT,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS observed_deposits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  network TEXT NOT NULL,
  tx_hash TEXT NOT NULL,
  block_height BIGINT NOT NULL,
  block_hash TEXT NOT NULL,
  address TEXT NOT NULL,
  asset TEXT NOT NULL,
  contract_address TEXT,
  amount NUMERIC(78,0) NOT NULL CHECK (amount > 0),
  confirmations BIGINT NOT NULL DEFAULT 0,
  status TEXT NOT NULL CHECK (status IN ('observed','confirming','credited','reorged','ignored')),
  ledger_transaction_id UUID REFERENCES ledger_transactions(id),
  first_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS observed_deposits_unique_idx
  ON observed_deposits(network,tx_hash,address,asset,COALESCE(contract_address,''));
CREATE INDEX IF NOT EXISTS observed_deposits_confirmation_idx ON observed_deposits(network,status,block_height);
CREATE INDEX IF NOT EXISTS observed_deposits_tx_idx ON observed_deposits(network,tx_hash);
