package txcodec

import("math/big";"testing")
func TestEIP1559Digest(t *testing.T){d,err:=EIP1559Digest(EVM1559Tx{ChainID:big.NewInt(1),Nonce:big.NewInt(1),MaxPriorityFeePerGas:big.NewInt(1000000000),MaxFeePerGas:big.NewInt(2000000000),GasLimit:big.NewInt(21000),To:make([]byte,20),Value:big.NewInt(1),Data:nil});if err!=nil{t.Fatal(err)};if len(d)!=32{t.Fatalf("digest length=%d",len(d))}}
func TestLegacyEIP155Digest(t *testing.T){d,err:=LegacyEIP155Digest(EVMLegacyTx{Nonce:big.NewInt(1),GasPrice:big.NewInt(1),GasLimit:big.NewInt(21000),To:make([]byte,20),Value:big.NewInt(1),ChainID:big.NewInt(1)});if err!=nil{t.Fatal(err)};if len(d)!=32{t.Fatalf("digest length=%d",len(d))}}
