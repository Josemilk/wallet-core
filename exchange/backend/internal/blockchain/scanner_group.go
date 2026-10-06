package blockchain

import "context"

type RangeScannerGroup []RangeScanner
func(g RangeScannerGroup)ScanRange(ctx context.Context,from,to uint64,sink DepositSink)error{for _,s:=range g{if s==nil{continue};if err:=s.ScanRange(ctx,from,to,sink);err!=nil{return err}};return nil}
