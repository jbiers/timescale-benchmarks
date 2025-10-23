package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jbiers/timescale-benchmark/internal/config"
)

var dbOnce sync.Once

func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %v", err)
	}

	poolCfg.MaxConns = int32(config.Workers) + 2
	poolCfg.MinConns = int32(config.Workers)/2 + 1
	poolCfg.MaxConnIdleTime = time.Minute * 10
	poolCfg.MaxConnLifetime = time.Hour * 24
	poolCfg.HealthCheckPeriod = time.Minute * 5

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
