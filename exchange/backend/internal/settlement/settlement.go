package settlement

import "errors"

type Fill struct { Symbol, MakerAccount, TakerAccount, BaseAsset, QuoteAsset string; Price, Quantity int64 }
type Service struct{}
func New() *Service { return &Service{} }
func (s *Service) Validate(f Fill) error { if f.Symbol==""||f.MakerAccount==""||f.TakerAccount==""||f.BaseAsset==""||f.QuoteAsset==""||f.Price<=0||f.Quantity<=0{return errors.New("invalid fill")}; return nil }
