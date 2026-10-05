package blockchain

import (
    "context"
    "errors"
    "time"
)

type Block struct { Height uint64; Hash string }
type Watcher interface { Run(context.Context) error }
type BlockSource interface { LatestBlock(context.Context) (Block,error) }
type DepositObserver interface { OnConfirmedDeposit(ctx context.Context, txHash, address, asset, amount string, confirmations uint64) error }

type ConfirmingWatcher struct {
    Source BlockSource
    Observer DepositObserver
    Confirmations uint64
    PollInterval time.Duration
}

func (w *ConfirmingWatcher) Run(ctx context.Context) error {
    if w == nil || w.Source == nil || w.Observer == nil { return errors.New("watcher dependencies missing") }
    if w.Confirmations == 0 { return errors.New("confirmation threshold must be positive") }
    interval := w.PollInterval; if interval <= 0 { interval = 5*time.Second }
    ticker := time.NewTicker(interval); defer ticker.Stop()
    for {
        if err := w.tick(ctx); err != nil { return err }
        select { case <-ctx.Done(): return ctx.Err(); case <-ticker.C: }
    }
}

func (w *ConfirmingWatcher) tick(ctx context.Context) error {
    _, err := w.Source.LatestBlock(ctx)
    return err
}
