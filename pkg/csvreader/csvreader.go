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

		IOreader = bufio.NewReader(os.Stdin)
	}

	CSVreader := csv.NewReader(IOreader)
	CSVreader.FieldsPerRecord = 3
	//CSVreader.ReuseRecord = true

	// TODO: use reader.Read to parse one line at a time without loading the entire file into memory.
	records, err := CSVreader.ReadAll()
	if err != nil {
		return fmt.Errorf("Error parsing the query data file %s as CSV: %w", filePath, err)
	}

	for _, record := range records {
		fmt.Println(record)
	}

	return nil
}
