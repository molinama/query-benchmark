package worker

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/molinama/timescale/src/model"
)

type MockTask struct {
	ID        int
	ResultsCh chan<- time.Duration
	wg        *sync.WaitGroup
}

func (mt *MockTask) Execute(worker model.Worker) {
	start := time.Now()
	// Simulate variable work time
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	duration := time.Since(start)
	log.Printf("Worker %d executed task %d in %v\n", worker, mt.ID, duration)
	mt.ResultsCh <- duration
	mt.wg.Done()
}
func (mt *MockTask) Hostname() string {
	return "host"
}

func TestWorkerPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	numberWorkers := 5
	numberTasks := 10
	wp := NewWorkerPool(ctx, numberTasks, numberWorkers)
	resultsCh := make(chan time.Duration, numberTasks)

	// Start the worker pool
	wp.Start()

	// Create and add tasks
	for i := 0; i < numberTasks; i++ {
		task := &MockTask{
			ID:        i,
			ResultsCh: resultsCh,
			wg:        &wp.WgTasks,
		}
		workerId := rand.Intn(numberWorkers-1) + 1
		go wp.Add(model.Worker(workerId), task)
	}

	time.Sleep(100 * time.Millisecond)

	// Collect results
	var results []time.Duration
	var resultMu sync.Mutex
	go func() {
		for res := range resultsCh {
			resultMu.Lock()
			results = append(results, res)
			resultMu.Unlock()
		}
	}()

	// Stop the worker pool
	wp.Stop(cancel)

	// Ensure no more results are sent to the channel
	close(resultsCh)

	// Ensure all workers stopped
	select {
	case <-wp.ctx.Done():
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
	ctx, cancel := context.WithCancel(context.Background())
	numberWorkers := 10
	numberTasks := 1000

	// Initialize the worker pool
	wp := NewWorkerPool(ctx, numberTasks, numberWorkers)
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
			workerId := rand.Intn(numberWorkers-1) + 1
			go wp.Add(model.Worker(workerId), task)
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
	wp.Stop(cancel)

	// Ensure no more results are sent to the channel
	close(resultsCh)
}
