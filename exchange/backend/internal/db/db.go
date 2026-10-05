package db

import (
    "context"
    "errors"
    "os"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DB struct { Pool *pgxpool.Pool }

func Open(ctx context.Context) (*DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" { return nil, errors.New("DATABASE_URL is required") }
    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil { return nil, err }
    cfg.MaxConns = 20
    cfg.MinConns = 2
    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil { return nil, err }
    if err := pool.Ping(ctx); err != nil { pool.Close(); return nil, err }
    return &DB{Pool: pool}, nil
}

func (d *DB) Close() { if d != nil && d.Pool != nil { d.Pool.Close() } }

func (d *DB) Health(ctx context.Context) error { if d == nil || d.Pool == nil { return errors.New("database not initialized") }; return d.Pool.Ping(ctx) }

func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
    tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
    if err != nil { return err }
    defer tx.Rollback(ctx)
    if err := fn(tx); err != nil { return err }
    return tx.Commit(ctx)
}
