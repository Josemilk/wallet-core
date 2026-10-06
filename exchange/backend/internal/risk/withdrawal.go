package risk

import (
    "context"
    "errors"
    "math/big"
)

type WithdrawalLimits struct { MaxAmount string; RequireReviewAbove string }

type WithdrawalChecker struct { Limits map[string]WithdrawalLimits }

func (c *WithdrawalChecker) Check(ctx context.Context, userID, asset, amount string) error {
    if c == nil { return errors.New("risk checker unavailable") }
    l, ok := c.Limits[asset]; if !ok { return errors.New("asset withdrawal disabled") }
    a, ok := new(big.Int).SetString(amount,10); if !ok || a.Sign() <= 0 { return errors.New("invalid withdrawal amount") }
    max, ok := new(big.Int).SetString(l.MaxAmount,10); if !ok || max.Sign() <= 0 { return errors.New("invalid configured withdrawal limit") }
    if a.Cmp(max) > 0 { return errors.New("withdrawal exceeds risk limit") }
    return nil
}
