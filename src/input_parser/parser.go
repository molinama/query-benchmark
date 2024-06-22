package inputparser

import "github.com/molinama/timescale/src/model"

type Reader interface {
	Parse() (*model.QueryParams, error)
	Close() error
}
