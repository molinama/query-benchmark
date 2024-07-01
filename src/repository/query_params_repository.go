package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/molinama/timescale/src/model"
)

type QueryParamsRepository struct {
	db *sql.DB
}

func NewQueryParamsRepository() (*QueryParamsRepository, error) {
	// Connect to the database
	dbConfig, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot create repository: %w", err)
	}

	db, err := sql.Open("pgx", dbConfig.dbConnString)
	if err != nil {
		return nil, errors.New("cannot create repository")
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

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
