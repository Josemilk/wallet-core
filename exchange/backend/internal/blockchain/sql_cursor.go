package blockchain

import("context";"errors";"github.com/Josemilk/wallet-core/exchange/backend/internal/db";"github.com/jackc/pgx/v5")

type SQLCursor struct{DB *db.DB}
func(s *SQLCursor)Load(ctx context.Context,network string)(uint64,string,error){if s==nil||s.DB==nil||s.DB.Pool==nil{return 0,"",errors.New("chain cursor database unavailable")};var h uint64;var hash string;err:=s.DB.Pool.QueryRow(ctx,`SELECT last_height,COALESCE(last_hash,'') FROM chain_cursors WHERE network=$1`,network).Scan(&h,&hash);if err==pgx.ErrNoRows{return 0,"",nil};if err!=nil{return 0,"",err};return h,hash,nil}
func(s *SQLCursor)Save(ctx context.Context,network string,height uint64,hash string)error{if s==nil||s.DB==nil||s.DB.Pool==nil{return errors.New("chain cursor database unavailable")};_,err:=s.DB.Pool.Exec(ctx,`INSERT INTO chain_cursors(network,last_height,last_hash,updated_at) VALUES($1,$2,$3,now()) ON CONFLICT(network) DO UPDATE SET last_height=EXCLUDED.last_height,last_hash=EXCLUDED.last_hash,updated_at=now()`,network,height,hash);return err}
