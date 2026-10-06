package settlement

import (
    "context"
    "errors"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/ledger"
    "github.com/jackc/pgx/v5"
)

type Trade struct { ID, BuyerAccount, SellerAccount, BaseAsset, QuoteAsset string; BaseAmount, QuoteAmount int64 }
type Service struct { DB *db.DB; Ledger *ledger.Postgres }

func (s *Service) Settle(ctx context.Context, t Trade) error {
    if s == nil || s.DB == nil || s.DB.Pool == nil || s.Ledger == nil { return errors.New("settlement dependencies missing") }
    if t.ID == "" || t.BuyerAccount == "" || t.SellerAccount == "" || t.BaseAsset == "" || t.QuoteAsset == "" || t.BaseAmount <= 0 || t.QuoteAmount <= 0 { return errors.New("invalid trade") }
    return db.WithTx(ctx, s.DB.Pool, func(tx pgx.Tx) error {
        var exists bool
        if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM settlement_records WHERE trade_id=$1)`, t.ID).Scan(&exists); err != nil { return err }
        if exists { return nil }
        if _, err := tx.Exec(ctx, `INSERT INTO settlement_records (trade_id, status) VALUES ($1,'pending')`, t.ID); err != nil { return err }
        if _, err := tx.Exec(ctx, `UPDATE settlement_records SET status='settled' WHERE trade_id=$1`, t.ID); err != nil { return err }
        return nil
    })
}
