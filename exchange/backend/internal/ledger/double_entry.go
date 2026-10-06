package ledger

import(
 "context"
 "errors"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
 "github.com/jackc/pgx/v5"
)

type Posting struct{TransactionID,IdempotencyKey,Type string;DebitAccountID,CreditAccountID,Asset string;Amount int64;DebitFromLocked bool}
type Poster struct{DB *db.DB}

func(p *Poster)Post(ctx context.Context,j Posting)error{if p==nil||p.DB==nil||p.DB.Pool==nil{return errors.New("database not initialized")};return db.WithTx(ctx,p.DB.Pool,func(tx pgx.Tx)error{return p.PostTx(ctx,tx,j)})}

func(p *Poster)PostTx(ctx context.Context,tx pgx.Tx,j Posting)error{
 if p==nil||tx==nil{return errors.New("ledger transaction unavailable")};if j.TransactionID==""||j.IdempotencyKey==""||j.Type==""||j.DebitAccountID==""||j.CreditAccountID==""||j.Asset==""||j.Amount<=0{return errors.New("invalid journal")};if j.DebitAccountID==j.CreditAccountID{return errors.New("journal accounts must differ")}
 var exists bool;if err:=tx.QueryRow(ctx,`SELECT EXISTS(SELECT 1 FROM ledger_transactions WHERE idempotency_key=$1)`,j.IdempotencyKey).Scan(&exists);err!=nil{return err};if exists{return nil}
 var debitKind,creditKind string;if err:=tx.QueryRow(ctx,`SELECT kind FROM accounts WHERE id=$1 FOR UPDATE`,j.DebitAccountID).Scan(&debitKind);err!=nil{return err};if err:=tx.QueryRow(ctx,`SELECT kind FROM accounts WHERE id=$1 FOR UPDATE`,j.CreditAccountID).Scan(&creditKind);err!=nil{return err}
 if _,err:=tx.Exec(ctx,`INSERT INTO ledger_transactions(id,idempotency_key,type) VALUES($1,$2,$3)`,j.TransactionID,j.IdempotencyKey,j.Type);err!=nil{return err};if _,err:=tx.Exec(ctx,`INSERT INTO ledger_entries(transaction_id,account_id,asset,amount) VALUES($1,$2,$3,$4),($1,$5,$3,$6)`,j.TransactionID,j.DebitAccountID,j.Asset,-j.Amount,j.CreditAccountID,j.Amount);err!=nil{return err}
 if debitKind!="SYSTEM"{column:="available";if j.DebitFromLocked{column="locked"};q:=`UPDATE balances SET `+column+`=`+column+`-$1,version=version+1 WHERE account_id=$2 AND asset=$3 AND `+column+` >= $1`;if r,err:=tx.Exec(ctx,q,j.Amount,j.DebitAccountID,j.Asset);err!=nil{return err}else if r.RowsAffected()!=1{return errors.New("insufficient debit account balance")}}
 if creditKind!="SYSTEM"{if r,err:=tx.Exec(ctx,`UPDATE balances SET available=available+$1,version=version+1 WHERE account_id=$2 AND asset=$3`,j.Amount,j.CreditAccountID,j.Asset);err!=nil{return err}else if r.RowsAffected()!=1{return errors.New("credit account balance missing")}}
 return nil
}
