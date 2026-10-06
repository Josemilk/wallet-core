package txcodec

import "testing"

func TestBIP143DigestDeterministic(t *testing.T){tx:=BitcoinTx{Version:2,SighashType:1,Inputs:[]BitcoinInput{{TxID:"0000000000000000000000000000000000000000000000000000000000000000",Vout:0,ScriptCode:[]byte{0x19,0x76,0xa9,0x14,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0x88,0xac},Amount:100000,Sequence:0xffffffff}},Outputs:[]BitcoinOutput{{Value:90000,ScriptPubKey:[]byte{0x00,0x14,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}},LockTime:0};a,err:=BIP143Digest(tx,0);if err!=nil{t.Fatal(err)};b,err:=BIP143Digest(tx,0);if err!=nil{t.Fatal(err)};if string(a)!=string(b){t.Fatal("digest is not deterministic")};if len(a)!=32{t.Fatalf("digest length=%d",len(a))}}
