package risk

import (
    "encoding/json"
    "errors"
    "net/http"
)

type Engine struct { MaxOrderNotional int64; MaxQuantity int64 }
func New() *Engine { return &Engine{MaxOrderNotional: 1_000_000_000_000, MaxQuantity: 1_000_000_000_000} }
func (e *Engine) Check(user, symbol, side string, price, quantity int64) error { if user==""||symbol=="" {return errors.New("identity required")}; if side!="BUY"&&side!="SELL" {return errors.New("invalid side")}; if price<=0||quantity<=0 {return errors.New("invalid amount")}; if quantity>e.MaxQuantity {return errors.New("quantity limit")}; if price > e.MaxOrderNotional/quantity {return errors.New("notional limit")}; return nil }
func HandleHTTP(w http.ResponseWriter,r *http.Request,e *Engine){var x struct{UserID,Symbol,Side string;Price,Quantity int64};if json.NewDecoder(r.Body).Decode(&x)!=nil{http.Error(w,"bad request",400);return};if err:=e.Check(x.UserID,x.Symbol,x.Side,x.Price,x.Quantity);err!=nil{http.Error(w,err.Error(),422);return};w.Header().Set("Content-Type","application/json");_ = json.NewEncoder(w).Encode(map[string]any{"allowed":true})}
