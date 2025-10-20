package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) ExecuteQuery(ctx context.Context, hostname string, startTime, endTime time.Time) error {
	_, err := r.pool.Exec(ctx, Query, hostname, startTime, endTime)
	if err != nil {
		return err
	}
	return nil
}

const Query = `
SELECT
  time_bucket('1 minute', ts, $2::TIMESTAMPTZ) AS minutes,
  MAX(usage) AS max_usage,
  MIN(usage) AS min_usage
FROM cpu_usage
WHERE host = $1
  AND ts > $2
  AND ts < $3
GROUP BY minutes
ORDER BY minutes;
`

func (r *PostgresRepository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}
