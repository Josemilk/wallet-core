package withdrawals

import (
    "context"
    "errors"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/custody"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/risk"
)

type Request struct { ID, UserID, Asset, Destination, Amount, Network, RawTransaction, KeyRef string }
type Service struct { Risk risk.Engine; Signer custody.Signer }

func (s *Service) Execute(ctx context.Context, req Request) (string, error) {
    if s == nil || s.Risk == nil || s.Signer == nil { return "", errors.New("withdrawal service dependencies missing") }
    if req.ID == "" || req.UserID == "" || req.Asset == "" || req.Destination == "" || req.Amount == "" || req.Network == "" || req.KeyRef == "" || req.RawTransaction == "" { return "", errors.New("invalid withdrawal request") }
    if err := s.Risk.Check(ctx, req.UserID, req.Asset, req.Amount); err != nil { return "", err }
    return s.Signer.Sign(ctx, custody.SignRequest{KeyRef:req.KeyRef, Network:req.Network, RawTransaction:req.RawTransaction})
}
