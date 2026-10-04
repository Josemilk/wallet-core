package ledger

import (
    "context"
    "errors"
    "sync"
)

type Entry struct { Account, Asset string; Amount int64 }
type Journal struct { ID, IdempotencyKey string; Entries []Entry }

type Ledger struct { mu sync.RWMutex; journals map[string]Journal }
func New() *Ledger { return &Ledger{journals: make(map[string]Journal)} }
func (l *Ledger) Ready(context.Context) error { return nil }
func (l *Ledger) Post(j Journal) error {
    if j.ID == "" || j.IdempotencyKey == "" || len(j.Entries) < 2 { return errors.New("invalid journal") }
    sums := map[string]int64{}
    for _, e := range j.Entries { if e.Account == "" || e.Asset == "" || e.Amount == 0 { return errors.New("invalid entry") }; sums[e.Asset] += e.Amount }
    for _, v := range sums { if v != 0 { return errors.New("unbalanced journal") } }
    l.mu.Lock(); defer l.mu.Unlock()
    if _, exists := l.journals[j.IdempotencyKey]; exists { return nil }
    l.journals[j.IdempotencyKey] = j
    return nil
}
