package events

import("context";"errors";"github.com/Josemilk/wallet-core/exchange/backend/internal/auth")

type CreateSeasonalEvent struct{TemplateID,Title,Asset,Sport,League,EventData,EntryPrice,PrizeAmount,PrizeAsset,StartsAt,EndsAt string}
type AdminStore interface{CreateFromTemplate(context.Context,CreateSeasonalEvent,string)(string,error);UpdateRewards(context.Context,string,string,string,string)error;Carousel(context.Context)([]CarouselItem,error)}
type CarouselItem struct{ID,Title,Type,Icon,Description,EntryPrice,Prize,PrizeAsset,StartsAt,EndsAt string}
type Service struct{Auth auth.Verifier;Store AdminStore}
func(s *Service)Create(ctx context.Context,bearer string,r CreateSeasonalEvent)(string,error){if s==nil||s.Auth==nil||s.Store==nil{return "",errors.New("event admin dependencies missing")};p,err:=auth.Require(ctx,s.Auth,bearer);if err!=nil{return "",err};if !hasAdmin(p){return "",errors.New("admin role required")};if r.TemplateID==""||r.PrizeAmount==""||r.PrizeAsset==""||r.StartsAt==""||r.EndsAt==""{return "",errors.New("incomplete event configuration")};return s.Store.CreateFromTemplate(ctx,r,p.UserID)}
func(s *Service)UpdateRewards(ctx context.Context,bearer,eventID,entryPrice,prizeAmount,prizeAsset string)error{if s==nil||s.Auth==nil||s.Store==nil{return errors.New("event admin dependencies missing")};p,err:=auth.Require(ctx,s.Auth,bearer);if err!=nil{return err};if !hasAdmin(p){return errors.New("admin role required")};if eventID==""||entryPrice==""||prizeAmount==""||prizeAsset==""{return errors.New("incomplete reward configuration")};return s.Store.UpdateRewards(ctx,eventID,entryPrice,prizeAmount,prizeAsset)}
func hasAdmin(p auth.Principal)bool{for _,r:=range p.Roles{if r=="admin"||r=="events_admin"{return true}};return false}
