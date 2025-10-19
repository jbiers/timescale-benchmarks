package main

import (
	"context"

	"github.com/jbiers/timescale-benchmark/config"
	"github.com/jbiers/timescale-benchmark/pkg/csvreader"
	"github.com/jbiers/timescale-benchmark/pkg/database"
	"github.com/jbiers/timescale-benchmark/pkg/query"
	wp "github.com/jbiers/timescale-benchmark/pkg/workerpool"
)

func init() {
	config.InitFlags()
	config.InitEnv()
	config.InitLogger()

	config.Logger.Infof("Program initialized with: Workers: %d, Debug: %t", config.Workers, config.Debug)
}

func main() {
	dataChannels := buildQueryDataChannels()

	go func() {
		err := csvreader.Stream(dataChannels)
		if err != nil {
			config.Logger.Fatalf("Failed to stream from CSV file: %v", err)
		}

		for _, w := range dataChannels {
			close(w)
		}
	}()

	ctx := context.Background()

	dbPool, err := database.InitDB(ctx)
	if err != nil {
		config.Logger.Fatalf("Database initialization failed: %v", err)
	}
	defer dbPool.Close()

	workerPool := wp.NewWorkerPool(dataChannels, dbPool)
	wpMetrics := workerPool.Dispatch(ctx)
	wpMetrics.ReportWorkerPoolMetrics()
}

func buildQueryDataChannels() []chan query.QueryData {
	queryDataChannels := make([]chan query.QueryData, config.Workers)

	for ch := range queryDataChannels {
		queryDataChannels[ch] = make(chan query.QueryData)
	}

	return queryDataChannels
}
