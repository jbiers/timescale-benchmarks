package query

import (
	"context"
	"hash/fnv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

type QueryData struct {
	Hostname  string
	StartTime time.Time
	EndTime   time.Time
}

func (qd *QueryData) GetHash(num int) int {
	h := fnv.New32a()
	h.Write([]byte(qd.Hostname))

	return int(h.Sum32()) % num
}

func (qd *QueryData) RunQuery(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, Query, qd.Hostname, qd.StartTime, qd.EndTime)
	if err != nil {
		return err
	}

	return nil
}
