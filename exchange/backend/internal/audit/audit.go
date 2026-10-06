package audit

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
    ID, ActorID, Action, Resource, RequestID string
    MetadataJSON string
    IPHash string
}

type Sink interface { Append(context.Context, Event) error }

func Validate(e Event) error {
    if e.ID == "" || e.ActorID == "" || e.Action == "" || e.Resource == "" || e.RequestID == "" {
        return errors.New("incomplete audit event")
    }
    if len(e.ID) > 128 || len(e.ActorID) > 256 || len(e.Action) > 256 || len(e.Resource) > 512 || len(e.RequestID) > 256 {
        return errors.New("audit field too long")
    }
    return nil
}

type PostgresSink struct { Pool *pgxpool.Pool }

func (s *PostgresSink) Append(ctx context.Context, e Event) error {
    if s == nil || s.Pool == nil { return errors.New("audit database unavailable") }
    if err := Validate(e); err != nil { return err }

    metadata := json.RawMessage(`{}`)
    if e.MetadataJSON != "" {
        metadata = json.RawMessage(e.MetadataJSON)
        if !json.Valid(metadata) { return errors.New("invalid audit metadata JSON") }
    }

    tx, err := s.Pool.Begin(ctx)
    if err != nil { return err }
    defer tx.Rollback(ctx)

    var prev []byte
    _ = tx.QueryRow(ctx, `SELECT event_hash FROM audit_log ORDER BY occurred_at DESC, event_id DESC LIMIT 1`).Scan(&prev)

    now := time.Now().UTC()
    canonical := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
        e.ID, e.ActorID, e.Action, e.Resource, e.RequestID, e.IPHash, now.Format(time.RFC3339Nano))
    h := sha256.New()
    h.Write(prev)
    h.Write([]byte(canonical))
    h.Write(metadata)
    eventHash := h.Sum(nil)

    _, err = tx.Exec(ctx, `INSERT INTO audit_log
        (event_id, occurred_at, actor_id, action, resource, request_id, ip_hash, metadata, prev_hash, event_hash)
        VALUES ($1::uuid,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
        e.ID, now, e.ActorID, e.Action, e.Resource, e.RequestID, nullIfEmpty(e.IPHash),
        metadata, nullIfEmptyBytes(prev), eventHash)
    if err != nil { return err }
    return tx.Commit(ctx)
}

func nullIfEmpty(v string) any {
    if v == "" { return nil }
    return v
}
func nullIfEmptyBytes(v []byte) any {
    if len(v) == 0 { return nil }
    return v
}
func HashIP(ip string) string {
    h := sha256.Sum256([]byte(ip))
    return hex.EncodeToString(h[:])
}
