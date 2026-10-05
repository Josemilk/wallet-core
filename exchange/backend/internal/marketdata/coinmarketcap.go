package marketdata

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "os"
 "strings"
 "time"
)

type Quote struct { ID int64 `json:"id"`; Name string `json:"name"`; Symbol string `json:"symbol"`; Price float64 `json:"price"`; Change24h float64 `json:"change24h"`; MarketCap float64 `json:"marketCap"`; Volume24h float64 `json:"volume24h"`; UpdatedAt time.Time `json:"updatedAt"` }
type Provider interface { Quotes(ctx context.Context, ids []string, convert string) ([]Quote,error) }
type CoinMarketCap struct { APIKey, BaseURL string; HTTP *http.Client }

func (c *CoinMarketCap) Quotes(ctx context.Context, ids []string, convert string) ([]Quote,error) {
 if c==nil{return nil,errors.New("market provider unavailable")}; if len(ids)==0{return nil,errors.New("no assets requested")}; if convert==""{convert="USD"}
 base:=c.BaseURL; if base==""{base="https://pro-api.coinmarketcap.com"}; path:="/v3/cryptocurrency/quotes/latest"; if c.APIKey=="" && os.Getenv("CMC_API_KEY")=="" { path="/public-api"+path }
 u,err:=url.Parse(strings.TrimRight(base,"/")+path);if err!=nil{return nil,err};q:=u.Query();q.Set("id",strings.Join(ids,","));q.Set("convert",convert);u.RawQuery=q.Encode()
 client:=c.HTTP;if client==nil{client=&http.Client{Timeout:10*time.Second}};req,err:=http.NewRequestWithContext(ctx,http.MethodGet,u.String(),nil);if err!=nil{return nil,err};key:=c.APIKey;if key==""{key=os.Getenv("CMC_API_KEY")};if key!=""{req.Header.Set("X-CMC_PRO_API_KEY",key)};req.Header.Set("Accept","application/json")
 resp,err:=client.Do(req);if err!=nil{return nil,err};defer resp.Body.Close();body,err:=io.ReadAll(io.LimitReader(resp.Body,4<<20));if err!=nil{return nil,err};if resp.StatusCode/100!=2{return nil,fmt.Errorf("market provider status=%d",resp.StatusCode)}
 var raw struct{Data map[string]struct{Name string `json:"name"`;Symbol string `json:"symbol"`;ID int64 `json:"id"`;Quote map[string]struct{Price float64 `json:"price"`;Change24h float64 `json:"percent_change_24h"`;MarketCap float64 `json:"market_cap"`;Volume24h float64 `json:"volume_24h"`;LastUpdated time.Time `json:"last_updated"`} `json:"quote"`} `json:"data"`}
 if err:=json.Unmarshal(body,&raw);err!=nil{return nil,err};out:=make([]Quote,0,len(raw.Data));for _,a:=range raw.Data{qq,ok:=raw.Data[fmt.Sprint(a.ID)].Quote[convert];if !ok{return nil,fmt.Errorf("missing %s quote for %d",convert,a.ID)};out=append(out,Quote{ID:a.ID,Name:a.Name,Symbol:a.Symbol,Price:qq.Price,Change24h:qq.Change24h,MarketCap:qq.MarketCap,Volume24h:qq.Volume24h,UpdatedAt:qq.LastUpdated})};return out,nil
}
