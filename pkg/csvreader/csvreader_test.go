package csvreader

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jbiers/timescale-benchmark/internal/config"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/stretchr/testify/assert"
)

func TestValidateHeaderLine(t *testing.T) {
	tests := []struct {
		name        string
		header      []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid header",
			header:      []string{"hostname", "start_time", "end_time"},
			expectError: false,
		},
		{
			name:        "wrong number of fields - too few",
			header:      []string{"hostname", "start_time"},
			expectError: true,
			errorMsg:    "expected three entries in header line, found 2",
		},
		{
			name:        "wrong number of fields - too many",
			header:      []string{"hostname", "start_time", "end_time", "extra"},
			expectError: true,
			errorMsg:    "expected three entries in header line, found 4",
		},
		{
			name:        "wrong first field name",
			header:      []string{"host", "start_time", "end_time"},
			expectError: true,
			errorMsg:    "expected first entry in header line to match 'hostname', found 'host'",
		},
		{
			name:        "wrong second field name",
			header:      []string{"hostname", "start", "end_time"},
			expectError: true,
			errorMsg:    "expected second entry in header line to match 'start_time', found 'start'",
		},
		{
			name:        "wrong third field name",
			header:      []string{"hostname", "start_time", "end"},
			expectError: true,
			errorMsg:    "expected third entry in header line to match 'end_time', found 'end'",
		},
		{
			name:        "empty header",
			header:      []string{},
			expectError: true,
			errorMsg:    "expected three entries in header line, found 0",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateHeaderLine(testCase.header)

			if testCase.expectError {
				assert.Error(t, err, "Expected an error but got nil")
				if testCase.errorMsg != "" {
					assert.Equal(t, err.Error(), testCase.errorMsg, "Expected error message to be %s, got %s", testCase.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}

func TestParseQueryData(t *testing.T) {
	mustParseTime := func(timeString string) time.Time {
		timeValue, err := time.Parse(timeStringFormat, timeString)
		if err != nil {
			panic(fmt.Errorf("failed to parse time as time.Time format: %w", err))
		}
		return timeValue
	}

	tests := []struct {
		name          string
		record        []string
		expectedQuery query.QueryData
		expectError   bool
		errorMsg      string
	}{
		{
			name:   "valid record",
			record: []string{"host_000001", "2017-01-01 08:59:22", "2017-01-01 09:59:22"},
			expectedQuery: query.QueryData{
				Hostname:  "host_000001",
				StartTime: mustParseTime("2017-01-01 08:59:22"),
				EndTime:   mustParseTime("2017-01-01 09:59:22"),
			},
			expectError: false,
		},
		{
			name:        "wrong number of fields - too few",
			record:      []string{"host_000001", "2017-01-01 08:59:22"},
			expectError: true,
			errorMsg:    "expected 3 entries in query data record, found 2",
		},
		{
			name:        "wrong number of fields - too many",
			record:      []string{"host_000001", "2017-01-01 08:59:22", "2017-01-01 09:59:22", "extra"},
			expectError: true,
			errorMsg:    "expected 3 entries in query data record, found 4",
		},
		{
			name:        "invalid start_time format",
			record:      []string{"host_000001", "2017-01-01T08:59:22", "2017-01-01 09:59:22"},
			expectError: true,
			errorMsg:    "failed to parse 'start_time' '2017-01-01T08:59:22' as time.Time format: parsing time",
		},
		{
			name:        "invalid end_time format",
			record:      []string{"host_000001", "2017-01-01 08:59:22", "2017-01-01T09:59:22"},
			expectError: true,
			errorMsg:    "failed to parse 'end_time' '2017-01-01T09:59:22' as time.Time format: parsing time",
		},
		{
			name:   "empty hostname",
			record: []string{"", "2017-01-01 08:59:22", "2017-01-01 09:59:22"},
			expectedQuery: query.QueryData{
				Hostname:  "",
				StartTime: mustParseTime("2017-01-01 08:59:22"),
				EndTime:   mustParseTime("2017-01-01 09:59:22"),
			},
			expectError: false,
		},
		{
			name:        "empty start_time",
			record:      []string{"host_000001", "", "2017-01-01 09:59:22"},
			expectError: true,
			errorMsg:    "failed to parse 'start_time' '' as time.Time format: parsing time",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := parseQueryData(testCase.record)

			if testCase.expectError {
				assert.Error(t, err, "Expected an error but got nil")
				if testCase.errorMsg != "" {
					assert.Contains(t, err.Error(), testCase.errorMsg, "Expected error message to contain %s, got %s", testCase.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
				assert.Equal(t, testCase.expectedQuery.Hostname, result.Hostname)
				assert.True(t, testCase.expectedQuery.StartTime.Equal(result.StartTime),
					"StartTime mismatch: expected %v, got %v", testCase.expectedQuery.StartTime, result.StartTime)
				assert.True(t, testCase.expectedQuery.EndTime.Equal(result.EndTime),
					"EndTime mismatch: expected %v, got %v", testCase.expectedQuery.EndTime, result.EndTime)
			}
		})
	}
}

func TestStream(t *testing.T) {
	tests := []struct {
		name          string
		csvContent    string
		expectError   bool
		errorMsg      string
		expectedCount int
	}{
		{
			name: "valid CSV with multiple records",
			csvContent: `hostname,start_time,end_time
host_000001,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000002,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000003,2017-01-01 08:59:22,2017-01-01 09:59:22`,
			expectError:   false,
			expectedCount: 3,
		},
		{
			name: "valid CSV with single record",
			csvContent: `hostname,start_time,end_time
host_000001,2017-01-01 08:59:22,2017-01-01 09:59:22`,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:          "empty CSV (only header)",
			csvContent:    `hostname,start_time,end_time`,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name: "invalid header - wrong field name",
			csvContent: `host,start_time,end_time
host_000001,2017-01-01 08:59:22,2017-01-01 09:59:22`,
			expectError: true,
			errorMsg:    "expected first entry in header line to match 'hostname'",
		},
		{
			name: "invalid header - wrong number of fields",
			csvContent: `hostname,start_time
host_000001,2017-01-01 08:59:22`,
			expectError: true,
			errorMsg:    "wrong number of fields",
		},
		{
			name: "valid header but invalid record",
			csvContent: `hostname,start_time,end_time
host_000001,invalid_time,2017-01-01 09:59:22`,
			expectError: true,
			errorMsg:    "failed to parse 'start_time'",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tmpFile := createTempCSVFile(t, testCase.csvContent)
			defer tmpFile.Close()

			channels := make([]chan query.QueryData, 2)
			for i := range channels {
				channels[i] = make(chan query.QueryData, 10)
			}

			originalFilePath := config.FilePath
			config.FilePath = tmpFile.Name()
			defer func() { config.FilePath = originalFilePath }()

			errChan := make(chan error, 1)
			go func() {
				errChan <- Stream(channels)
			}()

			receivedCount := make([]int, len(channels))
			var streamErr error

			for {
				select {
				case qd := <-channels[0]:
					receivedCount[0]++
					assert.NotEmpty(t, qd.Hostname, "Hostname should not be empty")
				case qd := <-channels[1]:
					receivedCount[1]++
					assert.NotEmpty(t, qd.Hostname, "Hostname should not be empty")
				default:
					select {
					case streamErr = <-errChan:
						goto checkResult
					default:
						continue
					}
				}
			}

		checkResult:
			if testCase.expectError {
				assert.Error(t, streamErr, "Expected an error but got nil")
				if testCase.errorMsg != "" {
					assert.Contains(t, streamErr.Error(), testCase.errorMsg,
						"Expected error message to contain %s, got %s", testCase.errorMsg, streamErr.Error())
				}
			} else {
				assert.NoError(t, streamErr, "Expected no error but got: %v", streamErr)
			}

			for _, ch := range channels {
				close(ch)
			}

			if !testCase.expectError {
				assert.Equal(t, testCase.expectedCount, receivedCount[0]+receivedCount[1],
					"Expected %d records, got %d", testCase.expectedCount, receivedCount[0]+receivedCount[1])
			}
		})
	}
}

func createTempCSVFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "test_csv_*.csv")
	assert.NoError(t, err, "Failed to create temporary file")

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err, "Failed to write to temporary file")

	_, err = tmpFile.Seek(0, 0)
	assert.NoError(t, err, "Failed to seek to beginning of file")

	return tmpFile
}
