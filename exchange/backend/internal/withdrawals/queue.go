package withdrawals

import("context";"errors";"github.com/Josemilk/wallet-core/exchange/backend/internal/db";"github.com/jackc/pgx/v5")

type Queue struct{DB *db.DB}
func(q *Queue) Claim(ctx context.Context,limit int)([]string,error){if q==nil||q.DB==nil||q.DB.Pool==nil{return nil,errors.New("database not initialized")};if limit<=0||limit>100{return nil,errors.New("invalid queue limit")};var ids []string;err:=db.WithTx(ctx,q.DB.Pool,func(tx pgx.Tx)error{rows,err:=tx.Query(ctx,`SELECT id FROM withdrawals WHERE status='requested' ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT $1`,limit);if err!=nil{return err};defer rows.Close();for rows.Next(){var id string;if err:=rows.Scan(&id);err!=nil{return err};ids=append(ids,id)};if err:=rows.Err();err!=nil{return err};for _,id:=range ids{if _,err:=tx.Exec(ctx,`UPDATE withdrawals SET status='processing' WHERE id=$1 AND status='requested'`,id);err!=nil{return err}};return nil});return ids,err}
