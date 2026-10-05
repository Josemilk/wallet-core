package withdrawals

import (
    "context"
    "errors"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
    "github.com/jackc/pgx/v5"
)

type StateStore struct { DB *db.DB }

func (s *StateStore) Transition(ctx context.Context, id, from, to string) error {
    if s == nil || s.DB == nil || s.DB.Pool == nil { return errors.New("database not initialized") }
    if id == "" || from == "" || to == "" { return errors.New("invalid withdrawal transition") }
    return db.WithTx(ctx, s.DB.Pool, func(tx pgx.Tx) error {
        var status string
        if err := tx.QueryRow(ctx, `SELECT status FROM withdrawal_state WHERE withdrawal_id=$1 FOR UPDATE`, id).Scan(&status); err != nil { return err }
        if status != from { return errors.New("withdrawal state conflict") }
        tag, err := tx.Exec(ctx, `UPDATE withdrawal_state SET status=$1, updated_at=now() WHERE withdrawal_id=$2 AND status=$3`, to, id, from)
        if err != nil { return err }
        if tag.RowsAffected() != 1 { return errors.New("withdrawal transition rejected") }
        return nil
    })
}
