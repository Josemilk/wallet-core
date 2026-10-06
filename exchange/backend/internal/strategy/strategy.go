package strategy

type Signal struct { Symbol string; Side string; ConfidenceBps int64; Price int64; Quantity int64 }
type BacktestResult struct { Trades int; NetPnl int64; MaxDrawdown int64; SharpeBps int64 }
type Engine interface { Signal(symbol string, candles any) (Signal,error) }

// AI models produce proposals only. Live execution must pass the same risk, limits,
// balance and authorization gates as manual orders.
