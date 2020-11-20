package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// database driver
	_ "github.com/lib/pq"
)

const (
	dbTimeoutShort = time.Second
	dbTimoutLong   = time.Minute
)

// Connect connects to database
func Connect(ctx context.Context) (*sql.DB, error) {
	conn, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/demoservice?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, dbTimeoutShort)
	defer cancel()

	if err := conn.PingContext(ctxPing); err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	return conn, nil
}
