package blockchain

import("context";"errors";"time")
type Cursor interface{Load(context.Context,string)(uint64,string,error);Save(context.Context,string,uint64,string)error}
type RangeScanner interface{ScanRange(context.Context,uint64,uint64,DepositSink)error}
type BlockScanner interface{ScanBlock(context.Context,uint64,uint64,DepositSink)error}
type ReorgHandler interface{RollbackBlock(context.Context,string,string)error}
type ConfirmationRefresher interface{Refresh(context.Context,string,uint64)error}

type PersistentWatcher struct{Network string;Head ChainHeadSource;Reader ChainBlockReader;Cursor Cursor;EVM RangeScanner;Bitcoin BlockScanner;Sink DepositSink;Reorg ReorgHandler;Refresh ConfirmationRefresher;PollInterval time.Duration}
func(w *PersistentWatcher)Run(ctx context.Context)error{if w==nil||w.Head==nil||w.Reader==nil||w.Cursor==nil||w.Sink==nil{return errors.New("persistent watcher dependencies missing")};if w.Network==""{return errors.New("network required")};interval:=w.PollInterval;if interval<=0{interval=5*time.Second};ticker:=time.NewTicker(interval);defer ticker.Stop();for{if err:=w.tick(ctx);err!=nil{return err};select{case<-ctx.Done():return ctx.Err();case<-ticker.C:}}}
func(w *PersistentWatcher)tick(ctx context.Context)error{last,lastHash,err:=w.Cursor.Load(ctx,w.Network);if err!=nil{return err};head,err:=w.Head.LatestChainBlock(ctx);if err!=nil{return err};if lastHash!=""&&last>0{next,err:=w.Reader.ChainBlock(ctx,last+1);if err==nil&&next.ParentHash!=lastHash{if w.Reorg!=nil{if err:=w.Reorg.RollbackBlock(ctx,w.Network,lastHash);err!=nil{return err}};last--;if last>0{b,err:=w.Reader.ChainBlock(ctx,last);if err!=nil{return err};lastHash=b.Hash}else{lastHash=""}}};from:=last+1;if from<=head.Height{if w.EVM!=nil{if err:=w.EVM.ScanRange(ctx,from,head.Height,w.Sink);err!=nil{return err}};if w.Bitcoin!=nil{for h:=from;h<=head.Height;h++{if err:=w.Bitcoin.ScanBlock(ctx,h,head.Height,w.Sink);err!=nil{return err}}}};if w.Refresh!=nil{if err:=w.Refresh.Refresh(ctx,w.Network,head.Height);err!=nil{return err}};return w.Cursor.Save(ctx,w.Network,head.Height,head.Hash)}
