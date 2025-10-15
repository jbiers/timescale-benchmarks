package query

import (
	"hash/fnv"
	"time"
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

func (qd *QueryData) Process() {

}
