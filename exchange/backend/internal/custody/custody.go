package custody

import "context"

type SignRequest struct { KeyRef, Network, RawTransaction string }
type Signer interface { Sign(ctx context.Context, req SignRequest) (string,error) }

// Production implementations must resolve KeyRef inside an HSM/KMS-backed signer.
// Raw private keys, seed phrases and decrypted key blobs must never cross this interface.
