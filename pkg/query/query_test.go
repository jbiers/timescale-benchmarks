package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHash_Distribution(t *testing.T) {
	queryData := []QueryData{
		{Hostname: "host_000001"},
		{Hostname: "host_000002"},
		{Hostname: "host_000003"},
		{Hostname: "host_000004"},
		{Hostname: "host_000005"},
		{Hostname: "host_000010"},
		{Hostname: "host_000100"},
		{Hostname: "host_001000"},
	}

	tests := []struct {
		name        string
		workerCount int
		expectedMin int
		expectedMax int
	}{
		{
			name:        "single worker",
			workerCount: 1,
			expectedMin: 0,
			expectedMax: 0,
		},
		{
			name:        "two workers",
			workerCount: 2,
			expectedMin: 0,
			expectedMax: 1,
		},
		{
			name:        "five workers",
			workerCount: 5,
			expectedMin: 0,
			expectedMax: 4,
		},
		{
			name:        "ten workers",
			workerCount: 10,
			expectedMin: 0,
			expectedMax: 9,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			workerAssignments := make(map[int]int)

			for _, qd := range queryData {
				hash := qd.GetHash(testCase.workerCount)

				assert.GreaterOrEqual(t, hash, testCase.expectedMin,
					"Hash value %d should be >= %d for worker count %d", hash, testCase.expectedMin, testCase.workerCount)
				assert.LessOrEqual(t, hash, testCase.expectedMax,
					"Hash value %d should be <= %d for worker count %d", hash, testCase.expectedMax, testCase.workerCount)

				workerAssignments[hash]++
			}

			if testCase.workerCount > 1 {
				assert.Greater(t, len(workerAssignments), 1,
					"With %d workers, should use more than 1 worker for load balancing", testCase.workerCount)
			}

			for _, count := range workerAssignments {
				assert.GreaterOrEqual(t, count, 1, "Each worker should be assigned at least one query")
			}
		})
	}
}
