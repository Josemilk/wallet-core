package ws

import("context";"net/http";"strings";"github.com/Josemilk/wallet-core/exchange/backend/internal/auth";"github.com/gorilla/websocket")

type VerifierAdapter struct{Verifier auth.Verifier}
func(a VerifierAdapter)Authenticate(ctx context.Context,token string)(string,error){p,err:=auth.Require(ctx,a.Verifier,token);if err!=nil{return "",err};return p.UserID,nil}

type HTTPHandler struct{Auth auth.Verifier;Hub *Hub;Upgrader websocket.Upgrader}
func(h *HTTPHandler)ServeHTTP(w http.ResponseWriter,r *http.Request){if h==nil||h.Auth==nil||h.Hub==nil{http.Error(w,"service unavailable",503);return};if r.Method!=http.MethodGet||r.URL.Path!="/v1/ws"{http.NotFound(w,r);return};token:=strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "));if token==""{http.Error(w,"unauthorized",401);return};p,err:=auth.Require(r.Context(),h.Auth,token);if err!=nil{http.Error(w,"unauthorized",401);return};up:=h.Upgrader;if up.CheckOrigin==nil{up.CheckOrigin=func(*http.Request)bool{return true}};conn,err:=up.Upgrade(w,r,nil);if err!=nil{return};client:=&Client{UserID:p.UserID,Send:make(chan []byte,32)};if err:=h.Hub.Register(r.Context(),VerifierAdapter{Verifier:h.Auth},token,client);err!=nil{_ = conn.Close();return};done:=make(chan struct{});go func(){defer close(done);for payload:=range client.Send{if err:=conn.WriteMessage(websocket.TextMessage,payload);err!=nil{return}}}();for{if _,_,err:=conn.ReadMessage();err!=nil{break}};h.Hub.Unregister(client);close(client.Send);<-done;_ = conn.Close()}
