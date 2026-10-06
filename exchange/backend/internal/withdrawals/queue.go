package withdrawals

import(
 "context"
 "errors"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
 "github.com/jackc/pgx/v5"
)

type Queue struct{DB *db.DB}
type EnqueueRequest struct{ID,UserID,AccountID,Asset,Destination,Network,KeyRef string; Amount int64}
type Job struct{ID,UserID,AccountID,Asset,Destination,Network,KeyRef string; Amount int64; AttemptCount int}

func(q *Queue) Enqueue(ctx context.Context,r EnqueueRequest)error{
 if q==nil||q.DB==nil||q.DB.Pool==nil{return errors.New("database not initialized")}
 if r.ID==""||r.UserID==""||r.AccountID==""||r.Asset==""||r.Destination==""||r.Network==""||r.KeyRef==""||r.Amount<=0{return errors.New("invalid withdrawal")}
 return db.WithTx(ctx,q.DB.Pool,func(tx pgx.Tx)error{
  res,err:=tx.Exec(ctx,`UPDATE balances SET available=available-$1,locked=locked+$1,version=version+1 WHERE account_id=$2 AND asset=$3 AND available >= $1`,r.Amount,r.AccountID,r.Asset);if err!=nil{return err};if res.RowsAffected()!=1{return errors.New("insufficient available balance")}
  _,err=tx.Exec(ctx,`INSERT INTO withdrawal_queue(id,user_id,account_id,asset,amount,destination,network,key_ref,status) VALUES($1,$2,$3,$4,$5,$6,$7,$8,'queued')`,r.ID,r.UserID,r.AccountID,r.Asset,r.Amount,r.Destination,r.Network,r.KeyRef);return err
 })
}

func(q *Queue) Claim(ctx context.Context)(Job,error){
 if q==nil||q.DB==nil||q.DB.Pool==nil{return Job{},errors.New("database not initialized")};var j Job
 err:=db.WithTx(ctx,q.DB.Pool,func(tx pgx.Tx)error{
  row:=tx.QueryRow(ctx,`SELECT id,user_id,account_id,asset,amount,destination,network,key_ref,attempt_count FROM withdrawal_queue WHERE status='queued' AND (available_after IS NULL OR available_after<=now()) ORDER BY created_at,id FOR UPDATE SKIP LOCKED LIMIT 1`)
  if err:=row.Scan(&j.ID,&j.UserID,&j.AccountID,&j.Asset,&j.Amount,&j.Destination,&j.Network,&j.KeyRef,&j.AttemptCount);err!=nil{return err}
  _,err:=tx.Exec(ctx,`UPDATE withdrawal_queue SET status='risk_approved',attempt_count=attempt_count+1,updated_at=now() WHERE id=$1`,j.ID);return err
 });return j,err
}

// Legacy batch claim retained for callers using the original withdrawals table.
func(q *Queue) ClaimLegacy(ctx context.Context,limit int)([]string,error){
 if q==nil||q.DB==nil||q.DB.Pool==nil{return nil,errors.New("database not initialized")};if limit<=0||limit>100{return nil,errors.New("invalid queue limit")};var ids []string
 err:=db.WithTx(ctx,q.DB.Pool,func(tx pgx.Tx)error{rows,err:=tx.Query(ctx,`SELECT id FROM withdrawals WHERE status='requested' ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT $1`,limit);if err!=nil{return err};defer rows.Close();for rows.Next(){var id string;if err:=rows.Scan(&id);err!=nil{return err};ids=append(ids,id)};if err:=rows.Err();err!=nil{return err};for _,id:=range ids{if _,err:=tx.Exec(ctx,`UPDATE withdrawals SET status='processing' WHERE id=$1 AND status='requested'`,id);err!=nil{return err}};return nil});return ids,err
}

func(q *Queue) SetStatus(ctx context.Context,id,status,txHash string)error{
 if q==nil||q.DB==nil||q.DB.Pool==nil{return errors.New("database not initialized")};allowed:=map[string]bool{"risk_approved":true,"signing":true,"broadcast":true,"confirmed":true,"failed":true,"cancelled":true};if !allowed[status]{return errors.New("invalid withdrawal status")};_,err:=q.DB.Pool.Exec(ctx,`UPDATE withdrawal_queue SET status=$1,tx_hash=NULLIF($2,''),updated_at=now() WHERE id=$3`,status,txHash,id);return err
}
