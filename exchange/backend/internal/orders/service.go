package orders

import (
    "context"
    "errors"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
    "github.com/jackc/pgx/v5"
)

type CreateRequest struct { ID, AccountID, Market, Side, Type, ClientOrderID string; Price, Quantity int64 }
type Service struct { DB *db.DB }

func (s *Service) Create(ctx context.Context, r CreateRequest) error {
    if s == nil || s.DB == nil || s.DB.Pool == nil { return errors.New("database not initialized") }
    if r.ID == "" || r.AccountID == "" || r.Market == "" || (r.Side != "buy" && r.Side != "sell") || (r.Type != "limit" && r.Type != "market") || r.Quantity <= 0 { return errors.New("invalid order") }
    if r.Type == "limit" && r.Price <= 0 { return errors.New("invalid price") }
    return db.WithTx(ctx, s.DB.Pool, func(tx pgx.Tx) error {
        if r.ClientOrderID != "" {
            var exists bool
            if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM orders WHERE account_id=$1 AND client_order_id=$2)`, r.AccountID, r.ClientOrderID).Scan(&exists); err != nil { return err }
            if exists { return nil }
        }
        if r.Side == "buy" {
            if _, err := tx.Exec(ctx, `UPDATE balances SET available=available-$1, locked=locked+$1, version=version+1 WHERE account_id=$2 AND asset=$3 AND available >= $1`, r.Price*r.Quantity, r.AccountID, quoteAsset(r.Market)); err != nil { return err }
        } else {
            if _, err := tx.Exec(ctx, `UPDATE balances SET available=available-$1, locked=locked+$1, version=version+1 WHERE account_id=$2 AND asset=$3 AND available >= $1`, r.Quantity, r.AccountID, baseAsset(r.Market)); err != nil { return err }
        }
        _, err := tx.Exec(ctx, `INSERT INTO orders(id,account_id,market,side,type,price,quantity,status,client_order_id) VALUES($1,$2,$3,$4,$5,$6,$7,'open',NULLIF($8,''))`, r.ID,r.AccountID,r.Market,r.Side,r.Type,r.Price,r.Quantity,r.ClientOrderID)
        return err
    })
}

func quoteAsset(m string) string { for i:=len(m)-1;i>=0;i-- { if m[i]=='/' { return m[i+1:] } }; return m }
func baseAsset(m string) string { for i:=0;i<len(m);i++ { if m[i]=='/' { return m[:i] } }; return m }
