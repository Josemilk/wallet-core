package settlement

import (
 "context"
 "errors"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
 "github.com/jackc/pgx/v5"
)

type Trade struct { ID, BuyAccount, SellAccount, BaseAsset, QuoteAsset string; Price, Quantity int64 }
type Service struct { DB *db.DB }

func (s *Service) Settle(ctx context.Context,t Trade) error {
 if s==nil||s.DB==nil||s.DB.Pool==nil{return errors.New("database not initialized")}
 if t.ID==""||t.BuyAccount==""||t.SellAccount==""||t.BaseAsset==""||t.QuoteAsset==""||t.Price<=0||t.Quantity<=0{return errors.New("invalid trade")}
 nocional:=t.Price*t.Quantity
 return db.WithTx(ctx,s.DB.Pool,func(tx pgx.Tx) error{
  var done bool
  if err:=tx.QueryRow(ctx,`SELECT EXISTS(SELECT 1 FROM settlement_records WHERE trade_id=$1)`,t.ID).Scan(&done);err!=nil{return err};if done{return nil}
  if _,err:=tx.Exec(ctx,`UPDATE balances SET locked=locked-$1,version=version+1 WHERE account_id=$2 AND asset=$3 AND locked >= $1`,nocional,t.BuyAccount,t.QuoteAsset);err!=nil{return err}
  if _,err:=tx.Exec(ctx,`UPDATE balances SET available=available+$1,version=version+1 WHERE account_id=$2 AND asset=$3`,t.BuyAccount,t.BaseAsset,t.Quantity);err!=nil{return err}
  if _,err:=tx.Exec(ctx,`UPDATE balances SET locked=locked-$1,version=version+1 WHERE account_id=$2 AND asset=$3 AND locked >= $1`,t.Quantity,t.SellAccount,t.BaseAsset);err!=nil{return err}
  if _,err:=tx.Exec(ctx,`UPDATE balances SET available=available+$1,version=version+1 WHERE account_id=$2 AND asset=$3`,nocional,t.SellAccount,t.QuoteAsset);err!=nil{return err}
  _,err:=tx.Exec(ctx,`INSERT INTO settlement_records(trade_id,status) VALUES($1,'settled')`,t.ID);return err
 })
}
