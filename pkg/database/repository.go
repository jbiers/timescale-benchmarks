package database

import (
	"context"
	"time"
)

type Repository interface {
	ExecuteQuery(ctx context.Context, hostname string, startTime, endTime time.Time) error

	Ping(ctx context.Context) error

	Close()
}
