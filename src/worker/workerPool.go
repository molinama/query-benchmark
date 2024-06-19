package worker

import (
	"log"
)

type WorkerPool struct {
	numberWorkers int
	numberTasks   int
	tasks         chan Task
	quit          chan struct{}
	stop          bool
}

func NewWorkerPool(numberTasks int, numberWorkers int) *WorkerPool {
	return &WorkerPool{
		numberTasks:   numberTasks,
		numberWorkers: numberWorkers,
		tasks:         make(chan Task, numberTasks),
		quit:          make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numberWorkers; i++ {
		go wp.run(i)
	}
}

func (wp *WorkerPool) run(id int) {
	log.Printf("Start Worker: %d\n", id)

	for {
		select {
		case task, ok := <-wp.tasks:
			if !ok {
				log.Printf("Channel is close, ending Worker: %d", id)
			}
			task.Execute(id)

		case <-wp.quit:
			log.Printf("Quitting Worker Id: #%d\n", id)

			return
		}

	}
}

func (wp *WorkerPool) Add(task Task) {
	if !wp.stop {
		wp.tasks <- task
	}
}

func (wp *WorkerPool) Stop() {
	log.Printf("Stopping Worker Pool")
	close(wp.quit)
	wp.stop = true
}
