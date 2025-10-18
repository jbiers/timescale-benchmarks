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

var (
	DatabaseURL string
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

func InitEnv() {
	DatabaseURL = os.Getenv("DB_URL")
	if DatabaseURL == "" {
		DatabaseURL = "postgres://postgres:password@localhost:5432/homework"
	}
}

//user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10 pool_max_conn_lifetime=1h30m

func InitLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	if Debug {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.Debug("Debug logging enabled")
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}
}
