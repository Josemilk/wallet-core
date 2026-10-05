CREATE TABLE IF NOT EXISTS balances (
    account_id UUID NOT NULL REFERENCES accounts(id),
    asset TEXT NOT NULL,
    available NUMERIC(78,0) NOT NULL DEFAULT 0 CHECK (available >= 0),
    locked NUMERIC(78,0) NOT NULL DEFAULT 0 CHECK (locked >= 0),
    version BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (account_id, asset)
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES accounts(id),
    market TEXT NOT NULL,
    side TEXT NOT NULL CHECK (side IN ('buy','sell')),
    type TEXT NOT NULL CHECK (type IN ('limit','market')),
    price NUMERIC(78,0),
    quantity NUMERIC(78,0) NOT NULL CHECK (quantity > 0),
    filled_quantity NUMERIC(78,0) NOT NULL DEFAULT 0 CHECK (filled_quantity >= 0),
    status TEXT NOT NULL CHECK (status IN ('open','partial','filled','cancelled','rejected')),
    client_order_id TEXT,
    sequence BIGSERIAL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (filled_quantity <= quantity)
);
CREATE UNIQUE INDEX IF NOT EXISTS orders_client_id ON orders(account_id, client_order_id) WHERE client_order_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS orders_book_idx ON orders(market, side, status, price, sequence);

CREATE TABLE IF NOT EXISTS market_ticks (
    market TEXT NOT NULL,
    sequence BIGSERIAL,
    price NUMERIC(78,0) NOT NULL,
    quantity NUMERIC(78,0) NOT NULL,
    event_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (market, sequence)
);

CREATE TABLE IF NOT EXISTS market_candles (
    market TEXT NOT NULL,
    interval_seconds INTEGER NOT NULL CHECK (interval_seconds > 0),
    bucket_start TIMESTAMPTZ NOT NULL,
    open NUMERIC(78,0) NOT NULL,
    high NUMERIC(78,0) NOT NULL,
    low NUMERIC(78,0) NOT NULL,
    close NUMERIC(78,0) NOT NULL,
    volume NUMERIC(78,0) NOT NULL DEFAULT 0,
    PRIMARY KEY (market, interval_seconds, bucket_start)
);
