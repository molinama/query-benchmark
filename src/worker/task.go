package worker

import "github.com/molinama/timescale/src/model"

type Task interface {
	Execute(worker model.Worker)
	Hostname() string
}
