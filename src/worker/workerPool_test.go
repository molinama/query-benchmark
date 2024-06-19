package worker

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type MockTask struct {
	ID        int
	ResultsCh chan<- time.Duration
}

func (mt *MockTask) Execute(worker int) {
	start := time.Now()
	// Simulate variable work time
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	duration := time.Since(start)
	log.Printf("Worker %d executed task %d in %v\n", worker, mt.ID, duration)
	mt.ResultsCh <- duration
}

func TestWorkerPool(t *testing.T) {
	numberWorkers := 5
	numberTasks := 10
	wp := NewWorkerPool(numberTasks, numberWorkers)
	resultsCh := make(chan time.Duration, numberTasks)
	var wg sync.WaitGroup

	// Create and add tasks
	for i := 0; i < numberTasks; i++ {
		task := &MockTask{
			ID:        i,
			ResultsCh: resultsCh,
		}
		wp.Add(task)
		wg.Add(1)
	}

	// Start the worker pool
	wp.Start()

	// Collect results
	var results []time.Duration
	var resultMu sync.Mutex
	go func() {
		for res := range resultsCh {
			resultMu.Lock()
			results = append(results, res)
			resultMu.Unlock()
			wg.Done()
		}
	}()

	// Wait for all tasks to complete
	wg.Wait()

	// Stop the worker pool
	wp.Stop()

	// Ensure no more results are sent to the channel
	close(resultsCh)

	// Ensure all workers stopped
	select {
	case <-wp.quit:
		t.Log("Worker pool stopped successfully")
	default:
		t.Error("Worker pool did not stop as expected")
	}

	// Validate results
	if len(results) != numberTasks {
		t.Errorf("Expected %d results, but got %d", numberTasks, len(results))
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	numberWorkers := 10
	numberTasks := 1000

	// Initialize the worker pool
	wp := NewWorkerPool(numberTasks, numberWorkers)
	resultsCh := make(chan time.Duration, numberTasks)
	var wg sync.WaitGroup

	// Start the worker pool
	wp.Start()

	b.ResetTimer() // Reset the timer to exclude setup time

	for n := 0; n < b.N; n++ {
		// Create and add tasks
		for i := 0; i < numberTasks; i++ {
			task := &MockTask{
				ID:        i,
				ResultsCh: resultsCh,
			}
			wg.Add(1)
			go wp.Add(task)
		}

		// Collect results
		var results []time.Duration
		var resultMux sync.Mutex
		go func() {
			for res := range resultsCh {
				resultMux.Lock()
				results = append(results, res)
				resultMux.Unlock()
				wg.Done()
			}
		}()

		// Wait for all tasks to complete
		wg.Wait()
	}

	// Stop the worker pool
	wp.Stop()

	// Ensure no more results are sent to the channel
	close(resultsCh)
}
