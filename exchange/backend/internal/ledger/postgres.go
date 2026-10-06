package ledger

import (
    "context"
    "errors"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
    "github.com/jackc/pgx/v5"
)

type Postgres struct { db *db.DB }
func NewPostgres(database *db.DB) *Postgres { return &Postgres{db: database} }

func (p *Postgres) Post(ctx context.Context, j Journal) error {
    if p == nil || p.db == nil || p.db.Pool == nil { return errors.New("database not initialized") }
    if j.ID == "" || j.IdempotencyKey == "" || len(j.Entries) < 2 { return errors.New("invalid journal") }
    sums := make(map[string]int64)
    for _, e := range j.Entries {
        if e.Account == "" || e.Asset == "" || e.Amount == 0 { return errors.New("invalid entry") }
        sums[e.Asset] += e.Amount
    }
    for _, v := range sums { if v != 0 { return errors.New("unbalanced journal") } }

    return db.WithTx(ctx, p.db.Pool, func(tx pgx.Tx) error {
        var exists bool
        if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM ledger_transactions WHERE idempotency_key=$1)`, j.IdempotencyKey).Scan(&exists); err != nil { return err }
        if exists { return nil }
        if _, err := tx.Exec(ctx, `INSERT INTO ledger_transactions (id, idempotency_key, type) VALUES ($1,$2,$3)`, j.ID, j.IdempotencyKey, "EXCHANGE_JOURNAL"); err != nil { return err }
        for _, e := range j.Entries {
            if _, err := tx.Exec(ctx, `INSERT INTO ledger_entries (transaction_id, account_id, asset, amount) VALUES ($1,$2,$3,$4)`, j.ID, e.Account, e.Asset, e.Amount); err != nil { return err }
        }
        return nil
    })
}
