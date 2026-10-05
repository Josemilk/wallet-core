package blockchain

import("context";"errors";"sync";"time")

type ChainBlock struct{Height uint64;Hash string;ParentHash string}
type ChainHeadSource interface{LatestChainBlock(context.Context)(ChainBlock,error)}
type ChainBlockReader interface{ChainBlock(ctx context.Context,height uint64)(ChainBlock,error)}
type DepositCandidate struct{TxHash,Address,Asset,Amount,BlockHash string;BlockHeight uint64}
type DepositScanner interface{ScanDeposits(context.Context,ChainBlock)([]DepositCandidate,error)}
type ConfirmationSink interface{ConfirmDeposit(context.Context,DepositCandidate,uint64)error}
type ReorgSink interface{RollbackDeposit(context.Context,string)error}

type ConfirmationEngine struct{Source ChainHeadSource;Reader ChainBlockReader;Scanner DepositScanner;Confirmed ConfirmationSink;Reorg ReorgSink;Confirmations uint64;PollInterval time.Duration;mu sync.Mutex;last ChainBlock;pending map[string]DepositCandidate}
func(e *ConfirmationEngine)Run(ctx context.Context)error{if e==nil||e.Source==nil||e.Scanner==nil||e.Confirmed==nil{return errors.New("confirmation engine dependencies missing")};if e.Confirmations==0{return errors.New("confirmation threshold must be positive")};if e.pending==nil{e.pending=map[string]DepositCandidate{}};interval:=e.PollInterval;if interval<=0{interval=5*time.Second};ticker:=time.NewTicker(interval);defer ticker.Stop();for{if err:=e.tick(ctx);err!=nil{return err};select{case<-ctx.Done():return ctx.Err();case<-ticker.C:}}}
func(e *ConfirmationEngine)tick(ctx context.Context)error{head,err:=e.Source.LatestChainBlock(ctx);if err!=nil{return err};e.mu.Lock();prev:=e.last;e.last=head;e.mu.Unlock();if prev.Hash!=""&&head.ParentHash!=""&&head.ParentHash!=prev.Hash{e.mu.Lock();for hash,d:=range e.pending{if d.BlockHash==prev.Hash&&e.Reorg!=nil{_ = e.Reorg.RollbackDeposit(ctx,hash);delete(e.pending,hash)}};e.mu.Unlock()};if e.Reader!=nil&&head.Height>prev.Height+1{for h:=prev.Height+1;h<=head.Height;h++{b,err:=e.Reader.ChainBlock(ctx,h);if err!=nil{return err};if err:=e.scan(ctx,b);err!=nil{return err}}}else{if err:=e.scan(ctx,head);err!=nil{return err}};return e.confirm(ctx,head.Height)}
func(e *ConfirmationEngine)scan(ctx context.Context,b ChainBlock)error{items,err:=e.Scanner.ScanDeposits(ctx,b);if err!=nil{return err};e.mu.Lock();defer e.mu.Unlock();for _,d:=range items{if d.TxHash==""||d.BlockHeight==0{continue};if _,ok:=e.pending[d.TxHash];!ok{e.pending[d.TxHash]=d}};return nil}
func(e *ConfirmationEngine)confirm(ctx context.Context,height uint64)error{e.mu.Lock();ready:=make([]DepositCandidate,0);for hash,d:=range e.pending{if height>=d.BlockHeight+e.Confirmations-1{ready=append(ready,d);delete(e.pending,hash)}};e.mu.Unlock();for _,d:=range ready{if err:=e.Confirmed.ConfirmDeposit(ctx,d,height-d.BlockHeight+1);err!=nil{return err}};return nil}
