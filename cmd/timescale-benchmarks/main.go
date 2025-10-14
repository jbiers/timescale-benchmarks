package main

import (
	"flag"

	"github.com/jbiers/timescale-benchmark/pkg/csvreader"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/jbiers/timescale-benchmark/pkg/workerpool"
	"github.com/sirupsen/logrus"
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
	// TODO: how much buffering should the channel really have?
	queryDataChannel := make(chan query.QueryData, 100)

	go func() {
		err := csvreader.Stream(*config.file, queryDataChannel)
		if err != nil {
			logrus.Fatal(err)
		}

		close(queryDataChannel)
	}()

	wp := workerpool.NewWorkerPool(&queryDataChannel, *config.workers)
	wp.Dispatch()

}
