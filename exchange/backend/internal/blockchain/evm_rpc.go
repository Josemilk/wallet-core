package blockchain

import("bytes";"context";"encoding/hex";"encoding/json";"errors";"fmt";"io";"net/http";"strings";"time")

type EVMRPC struct{URL,Token string;HTTP *http.Client}
type rpcReq struct{JSONRPC string `json:"jsonrpc"`;ID uint64 `json:"id"`;Method string `json:"method"`;Params any `json:"params"`}
type rpcResp struct{JSONRPC string `json:"jsonrpc"`;ID uint64 `json:"id"`;Result json.RawMessage `json:"result"`;Error *struct{Code int `json:"code"`;Message string `json:"message"`} `json:"error"`}
func(n *EVMRPC)call(ctx context.Context,method string,params any,out any)error{if n==nil||n.URL==""{return errors.New("EVM RPC URL is required")};body,_:=json.Marshal(rpcReq{"2.0",1,method,params});req,err:=http.NewRequestWithContext(ctx,http.MethodPost,n.URL,bytes.NewReader(body));if err!=nil{return err};req.Header.Set("Content-Type","application/json");if n.Token!=""{req.Header.Set("Authorization","Bearer "+n.Token)};c:=n.HTTP;if c==nil{c=&http.Client{Timeout:15*time.Second}};resp,err:=c.Do(req);if err!=nil{return err};defer resp.Body.Close();raw,err:=io.ReadAll(io.LimitReader(resp.Body,8<<20));if err!=nil{return err};if resp.StatusCode/100!=2{return fmt.Errorf("EVM RPC HTTP %d",resp.StatusCode)};var r rpcResp;if err:=json.Unmarshal(raw,&r);err!=nil{return err};if r.Error!=nil{return fmt.Errorf("EVM RPC %d: %s",r.Error.Code,r.Error.Message)};if out==nil{return nil};return json.Unmarshal(r.Result,out)}
func(n *EVMRPC)LatestBlock(ctx context.Context)(Block,error){var hexHeight string;if err:=n.call(ctx,"eth_blockNumber",[]any{},&hexHeight);err!=nil{return Block{},err};v,err:=parseHexUint(hexHeight);if err!=nil{return Block{},err};return Block{Height:v},nil}
func(n *EVMRPC)BlockByNumber(ctx context.Context,height uint64)(map[string]any,error){var out map[string]any;if err:=n.call(ctx,"eth_getBlockByNumber",[]any{fmt.Sprintf("0x%x",height),true},&out);err!=nil{return nil,err};return out,nil}
func(n *EVMRPC)Transaction(ctx context.Context,hash string)(map[string]any,error){var out map[string]any;if err:=n.call(ctx,"eth_getTransactionByHash",[]any{hash},&out);err!=nil{return nil,err};return out,nil}
func(n *EVMRPC)Logs(ctx context.Context,from,to uint64,address string,topics []string)([]map[string]any,error){filter:=map[string]any{"fromBlock":fmt.Sprintf("0x%x",from),"toBlock":fmt.Sprintf("0x%x",to)};if address!=""{filter["address"]=address};if len(topics)>0{filter["topics"]=topics};var out []map[string]any;if err:=n.call(ctx,"eth_getLogs",[]any{filter},&out);err!=nil{return nil,err};return out,nil}
func(n *EVMRPC)SendRawTransaction(ctx context.Context,raw string)(string,error){raw=strings.TrimSpace(raw);if !strings.HasPrefix(raw,"0x"){raw="0x"+raw};body:=strings.TrimPrefix(raw,"0x");if body==""||len(body)%2!=0{return "",errors.New("invalid EVM raw transaction hex")};if _,err:=hex.DecodeString(body);err!=nil{return "",errors.New("invalid EVM raw transaction hex")};var txHash string;if err:=n.call(ctx,"eth_sendRawTransaction",[]any{raw},&txHash);err!=nil{return "",err};return txHash,nil}
func parseHexUint(v string)(uint64,error){v=strings.TrimPrefix(v,"0x");if v==""{return 0,nil};b,err:=hex.DecodeString(func()string{if len(v)%2==1{return "0"+v};return v}());if err!=nil{return 0,err};var n uint64;if len(b)>8{return 0,errors.New("hex integer overflow")};for _,x:=range b{n=n<<8|uint64(x)};return n,nil}

func(n *EVMRPC) TransactionReceipt(ctx context.Context, hash string) (map[string]any, error) {
    hash=strings.TrimSpace(hash)
    if hash=="" { return nil, errors.New("transaction hash is required") }
    var out map[string]any
    if err:=n.call(ctx,"eth_getTransactionReceipt",[]any{hash},&out); err!=nil { return nil,err }
    return out,nil
}

func(n *EVMRPC) WaitForReceipt(ctx context.Context, hash string, interval time.Duration) (map[string]any,error) {
    if interval<=0 { interval=3*time.Second }
    ticker:=time.NewTicker(interval); defer ticker.Stop()
    for {
        receipt,err:=n.TransactionReceipt(ctx,hash)
        if err==nil && receipt!=nil { return receipt,nil }
        select { case <-ctx.Done(): return nil,ctx.Err(); case <-ticker.C: }
    }
}
