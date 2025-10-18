package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jbiers/timescale-benchmark/config"
)

var dbOnce sync.Once

// TODO: get db info based on env variables
func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig("postgres://postgres:password@localhost:5432/homework")
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %v", err)
	}

	poolCfg.MaxConns = int32(config.Workers)
	poolCfg.MinConns = int32(config.Workers)
	poolCfg.MaxConnIdleTime = time.Minute * 10

	var pool *pgxpool.Pool
	dbOnce.Do(func() {
		pool, err = pgxpool.NewWithConfig(ctx, poolCfg)
	})

	if err != nil {
		return nil, fmt.Errorf("error creating database pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return pool, nil
}
