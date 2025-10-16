package workerpool

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	Jobs         []chan query.QueryData
	Workers      int
	DBPool       *pgxpool.Pool
	Results      []time.Duration
	ResultsMutex sync.Mutex
}

func NewWorkerPool(jobs []chan query.QueryData, workers int, dbPool *pgxpool.Pool) *WorkerPool {
	return &WorkerPool{
		Jobs:    jobs,
		Workers: workers,
		DBPool:  dbPool,
		Results: []time.Duration{},
	}
}

func (wp *WorkerPool) Dispatch() {
	var wg sync.WaitGroup

	for w := 0; w < wp.Workers; w++ {
		wg.Add(1)
		go wp.Worker(w, &wg)
	}

	wg.Wait()
}

func (wp *WorkerPool) Worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range wp.Jobs[id] {
		start := time.Now()

		err := job.RunQuery(context.Background(), wp.DBPool)
		if err != nil {
			logrus.Errorf("worker %d failed to run query: %v", id, err)
			return
		}
		total := time.Since(start)

		wp.ResultsMutex.Lock()
		wp.Results = append(wp.Results, total)
		wp.ResultsMutex.Unlock()
	}
}
