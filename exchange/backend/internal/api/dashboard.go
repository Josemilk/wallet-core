package api

import (
 "context"
 "errors"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/auth"
 "github.com/Josemilk/wallet-core/exchange/backend/internal/marketdata"
)

type Dashboard struct { UserID string `json:"userId"`; Portfolio Portfolio `json:"portfolio"`; Markets []marketdata.Quote `json:"markets"`; Events []EventSummary `json:"events"` }
type Portfolio struct { TotalUSD string `json:"totalUsd"`; InvestedUSD string `json:"investedUsd"`; ProfitUSD string `json:"profitUsd"`; DayPnLUSD string `json:"dayPnlUsd"` }
type EventSummary struct { ID string `json:"id"`; Title string `json:"title"`; Type string `json:"type"`; EntryPrice string `json:"entryPrice"`; Prize string `json:"prize"`; StartsAt string `json:"startsAt"`; EndsAt string `json:"endsAt"` }
type DashboardSource interface { Portfolio(ctx context.Context,userID string)(Portfolio,error); Events(ctx context.Context,userID string)([]EventSummary,error) }
type Service struct { Auth auth.Verifier; Market marketdata.Provider; Source DashboardSource }
func(s *Service) Load(ctx context.Context,bearer string,assetIDs []string) (Dashboard,error){if s==nil||s.Auth==nil||s.Market==nil||s.Source==nil{return Dashboard{},errors.New("dashboard service dependencies missing")};p,err:=auth.Require(ctx,s.Auth,bearer);if err!=nil{return Dashboard{},err};pf,err:=s.Source.Portfolio(ctx,p.UserID);if err!=nil{return Dashboard{},err};ev,err:=s.Source.Events(ctx,p.UserID);if err!=nil{return Dashboard{},err};quotes,err:=s.Market.Quotes(ctx,assetIDs,"USD");if err!=nil{return Dashboard{},err};return Dashboard{UserID:p.UserID,Portfolio:pf,Markets:quotes,Events:ev},nil}
