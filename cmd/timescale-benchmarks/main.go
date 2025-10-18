package main

import (
	"github.com/jbiers/timescale-benchmark/config"
	"github.com/jbiers/timescale-benchmark/pkg/csvreader"
	"github.com/jbiers/timescale-benchmark/pkg/database"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	wp "github.com/jbiers/timescale-benchmark/pkg/workerpool"
	"github.com/sirupsen/logrus"
)

func init() {
	config.ParseConfig()
}

// TODO: should start thinking about graceful shutdown
func main() {
	dataChannels := buildQueryDataChannels()

	go func() {
		err := csvreader.Stream(dataChannels)
		if err != nil {
			logrus.Fatalf("failed to stream from CSV file: %v", err)
		}

		for _, w := range dataChannels {
			close(w)
		}
	}()

	dbPool, err := database.InitDB()
	if err != nil {
		logrus.Fatalf("database initialization failed: %v", err)
	}
	defer dbPool.Close()

	workerPool := wp.NewWorkerPool(dataChannels, dbPool)
	wpMetrics := workerPool.Dispatch()
	wpMetrics.ReportWorkerPoolMetrics()
}

func buildQueryDataChannels() []chan query.QueryData {
	queryDataChannels := make([]chan query.QueryData, config.Workers)

	for ch := range queryDataChannels {
		queryDataChannels[ch] = make(chan query.QueryData)
	}

	return queryDataChannels
}
