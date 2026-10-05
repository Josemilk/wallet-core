package custody

import("context";"errors")

type SigningGuard struct{Signer Signer; AllowedNetworks map[string]bool}
func(g *SigningGuard) Sign(ctx context.Context,r SignRequest)(string,error){if g==nil||g.Signer==nil{return "",errors.New("custody signer unavailable")};if !g.AllowedNetworks[r.Network]{return "",errors.New("network not enabled for custody")};if r.KeyRef==""||r.RawTransaction==""{return "",errors.New("missing signing material reference")};return g.Signer.Sign(ctx,r)}
