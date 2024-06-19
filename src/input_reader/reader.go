package inputreader

import "github.com/molinama/timescale/src/repository"

type Reader interface {
	Read() (*repository.QueryParams, error)
	Close() error
}
