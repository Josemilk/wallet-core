package reorg

import "sync"

type Block struct { Height uint64; Hash string; ParentHash string }
type Tracker struct { mu sync.Mutex; byHeight map[uint64]Block }
func NewTracker() *Tracker { return &Tracker{byHeight:map[uint64]Block{}} }
func (t *Tracker) Add(b Block) (reorg bool, replaced Block) { t.mu.Lock(); defer t.mu.Unlock(); old,ok:=t.byHeight[b.Height]; if ok && old.Hash!=b.Hash { t.byHeight[b.Height]=b; return true,old }; t.byHeight[b.Height]=b; return false,Block{} }
func (t *Tracker) Get(height uint64) (Block,bool) { t.mu.Lock(); defer t.mu.Unlock(); b,ok:=t.byHeight[height]; return b,ok }
