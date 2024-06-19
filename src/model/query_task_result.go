package model

import "time"

type QueryTaskResult struct {
	Worker   int
	Hostname string
	RawQuery string
	time.Duration
	Err error
}
