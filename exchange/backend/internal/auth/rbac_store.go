package auth

import (
    "context"
    "errors"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"

)

type RoleStore struct { Pool *pgxpool.Pool }

func (s *RoleStore) Roles(ctx context.Context, userID string) ([]string, error) {
    if s == nil || s.Pool == nil { return nil, errors.New("role store unavailable") }
    userID = strings.TrimSpace(userID)
    if userID == "" { return nil, errors.New("user id required") }
    rows, err := s.Pool.Query(ctx, "SELECT role_name FROM user_roles WHERE user_id=$1 ORDER BY role_name", userID)
    if err != nil { return nil, err }
    defer rows.Close()
    roles := make([]string, 0)
    for rows.Next() {
        var role string
        if err := rows.Scan(&role); err != nil { return nil, err }
        roles = append(roles, role)
    }
    if err := rows.Err(); err != nil { return nil, err }
    return roles, nil
}

func MergeRoles(tokenRoles, dbRoles []string) []string {
    seen := map[string]struct{}{}
    out := make([]string, 0, len(tokenRoles)+len(dbRoles))
    for _, list := range [][]string{tokenRoles, dbRoles} {
        for _, role := range list {
            role = strings.TrimSpace(role)
            if role == "" { continue }
            key := strings.ToLower(role)
            if _, ok := seen[key]; ok { continue }
            seen[key] = struct{}{}
            out = append(out, role)
        }
    }
    return out
}
