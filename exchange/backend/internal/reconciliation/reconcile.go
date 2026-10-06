package reconciliation

import (
    "math/big"
)

type Position struct {
    Asset string
    Ledger *big.Int
    Custody *big.Int
}

type Discrepancy struct {
    Asset string
    Ledger *big.Int
    Custody *big.Int
    Difference *big.Int
}

func Compare(positions []Position) []Discrepancy {
    out := make([]Discrepancy, 0)
    for _, p := range positions {
        if p.Ledger == nil || p.Custody == nil { continue }
        d := new(big.Int).Sub(p.Ledger, p.Custody)
        if d.Sign() != 0 {
            out = append(out, Discrepancy{
                Asset:p.Asset,
                Ledger:new(big.Int).Set(p.Ledger),
                Custody:new(big.Int).Set(p.Custody),
                Difference:d,
            })
        }
    }
    return out
}

func Balanced(amounts []*big.Int) bool {
    total:=new(big.Int)
    for _,a:=range amounts { if a==nil{return false}; total.Add(total,a) }
    return total.Sign()==0
}
