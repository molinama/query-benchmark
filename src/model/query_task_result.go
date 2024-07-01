package model

import "time"

type QueryTaskResult struct {
	Worker   Worker
	Hostname string
	time.Duration
}
