package workerpool

import (
	"sync"

	"github.com/jbiers/timescale-benchmark/pkg/query"
)

type WorkerPool struct {
	Jobs    []chan query.QueryData
	Workers int
}

func NewWorkerPool(jobs []chan query.QueryData, workers int) *WorkerPool {
	return &WorkerPool{
		Jobs:    jobs,
		Workers: workers,
	}
}

func (wp *WorkerPool) Dispatch() {
	var wg sync.WaitGroup

	for w := 0; w < wp.Workers; w++ {
		wg.Add(1)
		go wp.Worker(w, wp.Jobs[w], &wg)
	}

	wg.Wait()
}

func (wp *WorkerPool) Worker(id int, jobs chan query.QueryData, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		job.Process()
	}
}
