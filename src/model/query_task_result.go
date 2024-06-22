package model

import "time"

type QueryTaskResult struct {
	Worker   int
	Hostname string
	time.Duration
}
