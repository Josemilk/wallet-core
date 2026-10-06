package ws

import (
    "context"
    "errors"
    "sync"
)

type Authenticator interface { Authenticate(ctx context.Context, token string) (string, error) }
type Client struct { UserID string; Send chan []byte }
type Hub struct { mu sync.RWMutex; clients map[string]map[*Client]struct{} }

func NewHub() *Hub { return &Hub{clients: make(map[string]map[*Client]struct{})} }

func (h *Hub) Register(ctx context.Context, auth Authenticator, token string, c *Client) error {
    if h == nil || c == nil || auth == nil { return errors.New("invalid websocket registration") }
    userID, err := auth.Authenticate(ctx, token); if err != nil { return err }
    if userID == "" { return errors.New("unauthenticated websocket") }
    c.UserID = userID
    h.mu.Lock(); defer h.mu.Unlock()
    if h.clients[userID] == nil { h.clients[userID] = make(map[*Client]struct{}) }
    h.clients[userID][c] = struct{}{}
    return nil
}

func (h *Hub) Unregister(c *Client) {
    if h == nil || c == nil { return }
    h.mu.Lock(); defer h.mu.Unlock()
    if set := h.clients[c.UserID]; set != nil { delete(set, c); if len(set)==0 { delete(h.clients,c.UserID) } }
}

func (h *Hub) PublishUser(userID string, payload []byte) {
    h.mu.RLock(); defer h.mu.RUnlock()
    for c := range h.clients[userID] { select { case c.Send <- payload: default: /* slow client; drop instead of blocking engine */ } }
}
