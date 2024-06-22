package model

type QueryTaskErr struct {
	QueryTaskResult
	RawQuery string
	Err      error
}
