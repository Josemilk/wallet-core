package blockchain

import("context";"fmt";"strconv")

type BitcoinChainReader struct{RPC *BitcoinRPC}
func(r *BitcoinChainReader)LatestChainBlock(ctx context.Context)(ChainBlock,error){if r==nil||r.RPC==nil{return ChainBlock{},fmt.Errorf("Bitcoin RPC unavailable")};var info map[string]any;if err:=r.RPC.call(ctx,"getblockchaininfo",[]any{},&info);err!=nil{return ChainBlock{},err};h,_:=info["blocks"].(float64);return r.ChainBlock(ctx,uint64(h))}
func(r *BitcoinChainReader)ChainBlock(ctx context.Context,height uint64)(ChainBlock,error){hash,err:=r.RPC.BlockHash(ctx,height);if err!=nil{return ChainBlock{},err};b,err:=r.RPC.Block(ctx,hash);if err!=nil{return ChainBlock{},err};parent,_:=b["previousblockhash"].(string);if v,ok:=b["height"].(float64);ok{height=uint64(v)};return ChainBlock{Height:height,Hash:hash,ParentHash:parent},nil}
func bitcoinHeight(v any)uint64{switch x:=v.(type){case float64:return uint64(x);case string:n,_:=strconv.ParseUint(x,10,64);return n;default:return 0}}
