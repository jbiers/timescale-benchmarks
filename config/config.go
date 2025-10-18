package config

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	FilePath string
	Workers  int
	Debug    bool
)

var Logger *logrus.Logger

const (
	fileFlag    = "file"
	fileDefault = ""
	fileMessage = "Path to CSV formatted file containing the query parameters to be run in the benchmark. Defaults to nil, waiting for STDIN."

	workersFlag    = "workers"
	workersDefault = 1
	workersMessage = "Number of concurrent workers for querying the database. Defaults to 1."

	debugFlag    = "debug"
	debugDefault = false
	debugMessage = "Debug mode will have the program log out the processing time for each individual query. Defaults to false."
)

func InitFlags() {
	flag.StringVar(&FilePath, fileFlag, fileDefault, fileMessage)
	flag.IntVar(&Workers, workersFlag, workersDefault, workersMessage)
	flag.BoolVar(&Debug, debugFlag, debugDefault, debugMessage)

	flag.Parse()
}

func InitLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	if Debug {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		Logger.Debug("Debug logging enabled")
	} else {
		Logger.SetLevel(logrus.InfoLevel)
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		})
	}
}
