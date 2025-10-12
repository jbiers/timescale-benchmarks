package main

import (
	"flag"

	"github.com/jbiers/timescale-benchmark/pkg/csvreader"
)

type benchmarkConfig struct {
	file    *string
	workers *int
}

var config benchmarkConfig

func init() {
	config.file = flag.String("file", "", "Path to CSV formatted file containing the query parameters to be run in the benchmark. Defaults to nil, waiting for STDIN.")
	config.workers = flag.Int("workers", 1, "Number of concurrent workers for querying the database. Defaults to 1.")

	flag.Parse()
}

func main() {
	// Opens the CSV file and parses it in a separate Goroutine
	// As you process it, each line will be converted into a QueryData struct and sent into a channel
	csvreader.Stream(*config.file)
}
