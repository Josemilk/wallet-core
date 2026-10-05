package auth

import (
    "context"
    "errors"
)

type Principal struct { UserID string; Roles []string }
type Verifier interface { Verify(ctx context.Context, bearer string) (Principal,error) }
func Require(ctx context.Context, v Verifier, bearer string) (Principal,error) { if v==nil{return Principal{},errors.New("auth verifier unavailable")}; if bearer==""{return Principal{},errors.New("missing bearer token")}; p,err:=v.Verify(ctx,bearer); if err!=nil{return Principal{},err}; if p.UserID==""{return Principal{},errors.New("invalid principal")}; return p,nil }
