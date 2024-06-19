package repository

import (
	"database/sql"
	"errors"
	"time"
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

func (repository *QueryParamsRepository) RawQuery(params *QueryParams) (time.Duration, error) {
	start := time.Now()
	_, err := repository.db.Query(params.RawQuery())
	if err != nil {
		return -1, err
	}
	return time.Since(start), nil
}
