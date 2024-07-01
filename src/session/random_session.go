package session

import (
	"math/rand"
	"sync"

	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/worker"
)

type randomSession struct {
	wp           *worker.WorkerPool
	workerByHost map[string]model.Worker
	sessionMux   sync.Mutex
}

func NewRandomSession(wp *worker.WorkerPool) Session {
	return &randomSession{
		wp:           wp,
		workerByHost: make(map[string]model.Worker),
	}
}

func (rs *randomSession) GetWorker(task worker.Task) model.Worker {
	rs.sessionMux.Lock()
	defer rs.sessionMux.Unlock()

	if _, ok := rs.workerByHost[task.Hostname()]; !ok {
		var workerId int
		if rs.wp.NumberWorkers() <= 1 {
			workerId = 1
		} else {
			workerId = rand.Intn(rs.wp.NumberWorkers()-1) + 1
		}
		rs.workerByHost[task.Hostname()] = model.Worker(workerId)
	}
	return rs.workerByHost[task.Hostname()]
}
