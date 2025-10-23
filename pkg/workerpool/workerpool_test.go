package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jbiers/timescale-benchmark/internal/config"
	"github.com/jbiers/timescale-benchmark/pkg/database"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestDispatch_ProcessesJobs(t *testing.T) {
	config.InitLogger()
	config.Workers = 1

	mockRepo := &database.MockRepository{}

	jobs := make([]chan query.QueryData, config.Workers)
	jobs[0] = make(chan query.QueryData, config.Workers)

	wp := NewWorkerPool(jobs, mockRepo)

	mockRepo.On("ExecuteQuery", mock.Anything, "host_000001", mock.Anything, mock.Anything).Return(nil)

	queryData := query.QueryData{
		Hostname:  "host_000001",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}
	jobs[0] <- queryData
	close(jobs[0])

	ctx := context.Background()
	metrics := wp.Dispatch(ctx)

	assert.Equal(t, 1, metrics.ProcessedJobs)
	assert.Greater(t, metrics.TotalTime, time.Duration(0))
	mockRepo.AssertExpectations(t)
}

func TestDispatch_HandlesErrors(t *testing.T) {
	config.InitLogger()
	config.Workers = 1

	mockRepo := &database.MockRepository{}

	jobs := make([]chan query.QueryData, config.Workers)
	jobs[0] = make(chan query.QueryData, config.Workers)

	wp := NewWorkerPool(jobs, mockRepo)

	expectedError := errors.New("database connection failed")
	mockRepo.On("ExecuteQuery", mock.Anything, "host_000001", mock.Anything, mock.Anything).Return(expectedError)

	queryData := query.QueryData{
		Hostname:  "host_000001",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}
	jobs[0] <- queryData
	close(jobs[0])

	ctx := context.Background()
	metrics := wp.Dispatch(ctx)

	assert.Equal(t, 0, metrics.ProcessedJobs)
	mockRepo.AssertExpectations(t)
}
