package repository

import (
	"time"

	"github.com/molinama/timescale/src/model"
)

type Repository interface {
	RawQuery(params *model.QueryParams) (time.Duration, error)
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}
