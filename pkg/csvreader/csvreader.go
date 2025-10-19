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

const timeStringFormat = "2006-01-02 15:04:05"

func Stream(chs []chan query.QueryData) error {
	var IOreader io.Reader

	if config.FilePath != "" {
		file, err := os.Open(config.FilePath)
		if err != nil {
			return fmt.Errorf("error opening the query data file %s: %w", config.FilePath, err)
		}

		IOreader = bufio.NewReader(file)
		defer file.Close()
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no input provided (expected --file flag or piped data)")
		}

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

func validateHeaderLine(header []string) error {

	if len(header) != 3 {
		return fmt.Errorf("expected three entries in header line, found %d", len(header))
	}

	if header[0] != "hostname" {
		return fmt.Errorf("expected first entry in header line to match 'hostname', found '%s'", header[0])
	}

	if header[1] != "start_time" {
		return fmt.Errorf("expected second entry in header line to match 'start_time', found '%s'", header[1])
	}

	if header[2] != "end_time" {
		return fmt.Errorf("expected third entry in header line to match 'end_time', found '%s'", header[2])
	}

	return nil
}

func parseQueryData(record []string) (query.QueryData, error) {
	if len(record) != 3 {
		return query.QueryData{}, fmt.Errorf("expected 3 entries in query data record, found %d", len(record))
	}

	startTime, err := time.Parse(timeStringFormat, record[1])
	if err != nil {
		return query.QueryData{}, fmt.Errorf("failed to parse 'start_time' '%s' as time.Time format: %w", record[1], err)
	}

	endTime, err := time.Parse(timeStringFormat, record[2])
	if err != nil {
		return query.QueryData{}, fmt.Errorf("failed to parse 'end_time' '%s' as time.Time format: %w", record[2], err)
	}

	return query.QueryData{
		Hostname:  record[0],
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}
