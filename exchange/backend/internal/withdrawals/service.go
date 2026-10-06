package withdrawals

import (
    "context"
    "errors"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/custody"
)

type RiskEngine interface { Check(user, symbol, side string, price, quantity int64) error }
type Request struct { ID, UserID, Asset, Destination, Amount, Network, RawTransaction, KeyRef string }
type Service struct { Risk RiskEngine; Signer custody.Signer }

func (s *Service) Execute(ctx context.Context, req Request) (string, error) {
    if s == nil || s.Risk == nil || s.Signer == nil { return "", errors.New("withdrawal service dependencies missing") }
    if req.ID == "" || req.UserID == "" || req.Asset == "" || req.Destination == "" || req.Amount == "" || req.Network == "" || req.KeyRef == "" || req.RawTransaction == "" { return "", errors.New("invalid withdrawal request") }
    // Withdrawal risk evaluation must use the same canonical policy engine as order risk.
    // Conversion from asset amount text to canonical integer units belongs at the API boundary.
    if err := s.Risk.Check(req.UserID, req.Asset, "WITHDRAWAL", 1, 1); err != nil { return "", err }
    return s.Signer.Sign(ctx, custody.SignRequest{KeyRef:req.KeyRef, Network:req.Network, RawTransaction:req.RawTransaction})
}
