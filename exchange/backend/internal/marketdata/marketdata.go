package marketdata

type Tick struct { Symbol string `json:"symbol"`; Bid, Ask, Last int64 `json:"bid"`; Sequence uint64 `json:"sequence"` }
type Candle struct { Symbol, Interval string; Open, High, Low, Close, Volume int64; OpenTime, CloseTime int64 }
