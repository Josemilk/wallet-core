CREATE TABLE IF NOT EXISTS fee_policies (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  operation_type TEXT NOT NULL UNIQUE,
  basis_points INTEGER NOT NULL CHECK (basis_points BETWEEN 0 AND 10000),
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  updated_by UUID REFERENCES users(id),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
INSERT INTO fee_policies(operation_type,basis_points) VALUES
 ('TRADING',1000),
 ('SWAP',1000),
 ('WITHDRAWAL',1000),
 ('EVENT_ENTRY',1000)
ON CONFLICT(operation_type) DO NOTHING;

CREATE TABLE IF NOT EXISTS treasury_accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id UUID NOT NULL UNIQUE REFERENCES accounts(id),
  asset TEXT NOT NULL,
  purpose TEXT NOT NULL CHECK (purpose IN ('FEE_REVENUE','POOL_LIQUIDITY','PRIZE_FORFEITURE')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS treasury_transfer_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  treasury_account_id UUID NOT NULL REFERENCES treasury_accounts(id),
  asset TEXT NOT NULL,
  amount NUMERIC(78,0) NOT NULL CHECK (amount > 0),
  destination TEXT NOT NULL,
  network TEXT NOT NULL,
  requested_by UUID NOT NULL REFERENCES users(id),
  approved_by UUID REFERENCES users(id),
  status TEXT NOT NULL CHECK (status IN ('requested','risk_review','approved','signing','broadcast','confirmed','rejected','cancelled')),
  idempotency_key TEXT NOT NULL UNIQUE,
  custody_key_ref TEXT,
  tx_hash TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (approved_by IS NULL OR approved_by <> requested_by)
);
CREATE INDEX IF NOT EXISTS treasury_transfer_ready_idx ON treasury_transfer_requests(status,created_at);

CREATE TABLE IF NOT EXISTS event_prize_escrows (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL UNIQUE REFERENCES events(id),
  account_id UUID NOT NULL REFERENCES accounts(id),
  asset TEXT NOT NULL,
  reserved_amount NUMERIC(78,0) NOT NULL CHECK (reserved_amount >= 0),
  claimed_amount NUMERIC(78,0) NOT NULL DEFAULT 0 CHECK (claimed_amount >= 0),
  forfeited_amount NUMERIC(78,0) NOT NULL DEFAULT 0 CHECK (forfeited_amount >= 0),
  claim_deadline TIMESTAMPTZ,
  status TEXT NOT NULL CHECK (status IN ('funding','funded','claimable','partially_claimed','claimed','forfeited','cancelled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (claimed_amount + forfeited_amount <= reserved_amount)
);

CREATE TABLE IF NOT EXISTS event_prize_claims (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  escrow_id UUID NOT NULL REFERENCES event_prize_escrows(id),
  event_entry_id UUID NOT NULL UNIQUE REFERENCES event_entries(id),
  user_id UUID NOT NULL REFERENCES users(id),
  amount NUMERIC(78,0) NOT NULL CHECK (amount > 0),
  status TEXT NOT NULL CHECK (status IN ('pending','paid','rejected','expired')),
  ledger_transaction_id UUID UNIQUE REFERENCES ledger_transactions(id),
  idempotency_key TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS fee_ledger_links (
  ledger_transaction_id UUID PRIMARY KEY REFERENCES ledger_transactions(id),
  operation_type TEXT NOT NULL,
  basis_points INTEGER NOT NULL,
  source_amount NUMERIC(78,0) NOT NULL CHECK (source_amount > 0),
  fee_amount NUMERIC(78,0) NOT NULL CHECK (fee_amount >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
