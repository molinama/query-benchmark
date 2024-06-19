package worker

import (
	"log"
	"time"

	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/repository"
)

type QueryTask struct {
	repository repository.Repository
	params     *repository.QueryParams
	results    chan<- model.QueryTaskResult
}

func NewQueryTask(repository repository.Repository, params *repository.QueryParams, results chan<- model.QueryTaskResult) *QueryTask {
	return &QueryTask{
		repository: repository,
		params:     params,
		results:    results,
	}
}

func (t *QueryTask) Execute(worker int) {
	start := time.Now()
	_, err := t.repository.RawQuery(t.params)
	if err != nil {
		log.Printf("Error running query: %v", err)
	}
	result := model.QueryTaskResult{
		Worker:   worker,
		Hostname: t.params.Hostname,
		RawQuery: t.params.RawQuery(),
		Duration: time.Since(start),
		Err:      err,
	}

	t.results <- result
}
