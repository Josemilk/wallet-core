ALTER TABLE deposits ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE deposits ADD COLUMN IF NOT EXISTS contract_address TEXT;
ALTER TABLE deposits DROP CONSTRAINT IF EXISTS deposits_network_tx_hash_key;
CREATE UNIQUE INDEX IF NOT EXISTS deposits_recipient_unique_idx ON deposits(network,tx_hash,address,asset,COALESCE(contract_address,''));
CREATE INDEX IF NOT EXISTS deposits_user_created_idx ON deposits(user_id,created_at DESC);
