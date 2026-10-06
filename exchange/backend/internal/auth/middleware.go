package auth

import (
    "context"
    "errors"
    "strings"
)

type Principal struct { UserID string; Roles []string }
type Verifier interface { Verify(ctx context.Context, bearer string) (Principal,error) }

func Require(ctx context.Context, v Verifier, bearer string) (Principal,error) {
    if v == nil { return Principal{}, errors.New("auth verifier unavailable") }
    if bearer == "" { return Principal{}, errors.New("missing bearer token") }
    p, err := v.Verify(ctx, bearer)
    if err != nil { return Principal{}, err }
    if p.UserID == "" { return Principal{}, errors.New("invalid principal") }
    return p,nil
}

func HasRole(p Principal, required ...string) bool {
    for _, want := range required {
        for _, got := range p.Roles {
            if strings.EqualFold(strings.TrimSpace(got), strings.TrimSpace(want)) { return true }
        }
    }
    return false
}

func RequireRole(p Principal, required ...string) error {
    if !HasRole(p, required...) { return errors.New("forbidden: required role missing") }
    return nil
}

func IsAdmin(p Principal) bool { return HasRole(p, "admin", "security") }
