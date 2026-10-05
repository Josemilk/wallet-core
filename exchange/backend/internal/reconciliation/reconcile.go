package reconciliation

import "math/big"

func Balanced(amounts []*big.Int) bool { total:=new(big.Int); for _,a:=range amounts { if a==nil{return false}; total.Add(total,a) }; return total.Sign()==0 }
