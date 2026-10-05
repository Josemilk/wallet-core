package deposits

import (
 "context"
 "errors"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
 "github.com/jackc/pgx/v5"
)

type LedgerCredit struct { TransactionID, AccountID, Asset string; Amount int64; IdempotencyKey string }
type LedgerCreditor interface { Credit(ctx context.Context, c LedgerCredit) error }
type PersistentCreditor struct { DB *db.DB }

func (p *PersistentCreditor) Credit(ctx context.Context,c LedgerCredit) error {
 if p==nil||p.DB==nil||p.DB.Pool==nil{return errors.New("database not initialized")}
 if c.TransactionID==""||c.AccountID==""||c.Asset==""||c.Amount<=0||c.IdempotencyKey==""{return errors.New("invalid ledger credit")}
 return db.WithTx(ctx,p.DB.Pool,func(tx pgx.Tx)error{
  var exists bool
  if err:=tx.QueryRow(ctx,`SELECT EXISTS(SELECT 1 FROM ledger_transactions WHERE idempotency_key=$1)`,c.IdempotencyKey).Scan(&exists);err!=nil{return err}
  if exists{return nil}
  if _,err:=tx.Exec(ctx,`INSERT INTO ledger_transactions(id,idempotency_key,type) VALUES($1,$2,'deposit')`,c.TransactionID,c.IdempotencyKey);err!=nil{return err}
  if _,err:=tx.Exec(ctx,`INSERT INTO ledger_entries(transaction_id,account_id,asset,amount) VALUES($1,$2,$3,$4)`,c.TransactionID,c.AccountID,c.Asset,c.Amount);err!=nil{return err}
  if _,err:=tx.Exec(ctx,`UPDATE balances SET available=available+$1,version=version+1 WHERE account_id=$2 AND asset=$3`,c.Amount,c.AccountID,c.Asset);err!=nil{return err}
  return nil
 })
}
