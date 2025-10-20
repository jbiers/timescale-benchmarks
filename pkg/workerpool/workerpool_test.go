package workerpool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWorkerPoolMetrics_EmptyResults(t *testing.T) {
	wp := &WorkerPool{
		results: []time.Duration{},
	}

	metrics := wp.getWorkerPoolMetrics()

	assert.Equal(t, 0, metrics.ProcessedJobs)
	assert.Equal(t, time.Duration(0), metrics.LongestTime)
	assert.Equal(t, time.Duration(0), metrics.ShortestTime)
	assert.Equal(t, time.Duration(0), metrics.TotalTime)
	assert.Equal(t, time.Duration(0), metrics.AverageTime)
	assert.Equal(t, time.Duration(0), metrics.MedianTime)
}

func TestGetWorkerPoolMetrics_BasicCalculations(t *testing.T) {
	results := []time.Duration{
		10 * time.Millisecond,
		30 * time.Millisecond,
		20 * time.Millisecond,
	}

	wp := &WorkerPool{
		results: results,
	}

	metrics := wp.getWorkerPoolMetrics()

	assert.Equal(t, 3, metrics.ProcessedJobs)
	assert.Equal(t, 10*time.Millisecond, metrics.ShortestTime)
	assert.Equal(t, 30*time.Millisecond, metrics.LongestTime)
	assert.Equal(t, 60*time.Millisecond, metrics.TotalTime)
	assert.Equal(t, 20*time.Millisecond, metrics.AverageTime)
	assert.Equal(t, 20*time.Millisecond, metrics.MedianTime)
}

func TestGetWorkerPoolMetrics_EvenCountMedian(t *testing.T) {
	results := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
	}

	wp := &WorkerPool{
		results: results,
	}

	metrics := wp.getWorkerPoolMetrics()

	assert.Equal(t, 25*time.Millisecond, metrics.MedianTime)
}
