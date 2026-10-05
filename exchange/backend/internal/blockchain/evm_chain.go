package blockchain

import("context";"fmt";"strconv";"strings")

type EVMChainReader struct{RPC *EVMRPC}
func(r *EVMChainReader)LatestChainBlock(ctx context.Context)(ChainBlock,error){if r==nil||r.RPC==nil{return ChainBlock{},fmt.Errorf("EVM RPC unavailable")};h,err:=r.RPC.LatestBlock(ctx);if err!=nil{return ChainBlock{},err};return r.ChainBlock(ctx,h.Height)}
func(r *EVMChainReader)ChainBlock(ctx context.Context,height uint64)(ChainBlock,error){b,err:=r.RPC.BlockByNumber(ctx,height);if err!=nil{return ChainBlock{},err};hash,_:=b["hash"].(string);parent,_:=b["parentHash"].(string);return ChainBlock{Height:height,Hash:hash,ParentHash:parent},nil}
func parseEVMQuantity(v any)uint64{switch x:=v.(type){case string:n,_:=strconv.ParseUint(strings.TrimPrefix(x,"0x"),16,64);return n;default:return 0}}
