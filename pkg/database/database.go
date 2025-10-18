package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbOnce sync.Once

// TODO: get db info based on env variables
func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgres://postgres:password@localhost:5432/homework")
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %v", err)
	}

	// TODO: properly configure connection pool
	config.MaxConns = 10

	var pool *pgxpool.Pool
	dbOnce.Do(func() {
		pool, err = pgxpool.NewWithConfig(ctx, config)
	})

	if err != nil {
		return nil, fmt.Errorf("error creating database pool: %v", err)
	}

	// TODO: handle contexts
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return pool, nil
}
