package workerpool

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jbiers/timescale-benchmark/config"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	jobs         []chan query.QueryData
	workers      int
	dbPool       *pgxpool.Pool
	results      []time.Duration
	resultsMutex sync.Mutex
}

func NewWorkerPool(jobs []chan query.QueryData, dbPool *pgxpool.Pool) *WorkerPool {
	return &WorkerPool{
		jobs:    jobs,
		workers: config.Workers,
		dbPool:  dbPool,
		results: []time.Duration{},
	}
}

func (wp *WorkerPool) Dispatch() *WorkerPoolMetrics {
	var wg sync.WaitGroup

	for w := 0; w < wp.workers; w++ {
		wg.Add(1)
		go wp.worker(w, &wg)
	}

	wg.Wait()

	return wp.getWorkerPoolMetrics()
}

func (wp *WorkerPool) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range wp.jobs[id] {
		start := time.Now()

		err := job.RunQuery(context.Background(), wp.dbPool)
		if err != nil {
			logrus.Errorf("worker %d failed to run query: %v", id, err)
			return
		}
		total := time.Since(start)

		wp.resultsMutex.Lock()
		wp.results = append(wp.results, total)
		wp.resultsMutex.Unlock()
	}
}

func (wp *WorkerPool) getWorkerPoolMetrics() *WorkerPoolMetrics {
	if len(wp.results) == 0 {
		return &WorkerPoolMetrics{}
	}

	slices.Sort(wp.results)

	processedJobs := len(wp.results)

	var totalTime time.Duration
	for _, result := range wp.results {
		totalTime += result
	}
	var medianTime time.Duration
	if processedJobs%2 == 1 {
		medianTime = wp.results[processedJobs/2]
	} else {
		mid1 := wp.results[(processedJobs/2)-1]
		mid2 := wp.results[processedJobs/2]
		medianTime = (mid1 + mid2) / 2
	}

	return &WorkerPoolMetrics{
		ProcessedJobs: processedJobs,
		LongestTime:   wp.results[processedJobs-1],
		ShortestTime:  wp.results[0],
		TotalTime:     totalTime,
		AverageTime:   totalTime / time.Duration(processedJobs),
		MedianTime:    medianTime,
	}
}
