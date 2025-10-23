package workerpool

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/jbiers/timescale-benchmark/internal/config"
	"github.com/jbiers/timescale-benchmark/pkg/database"
	"github.com/jbiers/timescale-benchmark/pkg/query"
)

type WorkerPool struct {
	jobs         []chan query.QueryData
	workers      int
	repo         database.Repository
	results      []time.Duration
	resultsMutex sync.Mutex
}

func NewWorkerPool(jobs []chan query.QueryData, repo database.Repository) *WorkerPool {
	return &WorkerPool{
		jobs:    jobs,
		workers: config.Workers,
		repo:    repo,
		results: []time.Duration{},
	}
}

func (wp *WorkerPool) Dispatch(ctx context.Context) *WorkerPoolMetrics {
	var wg sync.WaitGroup

	for w := 0; w < wp.workers; w++ {
		wg.Add(1)
		go wp.worker(ctx, w, &wg)
	}

	wg.Wait()

	return wp.getWorkerPoolMetrics()
}

func (wp *WorkerPool) worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range wp.jobs[id] {

		totalTime, err := job.RunQuery(ctx, wp.repo)
		if err != nil {
			config.Logger.Errorf("Worker %d failed to run query: %v. The processing time will not be included in the result metrics.", id, err)
			return
		}

		config.Logger.Debugf("Worker %d - Query Params - (Hostname: %s, StartTime: %v, EndTime: %v) - Total processing time: %v", id, job.Hostname, job.StartTime, job.EndTime, totalTime)

		wp.resultsMutex.Lock()
		wp.results = append(wp.results, totalTime)
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
