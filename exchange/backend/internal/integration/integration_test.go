package integration_test

import (
    "testing"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/reorg"
)

func TestReorgTrackerReplacesCompetingBlock(t *testing.T) {
    tr:=reorg.NewTracker()
    if changed,_:=tr.Add(reorg.Block{Height:10,Hash:"a",ParentHash:"p"}); changed { t.Fatal("first block cannot be a reorg") }
    if changed,old:=tr.Add(reorg.Block{Height:10,Hash:"b",ParentHash:"p"}); !changed || old.Hash!="a" { t.Fatal("replacement block was not detected") }
}
