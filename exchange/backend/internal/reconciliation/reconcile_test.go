package reconciliation

import (
 "math/big"
 "testing"
)

func TestCompareFindsDiscrepancy(t *testing.T) {
 d:=Compare([]Position{{Asset:"BTC",Ledger:big.NewInt(100),Custody:big.NewInt(99)}})
 if len(d)!=1 || d[0].Difference.Cmp(big.NewInt(1))!=0 {t.Fatalf("unexpected discrepancy: %#v",d)}
}
func TestCompareClean(t *testing.T) {
 if got:=Compare([]Position{{Asset:"BTC",Ledger:big.NewInt(100),Custody:big.NewInt(100)}}); len(got)!=0 {t.Fatal(got)}
}
