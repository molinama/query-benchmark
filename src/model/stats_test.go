package model

import (
	"reflect"
	"testing"
	"time"
)

func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name               string
		queryTaskResults   []QueryTaskResult
		expectedQueryStats QueryStats
	}{
		{
			name: "Single Task Result",
			queryTaskResults: []QueryTaskResult{
				{Worker: 1, Hostname: "host1", Duration: 2 * time.Second},
			},
			expectedQueryStats: QueryStats{
				TotalQueries:        1,
				TotalProcessingTime: 2 * time.Second,
				MinQueryTime:        2 * time.Second,
				MedianQueryTime:     2 * time.Second,
				AvgQueryTime:        2 * time.Second,
				MaxQueryTime:        2 * time.Second,
				QueryWorkerStats: map[int]*QueryWorkerStats{
					1: {
						QueryStats: QueryStats{
							TotalQueries:        1,
							TotalProcessingTime: 2 * time.Second,
							MinQueryTime:        2 * time.Second,
							MedianQueryTime:     2 * time.Second,
							AvgQueryTime:        2 * time.Second,
							MaxQueryTime:        2 * time.Second,
						},
						QueryHostnameStats: map[string]*QueryStats{
							"host1": {
								TotalQueries:        1,
								TotalProcessingTime: 2 * time.Second,
								MinQueryTime:        2 * time.Second,
								MedianQueryTime:     2 * time.Second,
								AvgQueryTime:        2 * time.Second,
								MaxQueryTime:        2 * time.Second,
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple Task Results Single Worker",
			queryTaskResults: []QueryTaskResult{
				{Worker: 1, Hostname: "host1", Duration: 2 * time.Second},
				{Worker: 1, Hostname: "host1", Duration: 4 * time.Second},
				{Worker: 1, Hostname: "host2", Duration: 1 * time.Second},
			},
			expectedQueryStats: QueryStats{
				TotalQueries:        3,
				TotalProcessingTime: 7 * time.Second,
				MinQueryTime:        1 * time.Second,
				MedianQueryTime:     2 * time.Second,
				AvgQueryTime:        7 * time.Second / 3,
				MaxQueryTime:        4 * time.Second,
				QueryWorkerStats: map[int]*QueryWorkerStats{
					1: {
						QueryStats: QueryStats{
							TotalQueries:        3,
							TotalProcessingTime: 7 * time.Second,
							MinQueryTime:        1 * time.Second,
							MedianQueryTime:     2 * time.Second,
							AvgQueryTime:        7 * time.Second / 3,
							MaxQueryTime:        4 * time.Second,
						},
						QueryHostnameStats: map[string]*QueryStats{
							"host1": {
								TotalQueries:        2,
								TotalProcessingTime: 6 * time.Second,
								MinQueryTime:        2 * time.Second,
								MedianQueryTime:     3 * time.Second,
								AvgQueryTime:        3 * time.Second,
								MaxQueryTime:        4 * time.Second,
							},
							"host2": {
								TotalQueries:        1,
								TotalProcessingTime: 1 * time.Second,
								MinQueryTime:        1 * time.Second,
								MedianQueryTime:     1 * time.Second,
								AvgQueryTime:        1 * time.Second,
								MaxQueryTime:        1 * time.Second,
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple Task Results Multiple Workers",
			queryTaskResults: []QueryTaskResult{
				{Worker: 1, Hostname: "host1", Duration: 2 * time.Second},
				{Worker: 1, Hostname: "host2", Duration: 4 * time.Second},
				{Worker: 2, Hostname: "host1", Duration: 3 * time.Second},
				{Worker: 2, Hostname: "host2", Duration: 5 * time.Second},
			},
			expectedQueryStats: QueryStats{
				TotalQueries:        4,
				TotalProcessingTime: 14 * time.Second,
				MinQueryTime:        2 * time.Second,
				MedianQueryTime:     (3*time.Second + 4*time.Second) / 2,
				AvgQueryTime:        14 * time.Second / 4,
				MaxQueryTime:        5 * time.Second,
				QueryWorkerStats: map[int]*QueryWorkerStats{
					1: {
						QueryStats: QueryStats{
							TotalQueries:        2,
							TotalProcessingTime: 6 * time.Second,
							MinQueryTime:        2 * time.Second,
							MedianQueryTime:     3 * time.Second,
							AvgQueryTime:        3 * time.Second,
							MaxQueryTime:        4 * time.Second,
						},
						QueryHostnameStats: map[string]*QueryStats{
							"host1": {
								TotalQueries:        1,
								TotalProcessingTime: 2 * time.Second,
								MinQueryTime:        2 * time.Second,
								MedianQueryTime:     2 * time.Second,
								AvgQueryTime:        2 * time.Second,
								MaxQueryTime:        2 * time.Second,
							},
							"host2": {
								TotalQueries:        1,
								TotalProcessingTime: 4 * time.Second,
								MinQueryTime:        4 * time.Second,
								MedianQueryTime:     4 * time.Second,
								AvgQueryTime:        4 * time.Second,
								MaxQueryTime:        4 * time.Second,
							},
						},
					},
					2: {
						QueryStats: QueryStats{
							TotalQueries:        2,
							TotalProcessingTime: 8 * time.Second,
							MinQueryTime:        3 * time.Second,
							MedianQueryTime:     4 * time.Second,
							AvgQueryTime:        4 * time.Second,
							MaxQueryTime:        5 * time.Second,
						},
						QueryHostnameStats: map[string]*QueryStats{
							"host1": {
								TotalQueries:        1,
								TotalProcessingTime: 3 * time.Second,
								MinQueryTime:        3 * time.Second,
								MedianQueryTime:     3 * time.Second,
								AvgQueryTime:        3 * time.Second,
								MaxQueryTime:        3 * time.Second,
							},
							"host2": {
								TotalQueries:        1,
								TotalProcessingTime: 5 * time.Second,
								MinQueryTime:        5 * time.Second,
								MedianQueryTime:     5 * time.Second,
								AvgQueryTime:        5 * time.Second,
								MaxQueryTime:        5 * time.Second,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := QueryStats{}
			qs.QueryWorkerStats = make(map[int]*QueryWorkerStats)
			qs.CalculateStats(tt.queryTaskResults)
			if !reflect.DeepEqual(qs, tt.expectedQueryStats) {
				t.Errorf("CalculateStats() = %+v, want %+v", qs, tt.expectedQueryStats)
			}
		})
	}
}
