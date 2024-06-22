package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type QueryParams struct {
	Hostname  string
	StartTime string
	EndTime   string
}

func NewQueryParams(data []string) (*QueryParams, error) {
	for i := range data {
		data[i] = strings.TrimSpace(data[i])
	}

	if err := validate(data); err != nil {
		return nil, err
	}
	return &QueryParams{
		Hostname:  data[0],
		StartTime: data[1],
		EndTime:   data[2],
	}, nil
}

func validate(data []string) error {
	if len(data) != 3 {
		return errors.New("invalid format: expected exactly 3 elements")
	}

	hostname := data[0]
	startTime := data[1]
	endTime := data[2]

	if hostname == "" {
		return errors.New("invalid format: hostname cannot be empty")
	}

	const timeLayout = "2006-01-02 15:04:05"
	if _, err := time.Parse(timeLayout, startTime); err != nil {
		return fmt.Errorf("invalid format: startTime: %s is not in the correct format (expected %s)", startTime, timeLayout)
	}

	if _, err := time.Parse(timeLayout, endTime); err != nil {
		return fmt.Errorf("invalid format: endTime: %s is not in the correct format (expected %s)", endTime, timeLayout)
	}

	return nil
}

func (params *QueryParams) RawQuery() string {
	return fmt.Sprintf(`
	SELECT 
		time_bucket('1 minute', ts) AS minute, 
		MAX(usage) AS max_cpu_usage, 
		MIN(usage) AS min_cpu_usage 
	FROM 
		cpu_usage 
	WHERE 
		host = '%s' AND 
		ts >= '%s' AND 
		ts <= '%s' 
	GROUP BY 
		minute 
	ORDER BY 
		minute;
	`, params.Hostname, params.StartTime, params.EndTime)
}
