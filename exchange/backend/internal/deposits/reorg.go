package deposits

import("context";"github.com/Josemilk/wallet-core/exchange/backend/internal/blockchain")

func(o *SQLObserver)RollbackBlock(ctx context.Context,network,blockHash string)error{if o==nil||o.DB==nil||o.DB.Pool==nil{return nil};rows,err:=o.DB.Pool.Query(ctx,`SELECT tx_hash,address,asset,COALESCE(contract_address,'') FROM observed_deposits WHERE network=$1 AND block_hash=$2 AND status='credited'`,network,blockHash);if err!=nil{return err};defer rows.Close();for rows.Next(){var txHash,address,asset,contract string;if err:=rows.Scan(&txHash,&address,&asset,&contract);err!=nil{return err};if err:=o.Rollback(ctx,blockchain.ObservedDeposit{Network:network,TxHash:txHash,Address:address,Asset:asset,ContractAddress:contract});err!=nil{return err}};return rows.Err()}
