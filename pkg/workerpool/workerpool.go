package workerpool

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	Jobs    []chan query.QueryData
	Workers int
	DBPool  *pgxpool.Pool
}

func NewWorkerPool(jobs []chan query.QueryData, workers int, dbPool *pgxpool.Pool) *WorkerPool {
	return &WorkerPool{
		Jobs:    jobs,
		Workers: workers,
		DBPool:  dbPool,
	}
}

func (wp *WorkerPool) Dispatch() {
	var wg sync.WaitGroup

	for w := 0; w < wp.Workers; w++ {
		wg.Add(1)
		go wp.Worker(w, wp.Jobs[w], &wg, wp.DBPool)
	}

	wg.Wait()
}

func (wp *WorkerPool) Worker(id int, jobs chan query.QueryData, wg *sync.WaitGroup, dbPool *pgxpool.Pool) {
	defer wg.Done()

	for job := range jobs {
		//start := time.Now()
		err := job.RunQuery(context.Background(), dbPool)
		if err != nil {
			logrus.Errorf("worker %d failed to run query: %v", id, err)

		}
		// total := time.Since(start)
		// add total to an array (use mutex)
		// calculate results from array
		// print results in the end
	}
}
