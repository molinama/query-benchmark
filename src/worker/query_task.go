package worker

import (
	"log"
	"sync"

	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/repository"
)

type QueryTask struct {
	repository repository.Repository
	params     *model.QueryParams
	results    chan<- model.QueryTaskResult
	errs       chan<- model.QueryTaskErr
	wg         *sync.WaitGroup
}

func NewQueryTask(config QueryTaskConfig) *QueryTask {
	return &QueryTask{
		repository: config.Repository,
		params:     config.Params,
		results:    config.Results,
		errs:       config.Errs,
		wg:         &config.WorkerPool.WgTasks,
	}
}

func (t *QueryTask) Execute(worker int) {
	defer t.wg.Done()

	duration, err := t.repository.RawQuery(t.params)
	//log.Printf("Query executed: %v", t.params.RawQuery())
	result := model.QueryTaskResult{
		Worker:   worker,
		Hostname: t.params.Hostname,
		Duration: duration,
	}

	if err != nil {
		queryTaskErr := model.QueryTaskErr{
			QueryTaskResult: result,
			RawQuery:        t.params.RawQuery(),
			Err:             err,
		}
		t.errs <- queryTaskErr
		log.Printf("Error running query: %v", err)
	} else {
		t.results <- result
	}

}
