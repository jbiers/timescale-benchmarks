package query

import (
	"context"
	"hash/fnv"
	"time"

	"github.com/jbiers/timescale-benchmark/pkg/database"
)

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

func (qd *QueryData) RunQuery(ctx context.Context, repo database.Repository) error {
	return repo.ExecuteQuery(ctx, qd.Hostname, qd.StartTime, qd.EndTime)
}
