package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	inputparser "github.com/molinama/timescale/src/input_parser"
	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/repository"
	"github.com/molinama/timescale/src/worker"
)

// Config struct to hold command line arguments
type Config struct {
	csvFilePath   string
	numberWorkers int
	dbConnString  string
	db            *sql.DB
}

func (c *Config) loadEnv() {
	err := godotenv.Load("./timescaledb/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("TIMESCALE_USER")
	pwd := os.Getenv("TIMESCALE_PASSWORD")
	host := os.Getenv("TIMESCALE_HOST")
	db := os.Getenv("TIMESCALES_DB")
	port := os.Getenv("TIMESCALE_PORT")

	c.dbConnString = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		pwd,
		host,
		port,
		db)
}

var config Config

const (
	WORKERS = 10 // Default number of workers.
	TASKS   = 10 // Default number of tasks in the channel.
)

func init() {
	// Parse command line arguments
	flag.StringVar(&config.csvFilePath, "csv", "./query_params.csv", "The file path to the CSV file containing query parameters.")
	flag.IntVar(&config.numberWorkers, "workers", WORKERS, "The number of workers for the pool.")
}

func main() {
	// Get flags
	flag.Usage = usage
	flag.Parse()

	// Connect to the database
	config.loadEnv()
	db, err := sql.Open("pgx", config.dbConnString)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Run main application
	config.db = db
	err = run(config)
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
	if config.csvFilePath == "" {
		return fmt.Errorf("CSV file path is empty")
	}

	// Initialize CSV reader
	reader, err := inputparser.NewCSVReader(config.csvFilePath)
	if err != nil {
		return fmt.Errorf("CSV file path is empty")
	}
	defer reader.Close()

	// Create repository
	repository, err := repository.NewQueryParamsRepository(config.db)
	if err != nil {
		return fmt.Errorf("failed to create repository: %v", err)
	}

	// Initialize and start the worker pool
	workerPool := worker.NewWorkerPool(TASKS, config.numberWorkers)
	workerPool.Start()

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
	processTasks(reader, workerConfig)

	// Wait all tasks to be completed.
	workerPool.WgTasks.Wait()
	fmt.Println("\nAll Tasks Completed")
	workerPool.Stop()

	// Wait all workers quit.
	workerPool.WgWorkers.Wait()
	fmt.Println("\nAll Workers Ended")

	// Calculate and print query statistics
	queryStats := model.Stats{}
	queryStats.CalculateStats(allResults, allErrs)
	fmt.Print(queryStats)

	return nil
}

// processTasks reads parameters from the reader, creates tasks, and adds them to the worker pool
func processTasks(reader inputparser.Reader, taskConfig worker.QueryTaskConfig) {
	for {
		params, err := reader.Parse()
		if err == io.EOF {
			break // Exit loop at end of file
		}
		if err != nil {
			log.Printf("Error reading CSV file: %v", err)
			continue // Skip to the next line on error
		}

		// Create a new query task and add it to the worker pool
		taskConfig.Params = params
		task := worker.NewQueryTask(taskConfig)

		taskConfig.WorkerPool.Add(task)
	}
}
