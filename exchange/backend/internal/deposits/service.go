package deposits

import (
    "context"
    "errors"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/ledger"
    "github.com/jackc/pgx/v5"
)

type Deposit struct { ID, UserID, AccountID, Asset, TxHash, Address string; Amount string; Confirmations uint64 }
type Service struct { DB *db.DB; Ledger *ledger.Postgres }

func (s *Service) Confirm(ctx context.Context, d Deposit) error {
    if s == nil || s.DB == nil || s.DB.Pool == nil || s.Ledger == nil { return errors.New("deposit dependencies missing") }
    if d.ID == "" || d.UserID == "" || d.AccountID == "" || d.Asset == "" || d.TxHash == "" || d.Address == "" || d.Amount == "" { return errors.New("invalid deposit") }
    if d.Confirmations == 0 { return errors.New("deposit is not confirmed") }
    return db.WithTx(ctx, s.DB.Pool, func(tx pgx.Tx) error {
        var status string
        err := tx.QueryRow(ctx, `SELECT status FROM deposits WHERE id=$1 FOR UPDATE`, d.ID).Scan(&status)
        if err != nil { return err }
        if status == "credited" { return nil }
        if _, err := tx.Exec(ctx, `UPDATE deposits SET status='credited', confirmations=$1 WHERE id=$2`, d.Confirmations, d.ID); err != nil { return err }
        return nil
    })
}
