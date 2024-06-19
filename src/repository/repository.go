package repository

import "time"

type Repository interface {
	RawQuery(params *QueryParams) (time.Duration, error)
}
