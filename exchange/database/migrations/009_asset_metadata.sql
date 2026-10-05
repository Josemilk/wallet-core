CREATE TABLE IF NOT EXISTS asset_metadata (
  asset TEXT PRIMARY KEY,
  decimals INTEGER NOT NULL CHECK (decimals BETWEEN 0 AND 38),
  market_data_id TEXT,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
INSERT INTO asset_metadata(asset,decimals,market_data_id) VALUES
 ('BTC',8,'1'),('ETH',18,'1027'),('SOL',9,'5426'),('USDT',6,'825')
ON CONFLICT(asset) DO NOTHING;
