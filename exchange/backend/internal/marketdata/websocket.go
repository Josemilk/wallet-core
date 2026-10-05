package marketdata

import("context";"github.com/Josemilk/wallet-core/exchange/backend/internal/ws")

type Publisher interface{PublishUser(string,[]byte)}
type UserStream struct{Hub Publisher}
func(s *UserStream) Publish(ctx context.Context,userID string,payload []byte){if s==nil||s.Hub==nil||userID==""{return};s.Hub.PublishUser(userID,payload)}
var _ = ws.Client{}
