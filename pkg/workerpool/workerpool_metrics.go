package workerpool

import (
	"fmt"
	"time"
)

type WorkerPoolMetrics struct {
	MedianTime    time.Duration
	AverageTime   time.Duration
	LongestTime   time.Duration
	ShortestTime  time.Duration
	TotalTime     time.Duration
	ProcessedJobs int
}

func (wpm *WorkerPoolMetrics) ReportWorkerPoolMetrics() {
	fmt.Println("==== Worker Pool Metrics ====")
	fmt.Printf("Processed Jobs : %d\n", wpm.ProcessedJobs)
	fmt.Printf("Total Time     : %v\n", wpm.TotalTime)
	fmt.Printf("Average Time   : %v\n", wpm.AverageTime)
	fmt.Printf("Median Time    : %v\n", wpm.MedianTime)
	fmt.Printf("Longest Time   : %v\n", wpm.LongestTime)
	fmt.Printf("Shortest Time  : %v\n", wpm.ShortestTime)
}
