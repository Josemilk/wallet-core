package api

import("context";"encoding/json";"net/http";"strings";"github.com/Josemilk/wallet-core/exchange/backend/internal/auth";"github.com/Josemilk/wallet-core/exchange/backend/internal/db")

type HistoryItem struct{Kind string `json:"kind"`;ID string `json:"id"`;Asset string `json:"asset"`;Amount string `json:"amount"`;Status string `json:"status"`;CreatedAt string `json:"createdAt"`;TxHash string `json:"txHash,omitempty"`}
type HistorySource interface{History(context.Context,string,int)([]HistoryItem,error)}
type PostgresHistorySource struct{DB *db.DB}
func(s *PostgresHistorySource)History(ctx context.Context,userID string,limit int)([]HistoryItem,error){if limit<=0||limit>200{limit=50};rows,err:=s.DB.Pool.Query(ctx,`SELECT kind,id,asset,amount,status,created_at,tx_hash FROM (SELECT 'deposit' kind,id::text,asset,amount::text,status,created_at,tx_hash FROM deposits WHERE user_id=$1 UNION ALL SELECT 'withdrawal',id::text,asset,amount::text,status,created_at,COALESCE(tx_hash,'') FROM withdrawals WHERE user_id=$1) x ORDER BY created_at DESC LIMIT $2`,userID,limit);if err!=nil{return nil,err};defer rows.Close();out:=[]HistoryItem{};for rows.Next(){var x HistoryItem;if err:=rows.Scan(&x.Kind,&x.ID,&x.Asset,&x.Amount,&x.Status,&x.CreatedAt,&x.TxHash);err!=nil{return nil,err};out=append(out,x)};return out,rows.Err()}

type HistoryHandler struct{Auth auth.Verifier;Source HistorySource}
func(h *HistoryHandler)ServeHTTP(w http.ResponseWriter,r *http.Request){if h==nil||h.Auth==nil||h.Source==nil{http.Error(w,"service unavailable",503);return};if r.Method!=http.MethodGet||r.URL.Path!="/v1/history"{http.NotFound(w,r);return};bearer:=strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "));p,err:=auth.Require(r.Context(),h.Auth,bearer);if err!=nil{http.Error(w,"unauthorized",401);return};items,err:=h.Source.History(r.Context(),p.UserID,50);if err!=nil{http.Error(w,"history unavailable",500);return};w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store");_ = json.NewEncoder(w).Encode(items)}
