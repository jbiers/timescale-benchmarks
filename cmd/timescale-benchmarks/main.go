package main

import (
	"flag"
	"log"

	"github.com/jbiers/timescale-benchmark/pkg/csvreader"
	"github.com/jbiers/timescale-benchmark/pkg/database"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/jbiers/timescale-benchmark/pkg/workerpool"
)

type benchmarkConfig struct {
	file              *string
	workers           *int
	queryDataChannels []chan query.QueryData
}

var config benchmarkConfig

func (cfg *benchmarkConfig) buildQueryDataChannels() {
	cfg.queryDataChannels = make([]chan query.QueryData, *config.workers)

	for ch := range cfg.queryDataChannels {
		cfg.queryDataChannels[ch] = make(chan query.QueryData)
	}
}

func init() {
	config.file = flag.String("file", "", "Path to CSV formatted file containing the query parameters to be run in the benchmark. Defaults to nil, waiting for STDIN.")
	config.workers = flag.Int("workers", 1, "Number of concurrent workers for querying the database. Defaults to 1.")

	flag.Parse()
}

// TODO: should start thinking about graceful shutdown
func main() {
	config.buildQueryDataChannels()

	go func() {
		err := csvreader.Stream(*config.file, config.queryDataChannels)
		if err != nil {
			log.Fatalf("failed to stream from CSV file: %v", err)
		}

		for _, w := range config.queryDataChannels {
			close(w)
		}
	}()

	databasePool, err := database.InitDB()
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}
	defer databasePool.Close()

	workerPool := workerpool.NewWorkerPool(config.queryDataChannels, *config.workers, databasePool)
	workerPool.Dispatch()
}
