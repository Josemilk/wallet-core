package blockchain

import("context";"testing")

type testSource struct{head ChainBlock}
func(s *testSource)LatestChainBlock(context.Context)(ChainBlock,error){return s.head,nil}
type testScanner struct{items []DepositCandidate}
func(s *testScanner)ScanDeposits(context.Context,ChainBlock)([]DepositCandidate,error){return s.items,nil}
type testConfirmed struct{count int}
func(s *testConfirmed)ConfirmDeposit(context.Context,DepositCandidate,uint64)error{s.count++;return nil}

func TestConfirmationEngineConfirmsAfterThreshold(t *testing.T){src:=&testSource{head:ChainBlock{Height:10,Hash:"h10"}};sc:=&testScanner{items:[]DepositCandidate{{TxHash:"tx1",Address:"a",Asset:"BTC",Amount:"1",BlockHash:"h9",BlockHeight:9}}};sink:=&testConfirmed{count:0};e:=&ConfirmationEngine{Source:src,Scanner:sc,Confirmed:sink,Confirmations:2,pending:map[string]DepositCandidate{}};if err:=e.scan(context.Background(),src.head);err!=nil{t.Fatal(err)};if err:=e.confirm(context.Background(),10);err!=nil{t.Fatal(err)};if sink.count!=1{t.Fatalf("expected one confirmation, got %d",sink.count)}}
