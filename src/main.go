package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	inputparser "github.com/molinama/timescale/src/input_parser"
	"github.com/molinama/timescale/src/logging"
	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/repository"
	"github.com/molinama/timescale/src/session"
	"github.com/molinama/timescale/src/worker"
)

// Config struct to hold command line arguments
type Config struct {
	csvFilePath   string
	numberWorkers int
	dbConnString  string
	db            *sql.DB
}

var config Config

const (
	WORKERS = 10 // Default number of workers.
	TASKS   = 10 // Default number of tasks in the channel.
)

func init() {
	// Parse command line arguments
	flag.StringVar(&config.csvFilePath, "csv", "./query_params.csv", "The file path to the CSV file containing query parameters.")
	flag.IntVar(&config.numberWorkers, "workers", WORKERS, "The number of workers for the pool. Must be >= 1")
}

func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	// Get flags
	flag.Usage = usage
	flag.Parse()

	if config.numberWorkers <= 0 || config.csvFilePath == "" {
		usage()
		log.Fatal("Error with application parameters")
	}

	// Run main application
	err := run(config)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func usage() {
	msg := fmt.Sprintf(`usage: %s [OPTIONS]
	%s is a simple tool to do query benchmark
	`, "query-benchmark", "query-benchmark")
	fmt.Println(msg)
	flag.PrintDefaults()
}

func run(config Config) error {
	// Initialize Logger
	initLogging()

	// Initialize CSV reader
	reader := initCsvReader(config)
	defer reader.Close()

	// Create repository
	repository := initRepository()

	// Initialize and start the worker pool
	context, cancel := context.WithCancel(context.Background())
	workerPool := startWorkerPool(context, TASKS, config.numberWorkers)

	// Channel to collect task results
	results := make(chan model.QueryTaskResult, TASKS)
	defer close(results)
	// Channel to collect task errors
	errs := make(chan model.QueryTaskErr, TASKS)
	defer close(errs)

	// Goroutine to collect results from the results channel
	var allResults []model.QueryTaskResult
	var resultsMux sync.Mutex
	go func() {
		for result := range results {
			resultsMux.Lock()
			allResults = append(allResults, result)
			resultsMux.Unlock()
		}
	}()
	// Goroutine to collect errors from the errs channel
	var allErrs []model.QueryTaskErr
	var errsMux sync.Mutex
	go func() {
		for err := range errs {
			errsMux.Lock()
			allErrs = append(allErrs, err)
			errsMux.Unlock()
		}
	}()

	// Process tasks from the CSV reader
	workerConfig := worker.QueryTaskConfig{
		WorkerPool: workerPool,
		Repository: repository,
		Results:    results,
		Errs:       errs,
	}
	// Create a session for the Worker Pool.
	session := session.NewRandomSession(workerPool)
	processTasks(session, reader, workerConfig)

	// Stop WorkerPool.
	workerPool.Stop(cancel)

	// Calculate and print query statistics
	queryStats := model.Stats{}

	resultsMux.Lock()
	errsMux.Lock()
	queryStats.CalculateStats(allResults, allErrs)
	resultsMux.Unlock()
	errsMux.Unlock()

	fmt.Print(queryStats)

	return nil
}

func initLogging() {
	err := logging.InitGlobalLogger()
	if err != nil {
		logging.SugaredLog.Errorf("Logging setup failed: %s", err.Error())
		os.Exit(501)
	}
}

func initCsvReader(config Config) inputparser.Reader {
	reader, err := inputparser.NewCSVReader(config.csvFilePath)
	if err != nil {
		logging.SugaredLog.Errorf("CSV file path is empty")
		os.Exit(501)
	}
	return reader
}

func initRepository() repository.Repository {
	repository, err := repository.NewQueryParamsRepository()
	if err != nil {
		logging.SugaredLog.Errorf("failed to create repository: %v", err)
		os.Exit(501)
	}
	return repository
}

func startWorkerPool(ctx context.Context, tasks int, workers int) *worker.WorkerPool {
	workerPool := worker.NewWorkerPool(ctx, tasks, workers)
	workerPool.Start()

	return workerPool
}

// processTasks reads parameters from the reader, creates tasks, and adds them to the worker pool
func processTasks(session session.Session, reader inputparser.Reader, taskConfig worker.QueryTaskConfig) {
	for {
		params, err := reader.Parse()
		if err == io.EOF {
			break // Exit loop at end of file
		}
		if err != nil {
			logging.SugaredLog.Errorf("Error reading CSV file: %v", err)
			continue // Skip to the next line on error
		}

		// Create a new query task and add it to the worker pool
		taskConfig.Params = params
		task := worker.NewQueryTask(taskConfig)
		worker := session.GetWorker(task)

		taskConfig.WorkerPool.Add(worker, task)
	}
}
