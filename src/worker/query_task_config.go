package worker

import (
	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/repository"
)

type QueryTaskConfig struct {
	Repository repository.Repository
	Params     *model.QueryParams
	Results    chan<- model.QueryTaskResult
	Errs       chan<- model.QueryTaskErr
	WorkerPool *WorkerPool
}
