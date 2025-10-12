package query

import "time"

type QueryData struct {
	Hostname  string
	StartTime time.Time
	EndTime   time.Time
}
