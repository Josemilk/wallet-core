package orders

import (
    "context"
    "errors"
    "math"
    "strings"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
    "github.com/jackc/pgx/v5"
)

type CreateRequest struct { ID, AccountID, Market, Side, Type, ClientOrderID string; Price, Quantity int64 }
type Service struct { DB *db.DB }

func (s *Service) Create(ctx context.Context, r CreateRequest) error {
    if s == nil || s.DB == nil || s.DB.Pool == nil { return errors.New("database not initialized") }
    if r.ID == "" || r.AccountID == "" || r.Market == "" || (r.Side != "buy" && r.Side != "sell") || (r.Type != "limit" && r.Type != "market") || r.Quantity <= 0 { return errors.New("invalid order") }
    if r.Type == "market" { return errors.New("market orders require an explicit quote reserve") }
    if r.Price <= 0 { return errors.New("invalid price") }
    if r.Price > math.MaxInt64/r.Quantity { return errors.New("order notional overflow") }
    notional := r.Price * r.Quantity
    side := strings.ToUpper(r.Side)
    orderType := strings.ToUpper(r.Type)

    return db.WithTx(ctx, s.DB.Pool, func(tx pgx.Tx) error {
        if r.ClientOrderID != "" {
            var exists bool
            if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM orders WHERE user_id=(SELECT user_id FROM accounts WHERE id=$1) AND client_order_id=$2)`, r.AccountID, r.ClientOrderID).Scan(&exists); err != nil { return err }
            if exists { return nil }
        }

        reserve := r.Quantity
        asset := baseAsset(r.Market)
        if side == "BUY" { reserve = notional; asset = quoteAsset(r.Market) }
        res, err := tx.Exec(ctx, `UPDATE balances SET available=available-$1, locked=locked+$1, version=version+1 WHERE account_id=$2 AND asset=$3 AND available >= $1`, reserve, r.AccountID, asset)
        if err != nil { return err }
        if res.RowsAffected() != 1 { return errors.New("insufficient available balance or balance account missing") }

        _, err = tx.Exec(ctx, `INSERT INTO orders(id,user_id,symbol,side,order_type,price,quantity,remaining,status,client_order_id,account_id,market,type) SELECT $1,user_id,$3,$4,$5,$6,$7,$7,'OPEN',$8,$2,$3,$5 FROM accounts WHERE id=$2`, r.ID, r.AccountID, r.Market, side, orderType, r.Price, r.Quantity, r.ClientOrderID)
        return err
    })
}

func quoteAsset(m string) string { for i:=len(m)-1;i>=0;i-- { if m[i]=='/' { return m[i+1:] } }; return m }
func baseAsset(m string) string { for i:=0;i<len(m);i++ { if m[i]=='/' { return m[:i] } }; return m }
