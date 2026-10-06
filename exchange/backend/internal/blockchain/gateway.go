package blockchain

import "context"

type Tx struct { Network, Hash, From, To, Asset string; Amount string; Confirmations uint64 }
type Gateway interface { LatestBlock(ctx context.Context) (uint64,error); Transaction(ctx context.Context, hash string) (Tx,error); Broadcast(ctx context.Context, rawTx string) (string,error) }

// Implementations must use allowlisted RPC endpoints and verify chain/network identity.
// No private key material belongs in this package.
