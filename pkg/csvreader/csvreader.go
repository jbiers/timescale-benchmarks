package csvreader

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// TODO: pass in a channel that will receive each line as a querydata type
func Stream(filePath string) error {
	var IOreader io.Reader

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening the query data file %s: %w", filePath, err)
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
	//CSVreader.ReuseRecord = true

	header, err := CSVreader.Read()
	if err != nil {
		return fmt.Errorf("Error parsing first line of query data file %s: %w", filePath, err)
	}

	err = validateHeaderLine(header)
	if err != nil {
		return fmt.Errorf("Header validation failed for file %s: %w", filePath, err)
	}

	for {
		record, err := CSVreader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error parsing CSV line from query data file %s: %w", filePath, err)
		}

		fmt.Println(record)
	}

	return nil
}

// TODO: can I assume that for all CSV files the first line will contain headers?
// this piece of code will invalidate any file without a header that matches the current setup
// which makes the program very specific and not adaptable to other queries
func validateHeaderLine(header []string) error {

	if len(header) != 3 {
		return fmt.Errorf("Expected three entries in header line")
	}

	if header[0] != "hostname" {
		return fmt.Errorf("Expected first entry in header line to match 'hostname'")
	}

	if header[1] != "start_time" {
		return fmt.Errorf("Expected first entry in header line to match 'start_time'")
	}

	if header[2] != "end_time" {
		return fmt.Errorf("Expected first entry in header line to match 'end_time'")
	}

	return nil
}
