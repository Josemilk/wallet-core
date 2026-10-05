package api

import (
 "encoding/json"
 "net/http"
 "strings"
)

type Handler struct { Service *Service; AssetIDs []string }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 if h == nil || h.Service == nil { http.Error(w, "service unavailable", http.StatusServiceUnavailable); return }
 if r.Method != http.MethodGet || r.URL.Path != "/v1/dashboard" { http.NotFound(w, r); return }
 authz := strings.TrimSpace(r.Header.Get("Authorization"))
 bearer := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
 dashboard, err := h.Service.Load(r.Context(), bearer, h.AssetIDs)
 if err != nil { http.Error(w, "unauthorized or unavailable", http.StatusUnauthorized); return }
 w.Header().Set("Content-Type", "application/json")
 w.Header().Set("Cache-Control", "no-store")
 _ = json.NewEncoder(w).Encode(dashboard)
}
