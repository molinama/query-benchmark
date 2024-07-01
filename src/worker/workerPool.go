package worker

import (
	"context"
	"sync"

	"github.com/molinama/timescale/src/logging"
	"github.com/molinama/timescale/src/model"
	"go.uber.org/zap"
)

type WorkerPool struct {
	numberWorkers     int
	numberTasks       int
	workerChannelsMap map[model.Worker]chan Task
	stop              bool
	WgTasks           sync.WaitGroup
	WgWorkers         sync.WaitGroup
	mu                sync.RWMutex
	ctx               context.Context
}

func NewWorkerPool(ctx context.Context, numberTasks int, numberWorkers int) *WorkerPool {
	return &WorkerPool{
		numberTasks:       numberTasks,
		numberWorkers:     numberWorkers,
		workerChannelsMap: make(map[model.Worker]chan Task, numberWorkers),
		WgTasks:           sync.WaitGroup{},
		WgWorkers:         sync.WaitGroup{},
		ctx:               ctx,
	}
}

func (wp *WorkerPool) NumberWorkers() int {
	return wp.numberWorkers
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numberWorkers; i++ {
		worker := model.Worker(i)
		wp.WgWorkers.Add(1)
		workerChannel := make(chan Task, wp.numberTasks)

		wp.mu.Lock()
		wp.workerChannelsMap[worker] = workerChannel
		wp.mu.Unlock()

		go wp.run(worker)
	}
}

func (wp *WorkerPool) run(worker model.Worker) {
	logging.Log.Debug("Start Worker", zap.Int("workerId", int(worker)))

	for {
		wp.mu.RLock()
		workerChannel := wp.workerChannelsMap[worker]
		wp.mu.RUnlock()

		select {
		case task, ok := <-workerChannel:
			if !ok {
				logging.Log.Debug("Channel is closed for Worker", zap.Int("workerId", int(worker)))
				wp.WgWorkers.Done()
				return
			}
			logging.Log.Debug("Worker running Hostname", zap.Int("workerId", int(worker)), zap.String("hostname", task.Hostname()))
			task.Execute(worker)

		case <-wp.ctx.Done():
			logging.Log.Debug("Quitting Worker", zap.Int("workerId", int(worker)))
			wp.WgWorkers.Done()
			close(wp.workerChannelsMap[worker])
			return
		}

	}
}

func (wp *WorkerPool) Add(worker model.Worker, task Task) {
	if !wp.stop {
		wp.WgTasks.Add(1)

		wp.mu.RLock()
		workerChannel, exists := wp.workerChannelsMap[worker]
		wp.mu.RUnlock()

		if !exists {
			logging.Log.Debug("Nil channel Worker", zap.Int("workerId", int(worker)))
			return
		}

		workerChannel <- task
	}
}

func (wp *WorkerPool) Stop(cancel context.CancelFunc) {
	// Wait all tasks to be completed.
	wp.WgTasks.Wait()
	logging.Log.Info("All Tasks Completed")
	cancel()
	logging.Log.Info("Stopping Worker Pool")
	wp.stop = true

	// Wait all workers quit.
	wp.WgWorkers.Wait()
	logging.Log.Info("All Workers Ended")
}
