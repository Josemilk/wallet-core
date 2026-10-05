package ledger

import (
    "context"
    "errors"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
    "github.com/jackc/pgx/v5"
)

type Posting struct { TransactionID, IdempotencyKey, Type string; DebitAccountID, CreditAccountID, Asset string; Amount int64 }
type Poster struct { DB *db.DB }

// Positive amounts credit the destination account; negative amounts debit it.
// Every journal is exactly two entries with equal and opposite amounts.
func (p *Poster) Post(ctx context.Context, j Posting) error {
    if p==nil || p.DB==nil || p.DB.Pool==nil { return errors.New("database not initialized") }
    if j.TransactionID=="" || j.IdempotencyKey=="" || j.Type=="" || j.DebitAccountID=="" || j.CreditAccountID=="" || j.Asset=="" || j.Amount<=0 { return errors.New("invalid journal") }
    if j.DebitAccountID==j.CreditAccountID { return errors.New("journal accounts must differ") }
    return db.WithTx(ctx,p.DB.Pool,func(tx pgx.Tx) error {
        var exists bool
        if err:=tx.QueryRow(ctx,`SELECT EXISTS(SELECT 1 FROM ledger_transactions WHERE idempotency_key=$1)`,j.IdempotencyKey).Scan(&exists);err!=nil{return err}
        if exists{return nil}
        if _,err:=tx.Exec(ctx,`INSERT INTO ledger_transactions(id,idempotency_key,type) VALUES($1,$2,$3)`,j.TransactionID,j.IdempotencyKey,j.Type);err!=nil{return err}
        if _,err:=tx.Exec(ctx,`INSERT INTO ledger_entries(transaction_id,account_id,asset,amount) VALUES($1,$2,$3,$4),($1,$5,$3,$6)`,j.TransactionID,j.DebitAccountID,j.Asset,-j.Amount,j.CreditAccountID,j.Amount);err!=nil{return err}
        if r,err:=tx.Exec(ctx,`UPDATE balances SET available=available-$1,version=version+1 WHERE account_id=$2 AND asset=$3 AND available >= $1`,j.Amount,j.DebitAccountID,j.Asset);err!=nil{return err}else if r.RowsAffected()!=1{return errors.New("insufficient debit account balance")}
        if r,err:=tx.Exec(ctx,`UPDATE balances SET available=available+$1,version=version+1 WHERE account_id=$2 AND asset=$3`,j.Amount,j.CreditAccountID,j.Asset);err!=nil{return err}else if r.RowsAffected()!=1{return errors.New("credit account balance missing")}
        return nil
    })
}
