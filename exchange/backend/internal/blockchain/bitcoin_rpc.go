package blockchain

import("bytes";"context";"encoding/base64";"encoding/hex";"encoding/json";"errors";"fmt";"io";"net/http";"strings";"time")

type BitcoinRPC struct{URL,Username,Password string;HTTP *http.Client}
type btcReq struct{JSONRPC string `json:"jsonrpc"`;ID uint64 `json:"id"`;Method string `json:"method"`;Params any `json:"params"`}
type btcResp struct{Result json.RawMessage `json:"result"`;Error *struct{Code int `json:"code"`;Message string `json:"message"`} `json:"error"`}
func(n *BitcoinRPC)call(ctx context.Context,method string,params any,out any)error{if n==nil||n.URL==""{return errors.New("Bitcoin RPC URL is required")};body,_:=json.Marshal(btcReq{"1.0",1,method,params});req,err:=http.NewRequestWithContext(ctx,http.MethodPost,n.URL,bytes.NewReader(body));if err!=nil{return err};req.Header.Set("Content-Type","application/json");if n.Username!=""{req.Header.Set("Authorization","Basic "+base64.StdEncoding.EncodeToString([]byte(n.Username+":"+n.Password)))};c:=n.HTTP;if c==nil{c=&http.Client{Timeout:15*time.Second}};resp,err:=c.Do(req);if err!=nil{return err};defer resp.Body.Close();raw,err:=io.ReadAll(io.LimitReader(resp.Body,8<<20));if err!=nil{return err};if resp.StatusCode/100!=2{return fmt.Errorf("Bitcoin RPC HTTP %d",resp.StatusCode)};var r btcResp;if err:=json.Unmarshal(raw,&r);err!=nil{return err};if r.Error!=nil{return fmt.Errorf("Bitcoin RPC %d: %s",r.Error.Code,r.Error.Message)};if out==nil{return nil};return json.Unmarshal(r.Result,out)}
func(n *BitcoinRPC)BlockHash(ctx context.Context,height uint64)(string,error){var h string;if err:=n.call(ctx,"getblockhash",[]any{height},&h);err!=nil{return "",err};return h,nil}
func(n *BitcoinRPC)Block(ctx context.Context,hash string)(map[string]any,error){var b map[string]any;if err:=n.call(ctx,"getblock",[]any{hash,2},&b);err!=nil{return nil,err};return b,nil}
func(n *BitcoinRPC)RawTransaction(ctx context.Context,hash string)(map[string]any,error){var t map[string]any;if err:=n.call(ctx,"getrawtransaction",[]any{hash,true},&t);err!=nil{return nil,err};return t,nil}
func(n *BitcoinRPC)SendRawTransaction(ctx context.Context,raw string)(string,error){raw=strings.TrimSpace(raw);if raw==""{return "",errors.New("empty raw transaction")};if len(raw)%2!=0{return "",errors.New("raw Bitcoin transaction must have even-length hex")};if _,err:=hex.DecodeString(raw);err!=nil{return "",errors.New("raw Bitcoin transaction is not valid hex")};var h string;if err:=n.call(ctx,"sendrawtransaction",[]any{raw},&h);err!=nil{return "",err};return h,nil}

func(n *BitcoinRPC) WaitForTransaction(ctx context.Context, hash string, interval time.Duration) (map[string]any, error) {
    if strings.TrimSpace(hash)=="" { return nil, errors.New("transaction hash is required") }
    if interval<=0 { interval=3*time.Second }
    ticker:=time.NewTicker(interval); defer ticker.Stop()
    for {
        tx,err:=n.RawTransaction(ctx,hash)
        if err==nil && tx!=nil { return tx,nil }
        select { case <-ctx.Done(): return nil,ctx.Err(); case <-ticker.C: }
    }
}
