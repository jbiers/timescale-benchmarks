package config

import "flag"

var (
	FilePath string
	Workers  int
	Debug    bool
)

const (
	fileFlag    = "file"
	fileDefault = ""
	fileMessage = "Path to CSV formatted file containing the query parameters to be run in the benchmark. Defaults to nil, waiting for STDIN."

	workersFlag    = "workers"
	workersDefault = 1
	workersMessage = "Number of concurrent workers for querying the database. Defaults to 1."

	debugFlag    = "debug"
	debugDefault = false
	debugMessage = "Debug mode will have the program log out the results and processing time for each individual query. Defaults to false."
)

func ParseConfig() {
	flag.StringVar(&FilePath, fileFlag, fileDefault, fileMessage)
	flag.IntVar(&Workers, workersFlag, workersDefault, workersMessage)
	flag.BoolVar(&Debug, debugFlag, debugDefault, debugMessage)

	flag.Parse()
}
