package session

import (
	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/worker"
)

type Session interface {
	GetWorker(task worker.Task) model.Worker
}
