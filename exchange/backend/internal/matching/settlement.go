package matching

import "errors"

type Fill struct { TradeID, BuyOrderID, SellOrderID string; Price, Quantity int64 }
type TradeSink interface { Settle(Fill) error }

func ValidateFill(f Fill) error { if f.TradeID==""||f.BuyOrderID==""||f.SellOrderID==""||f.Price<=0||f.Quantity<=0{return errors.New("invalid fill")};return nil }
