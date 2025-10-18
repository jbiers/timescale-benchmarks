package csvreader

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jbiers/timescale-benchmark/config"
	"github.com/jbiers/timescale-benchmark/pkg/query"
)

// TODO: pass in a channel that will receive each line as a querydata type
func Stream(chs []chan query.QueryData) error {
	var IOreader io.Reader

	if config.FilePath != "" {
		file, err := os.Open(config.FilePath)
		if err != nil {
			return fmt.Errorf("error opening the query data file %s: %w", config.FilePath, err)
		}

		// TODO: what are the real performance differences of using bufio.Reader vs a regular os.File?
		IOreader = bufio.NewReader(file)
		defer file.Close()
	} else {
		// TODO: does the program work properly with piped input and with copying text into stdin after it's started?
		IOreader = bufio.NewReader(os.Stdin)
	}

	CSVreader := csv.NewReader(IOreader)
	CSVreader.FieldsPerRecord = 3
	CSVreader.ReuseRecord = true

	header, err := CSVreader.Read()
	if err != nil {
		return fmt.Errorf("error parsing first line of query data file %s: %w", config.FilePath, err)
	}

	err = validateHeaderLine(header)
	if err != nil {
		return fmt.Errorf("header validation failed for file %s: %w", config.FilePath, err)
	}

	for {
		record, err := CSVreader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV line from query data file %s: %w", config.FilePath, err)
		}

		queryData, err := parseQueryData(record)
		if err != nil {
			return fmt.Errorf("error parsing query data: %w", err)
		}

		idx := queryData.GetHash(len(chs))
		chs[idx] <- queryData
	}

	return nil
}

// TODO: can I assume that for all CSV files the first line will contain headers?
// this piece of code will invalidate any file without a header that matches the current setup
// which makes the program very specific and not adaptable to other queries
func validateHeaderLine(header []string) error {

	if len(header) != 3 {
		return fmt.Errorf("expected three entries in header line")
	}

	if header[0] != "hostname" {
		return fmt.Errorf("expected first entry in header line to match 'hostname'")
	}

	if header[1] != "start_time" {
		return fmt.Errorf("expected first entry in header line to match 'start_time'")
	}

	if header[2] != "end_time" {
		return fmt.Errorf("expected first entry in header line to match 'end_time'")
	}

	return nil
}

func parseQueryData(record []string) (query.QueryData, error) {
	if len(record) != 3 {
		return query.QueryData{}, fmt.Errorf("expected 3 entries in query data record, found %d", len(record))
	}

	timeStringFormat := "2006-01-02 15:04:05"
	startTime, err := time.Parse(timeStringFormat, record[1])
	if err != nil {
		return query.QueryData{}, fmt.Errorf("failed to parse 'start_time' as time.Time format: %w", err)
	}

	endTime, err := time.Parse(timeStringFormat, record[2])
	if err != nil {
		return query.QueryData{}, fmt.Errorf("failed to parse 'end_time' as time.Time format: %w", err)
	}

	return query.QueryData{
		Hostname:  record[0],
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}
