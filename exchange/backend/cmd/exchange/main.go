package main

import("context";"encoding/json";"log";"net/http";"os";"strings";"time";"github.com/Josemilk/wallet-core/exchange/backend/internal/api";"github.com/Josemilk/wallet-core/exchange/backend/internal/auth";"github.com/Josemilk/wallet-core/exchange/backend/internal/db";"github.com/Josemilk/wallet-core/exchange/backend/internal/ledger";"github.com/Josemilk/wallet-core/exchange/backend/internal/marketdata";"github.com/Josemilk/wallet-core/exchange/backend/internal/matching";"github.com/Josemilk/wallet-core/exchange/backend/internal/risk";"github.com/Josemilk/wallet-core/exchange/backend/internal/ws";"github.com/gorilla/websocket")

type server struct{ledger *ledger.Ledger;book *matching.Book;risk *risk.Engine;db *db.DB;dashboard *api.Handler;history *api.HistoryHandler;socket *ws.HTTPHandler;auth *auth.OIDCVerifier}
func main(){addr:=os.Getenv("EXCHANGE_HTTP_ADDR");if addr==""{addr=":8080"};s:=&server{ledger:ledger.New(),book:matching.NewBook(),risk:risk.New()};ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel();if os.Getenv("DATABASE_URL")!=""{d,err:=db.Open(ctx);if err!=nil{log.Fatal(err)};s.db=d;defer d.Close()};if s.db!=nil&&os.Getenv("OIDC_ISSUER")!=""&&os.Getenv("OIDC_AUDIENCE")!=""&&os.Getenv("OIDC_JWKS_URL")!=""{provider:=&marketdata.CoinMarketCap{APIKey:os.Getenv("CMC_API_KEY")};verifier:=&auth.OIDCVerifier{Issuer:os.Getenv("OIDC_ISSUER"),Audience:os.Getenv("OIDC_AUDIENCE"),JWKSURL:os.Getenv("OIDC_JWKS_URL")};s.auth=verifier;source:=&api.PostgresDashboardSource{DB:s.db,Market:provider,AssetIDs:splitCSV(getenv("CMC_ASSET_IDS","1,1027,5426"))};s.dashboard=&api.Handler{Service:&api.Service{Auth:verifier,Market:provider,Source:source},AssetIDs:source.AssetIDs};s.history=&api.HistoryHandler{Auth:verifier,Source:&api.PostgresHistorySource{DB:s.db}};s.socket=&ws.HTTPHandler{Auth:verifier,Hub:ws.NewHub(),Upgrader:websocket.Upgrader{ReadBufferSize:4096,WriteBufferSize:4096,EnableCompression:true}}}
mux:=http.NewServeMux();mux.HandleFunc("GET /healthz",s.health);mux.HandleFunc("GET /readyz",s.ready);mux.HandleFunc("GET /v1/orderbook",s.orderbook);mux.Handle("POST /v1/orders",s.requireUser(http.HandlerFunc(s.placeOrder)));mux.Handle("POST /v1/risk/check",s.requireUser(http.HandlerFunc(s.riskCheck)));if s.dashboard!=nil{mux.Handle("/v1/dashboard",s.dashboard);mux.Handle("/v1/history",s.history);mux.Handle("/v1/ws",s.socket)};srv:=&http.Server{Addr:addr,Handler:securityHeaders(mux),ReadHeaderTimeout:5*time.Second,ReadTimeout:30*time.Second,WriteTimeout:30*time.Second,IdleTimeout:60*time.Second};log.Printf("exchange service listening on %s",addr);log.Fatal(srv.ListenAndServe())}
func getenv(k,d string)string{if v:=os.Getenv(k);v!=""{return v};return d}
func splitCSV(v string)[]string{parts:=strings.Split(v,",");out:=make([]string,0,len(parts));for _,p:=range parts{p=strings.TrimSpace(p);if p!=""{out=append(out,p)}};return out}
func securityHeaders(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("Cache-Control","no-store");w.Header().Set("X-Content-Type-Options","nosniff");w.Header().Set("X-Frame-Options","DENY");w.Header().Set("Referrer-Policy","no-referrer");next.ServeHTTP(w,r)})}
func(s *server)health(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,map[string]any{"status":"ok"})}
func(s *server)ready(w http.ResponseWriter,r *http.Request){if s.db!=nil{if err:=s.db.Health(r.Context());err!=nil{writeJSON(w,503,map[string]any{"status":"not_ready","database":"down"});return}};if err:=s.ledger.Ready(context.Background());err!=nil{writeJSON(w,503,map[string]any{"status":"not_ready"});return};writeJSON(w,200,map[string]any{"status":"ready"})}
func(s *server)orderbook(w http.ResponseWriter,r *http.Request){writeJSON(w,200,s.book.Snapshot(r.URL.Query().Get("symbol")))}
func(s *server)placeOrder(w http.ResponseWriter,r *http.Request){if p,ok:=principalFromContext(r.Context());ok{r.Header.Set("X-Authenticated-User",p.UserID)};matching.HandleOrderHTTP(w,r,s.book,s.ledger,s.risk)}
func(s *server)riskCheck(w http.ResponseWriter,r *http.Request){risk.HandleHTTP(w,r,s.risk)}
func writeJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}

func(s *server)requireUser(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  if s.auth==nil { writeJSON(w,http.StatusServiceUnavailable,map[string]any{"error":"authentication is not configured"}); return }
  token:=strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "))
  p,err:=auth.Require(r.Context(),s.auth,token)
  if err!=nil { writeJSON(w,http.StatusUnauthorized,map[string]any{"error":"unauthorized"}); return }
  ctx:=context.WithValue(r.Context(),principalKey{},p)
  next.ServeHTTP(w,r.WithContext(ctx))
 })
}
type principalKey struct{}
func principalFromContext(ctx context.Context)(auth.Principal,bool){p,ok:=ctx.Value(principalKey{}).(auth.Principal);return p,ok}
