package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/molinama/timescale/src/model"
)

type QueryParamsRepository struct {
	db *sql.DB
}

func NewQueryParamsRepository(db *sql.DB) (*QueryParamsRepository, error) {
	if db == nil {
		return nil, errors.New("cannot create repository")
	}
	return &QueryParamsRepository{
		db: db,
	}, nil
}

func (repository *QueryParamsRepository) RawQuery(params *model.QueryParams) (time.Duration, error) {
	start := time.Now()
	rows, err := repository.db.Query(params.RawQuery())
	if err != nil {
		return time.Since(start), err
	}
	rows.Close()
	return time.Since(start), nil
}
