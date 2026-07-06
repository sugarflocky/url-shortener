// Package postgres implements storage and defines methods to save URL - code links.
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugarflocky/url-shortener/internal/storage"
)

// schema is applied on startup; IF NOT EXISTS makes it idempotent.
const schema = `
CREATE TABLE IF NOT EXISTS links (
    code   TEXT PRIMARY KEY,
    url TEXT NOT NULL UNIQUE
)`

// Storage is a postgres implementation of the code - url links storage.
type Storage struct {
	pool *pgxpool.Pool
}

// uniqueViolation is the postgres error code for unique constraint violations.
const uniqueViolation = "23505"

// New connects to postgres by dsn, applies the schema
// and returns a ready-to-use Storage.
func New(ctx context.Context, dsn string) (*Storage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if _, err := pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &Storage{pool: pool}, nil
}

// Save stores the code - url pair.
// It returns storage.ErrCodeTaken if the code is taken,
// and storage.ErrURLExists if the url already has an owner.
func (s *Storage) Save(ctx context.Context, url string, code string) error {
	_, err := s.pool.Exec(ctx,
		"INSERT INTO links (code, url) VALUES ($1, $2)", code, url)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		switch pgErr.ConstraintName {
		case "links_pkey":
			return storage.ErrCodeTaken
		case "links_url_key":
			return storage.ErrURLExists
		}
	}
	return fmt.Errorf("insert pair: %w", err)
}

// URLByCode returns URL by given code.
// Returns storage.ErrNotFound if there is no pair.
func (s *Storage) URLByCode(ctx context.Context, code string) (url string, err error) {
	err = s.pool.QueryRow(ctx,
		"SELECT url FROM links WHERE code = $1", code).Scan(&url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("select url: %w", err)
	}
	return url, nil
}

// CodeByURL returns code by given URL.
// Returns storage.ErrNotFound if there is no pair.
func (s *Storage) CodeByURL(ctx context.Context, url string) (code string, err error) {
	err = s.pool.QueryRow(ctx,
		"SELECT code FROM links WHERE url = $1", url).Scan(&code)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("select code: %w", err)
	}
	return code, nil
}
