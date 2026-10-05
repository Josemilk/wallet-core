package deposits

import(
 "context"
 "errors"
 "strconv"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/ledger"
 "github.com/jackc/pgx/v5"
)

type Deposit struct{ID,UserID,AccountID,Asset,TxHash,Address string;Amount string;Confirmations uint64;RequiredConfirmations uint64}
type Service struct{DB *db.DB;Ledger *ledger.Poster}

func(s *Service)Confirm(ctx context.Context,d Deposit)error{
 if s==nil||s.DB==nil||s.DB.Pool==nil||s.Ledger==nil{return errors.New("deposit dependencies missing")};if d.ID==""||d.UserID==""||d.AccountID==""||d.Asset==""||d.TxHash==""||d.Address==""||d.Amount==""{return errors.New("invalid deposit")};if d.RequiredConfirmations>0&&d.Confirmations<d.RequiredConfirmations{return errors.New("insufficient confirmations")};amount,err:=strconv.ParseInt(d.Amount,10,64);if err!=nil||amount<=0{return errors.New("invalid deposit amount")}
 return db.WithTx(ctx,s.DB.Pool,func(tx pgx.Tx)error{
  var status string;if err:=tx.QueryRow(ctx,`SELECT status FROM deposits WHERE id=$1 FOR UPDATE`,d.ID).Scan(&status);err!=nil{return err};if status=="credited"{return nil}
  if _,err:=tx.Exec(ctx,`INSERT INTO accounts(user_id,asset,kind) VALUES(NULL,$1,'SYSTEM') ON CONFLICT DO NOTHING`,d.Asset);err!=nil{return err};var systemID string;if err:=tx.QueryRow(ctx,`SELECT id FROM accounts WHERE asset=$1 AND kind='SYSTEM'`,d.Asset).Scan(&systemID);err!=nil{return err}
  if _,err:=tx.Exec(ctx,`INSERT INTO balances(account_id,asset,available,locked,version) VALUES($1,$2,0,0,0) ON CONFLICT(account_id,asset) DO NOTHING`,d.AccountID,d.Asset);err!=nil{return err}
  if err:=s.Ledger.PostTx(ctx,tx,ledger.Posting{TransactionID:d.ID,IdempotencyKey:"deposit:"+d.ID,Type:"deposit",DebitAccountID:systemID,CreditAccountID:d.AccountID,Asset:d.Asset,Amount:amount});err!=nil{return err}
  _,err=tx.Exec(ctx,`UPDATE deposits SET status='credited',confirmations=$1 WHERE id=$2`,d.Confirmations,d.ID);return err
 })
}
