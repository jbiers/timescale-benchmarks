package query

import (
	"fmt"
	"hash/fnv"
	"time"
)

type QueryData struct {
	Hostname  string
	StartTime time.Time
	EndTime   time.Time
}

func (qd *QueryData) GetIndex(maxIndex int) int {
	h := fnv.New32a()
	h.Write([]byte(qd.Hostname))

	return int(h.Sum32()) % maxIndex
}

func (qd *QueryData) Process() {
	fmt.Println()
}
