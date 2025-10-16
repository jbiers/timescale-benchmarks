package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jbiers/timescale-benchmark/pkg/csvreader"
	"github.com/jbiers/timescale-benchmark/pkg/database"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	"github.com/jbiers/timescale-benchmark/pkg/workerpool"
	"github.com/sirupsen/logrus"
)

type benchmarkConfig struct {
	file    *string
	workers *int
}

var config benchmarkConfig

func buildQueryDataChannels() []chan query.QueryData {
	queryDataChannels := make([]chan query.QueryData, *config.workers)

	for ch := range queryDataChannels {
		queryDataChannels[ch] = make(chan query.QueryData)
	}

	return queryDataChannels
}

// TODO: clean this ugly thing
func getMetrics(r []time.Duration) (int, time.Duration, time.Duration, time.Duration, time.Duration) {
	var total time.Duration
	var longest time.Duration
	var shortest time.Duration
	var average time.Duration

	for i, t := range r {
		total += t

		if t > longest {
			longest = t
		}

		if i == 0 {
			shortest = t
		} else {
			if t < shortest {
				shortest = t
			}
		}
	}

	average = total / time.Duration(len(r))

	return len(r), total, longest, shortest, average
}

func init() {
	config.file = flag.String("file", "", "Path to CSV formatted file containing the query parameters to be run in the benchmark. Defaults to nil, waiting for STDIN.")
	config.workers = flag.Int("workers", 1, "Number of concurrent workers for querying the database. Defaults to 1.")

	flag.Parse()
}

// TODO: should start thinking about graceful shutdown
func main() {
	dataChannels := buildQueryDataChannels()

	go func() {
		err := csvreader.Stream(*config.file, dataChannels)
		if err != nil {
			logrus.Fatalf("failed to stream from CSV file: %v", err)
		}

		for _, w := range dataChannels {
			close(w)
		}
	}()

	databasePool, err := database.InitDB()
	if err != nil {
		logrus.Fatalf("database initialization failed: %v", err)
	}
	defer databasePool.Close()

	workerPool := workerpool.NewWorkerPool(dataChannels, *config.workers, databasePool)
	workerPool.Dispatch()

	num, totalTime, longest, shortest, average := getMetrics(workerPool.Results)
	fmt.Println("Number of queries:", num)
	fmt.Println("Total time:", totalTime)
	fmt.Println("Average time:", average)
	fmt.Println("Longest time:", longest)
	fmt.Println("Shortest time:", shortest)
	fmt.Println("Median time:")
}
