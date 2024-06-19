package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"

	inputreader "github.com/molinama/timescale/src/input_reader"
	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/repository"
	"github.com/molinama/timescale/src/worker"
)

// Config struct to hold command line arguments
type Config struct {
	CSVFilePath   string
	NumberWorkers int
	DBConnString  string
}

var config Config

const (
	WORKERS = 10 // Default number of workers.
	TASKS   = 10 // Default number of tasks in the channel.
)

func init() {
	// Parse command line arguments
	flag.StringVar(&config.CSVFilePath, "csv", "query_params.csv", "The file path to the CSV file containing query parameters.")
	flag.IntVar(&config.NumberWorkers, "workers", WORKERS, "The number of workers for the pool.")
	flag.StringVar(&config.DBConnString, "db", "your_db_connection_string", "The database connection string")
}

func main() {
	flag.Parse()

	if config.CSVFilePath == "" {
		log.Fatal("CSV file path is empty")
	}

	// Initialize CSV reader
	reader, err := inputreader.NewCSVReader(config.CSVFilePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer reader.Close()

	// Connect to the database
	db, err := sql.Open("pgx", config.DBConnString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Create repository
	repository, err := repository.NewQueryParamsRepository(db)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Initialize and start the worker pool
	workerPool := worker.NewWorkerPool(TASKS, config.NumberWorkers)
	workerPool.Start()
	defer workerPool.Stop()

	// Channel to collect task results
	results := make(chan model.QueryTaskResult, TASKS)
	defer close(results)

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

	// Process tasks from the CSV reader
	processTasks(reader, workerPool, repository, results)

	// Calculate and print query statistics
	queryStats := model.QueryStats{}
	queryStats.CalculateStats(allResults)
	fmt.Print(queryStats)
}

// processTasks reads parameters from the reader, creates tasks, and adds them to the worker pool
func processTasks(reader inputreader.Reader, workerPool *worker.WorkerPool, repository *repository.QueryParamsRepository, results chan<- model.QueryTaskResult) {
	for {
		params, err := reader.Read()
		if err == io.EOF {
			break // Exit loop at end of file
		}
		if err != nil {
			log.Printf("Error reading CSV file: %v", err)
			continue // Skip to the next line on error
		}

		// Create a new query task and add it to the worker pool
		task := worker.NewQueryTask(repository, params, results)
		workerPool.Add(task)
	}
}
