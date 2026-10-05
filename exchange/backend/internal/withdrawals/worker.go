package withdrawals

import(
 "context"
 "errors"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/custody"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/db"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/ledger"
 "github.com/jackc/pgx/v5"
)

type TransactionBuilder interface{Build(ctx context.Context,job Job)(string,error)}
type Broadcaster interface{Broadcast(ctx context.Context,network,signedTransaction string)(string,error)}
type Worker struct{DB *db.DB;Queue *Queue;Signer custody.Signer;Builder TransactionBuilder;Broadcaster Broadcaster;Ledger *ledger.Poster}

func(w *Worker) Process(ctx context.Context)(string,error){
 if w==nil||w.Queue==nil||w.Signer==nil||w.Builder==nil||w.Broadcaster==nil{return "",errors.New("withdrawal worker dependencies missing")}
 job,err:=w.Queue.Claim(ctx);if err!=nil{return "",err}
 if err:=w.Queue.SetStatus(ctx,job.ID,"signing","");err!=nil{return "",err}
 unsigned,err:=w.Builder.Build(ctx,job);if err!=nil{_ = w.fail(ctx,job);return "",err}
 signed,err:=w.Signer.Sign(ctx,custody.SignRequest{KeyRef:job.KeyRef,Network:job.Network,RawTransaction:unsigned});if err!=nil{_ = w.fail(ctx,job);return "",err}
 txHash,err:=w.Broadcaster.Broadcast(ctx,job.Network,signed);if err!=nil{_ = w.fail(ctx,job);return "",err}
 if err:=w.Queue.SetStatus(ctx,job.ID,"broadcast",txHash);err!=nil{return "",err};return txHash,nil
}

func(w *Worker) SettleConfirmed(ctx context.Context,job Job,systemAccount,txHash string)error{
 if w==nil||w.Ledger==nil||w.Queue==nil{return errors.New("withdrawal settlement dependencies missing")};if systemAccount==""||txHash==""{return errors.New("system account and tx hash required")}
 if err:=w.Ledger.Post(ctx,ledger.Posting{TransactionID:job.ID,IdempotencyKey:"withdrawal:"+job.ID,Type:"withdrawal",DebitAccountID:job.AccountID,CreditAccountID:systemAccount,Asset:job.Asset,Amount:job.Amount,DebitFromLocked:true});err!=nil{return err}
 return w.Queue.SetStatus(ctx,job.ID,"confirmed",txHash)
}

func(w *Worker) fail(ctx context.Context,job Job)error{
 if w==nil||w.DB==nil||w.DB.Pool==nil{return errors.New("database not initialized")}
 return db.WithTx(ctx,w.DB.Pool,func(tx pgx.Tx)error{
  res,err:=tx.Exec(ctx,`UPDATE balances SET available=available+$1,locked=locked-$1,version=version+1 WHERE account_id=$2 AND asset=$3 AND locked >= $1`,job.Amount,job.AccountID,job.Asset);if err!=nil{return err};if res.RowsAffected()!=1{return errors.New("locked withdrawal balance missing")}
  _,err=tx.Exec(ctx,`UPDATE withdrawal_queue SET status='failed',updated_at=now() WHERE id=$1 AND status IN ('risk_approved','signing')`,job.ID);return err
 })
}
