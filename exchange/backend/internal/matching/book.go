package matching

import (
    "encoding/json"
    "net/http"
    "sort"
    "sync"
    "time"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/ledger"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/risk"
)

type Order struct { ID, UserID, Symbol string; Side string; Price, Quantity int64; Remaining int64; Time time.Time }
type Book struct { mu sync.Mutex; orders map[string][]Order }
func NewBook() *Book { return &Book{orders: make(map[string][]Order)} }
func (b *Book) Snapshot(symbol string) map[string]any { b.mu.Lock(); defer b.mu.Unlock(); o:=append([]Order(nil), b.orders[symbol]...); sort.Slice(o,func(i,j int)bool{if o[i].Side==o[j].Side { if o[i].Side=="BUY" { return o[i].Price>o[j].Price }; return o[i].Price<o[j].Price }; return o[i].Side=="BUY"}); return map[string]any{"symbol":symbol,"orders":o} }
func (b *Book) Place(o Order) []Order { b.mu.Lock(); defer b.mu.Unlock(); if o.Remaining==0{o.Remaining=o.Quantity}; fills:=[]Order{}; for len(fills)==0 && o.Remaining>0 { idx:=-1; for i,x:=range b.orders[o.Symbol] { if x.Side==o.Side {continue}; if o.Side=="BUY" && o.Price<x.Price {continue}; if o.Side=="SELL" && o.Price>x.Price {continue}; idx=i; break }; if idx<0 {break}; x:=b.orders[o.Symbol][idx]; q:=o.Remaining; if x.Remaining<q{q=x.Remaining}; o.Remaining-=q; x.Remaining-=q; fills=append(fills,Order{ID:x.ID,UserID:x.UserID,Symbol:x.Symbol,Side:x.Side,Price:x.Price,Quantity:q,Remaining:x.Remaining,Time:x.Time}); if x.Remaining==0 {b.orders[o.Symbol]=append(b.orders[o.Symbol][:idx],b.orders[o.Symbol][idx+1:]...)} else {b.orders[o.Symbol][idx]=x} }; if o.Remaining>0 {b.orders[o.Symbol]=append(b.orders[o.Symbol],o)}; return fills }
func HandleOrderHTTP(w http.ResponseWriter,r *http.Request,b *Book,l *ledger.Ledger,re *risk.Engine){var o Order;if json.NewDecoder(r.Body).Decode(&o)!=nil {http.Error(w,"bad request",400);return}; if err:=re.Check(o.UserID,o.Symbol,o.Side,o.Price,o.Quantity);err!=nil{http.Error(w,err.Error(),422);return}; o.Time=time.Now().UTC();fills:=b.Place(o);_ = l; w.Header().Set("Content-Type","application/json");_ = json.NewEncoder(w).Encode(map[string]any{"order":o,"fills":fills})}
