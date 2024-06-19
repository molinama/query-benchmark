package model

import (
	"fmt"
	"sort"
	"time"
)

type QueryStats struct {
	TotalQueries        int
	TotalProcessingTime time.Duration
	MinQueryTime        time.Duration
	MedianQueryTime     time.Duration
	AvgQueryTime        time.Duration
	MaxQueryTime        time.Duration
	QueryWorkerStats    map[int]*QueryWorkerStats
}

type QueryWorkerStats struct {
	QueryStats
	QueryHostnameStats map[string]*QueryStats
}

type QueryError struct {
	Err      error
	RawQuery string
}

func (qs QueryStats) String() string {
	return fmt.Sprintf(
		"Number of queries processed: %d\n"+
			"Total processing time: %v\n"+
			"Minimum query time: %v\n"+
			"Median query time: %v\n"+
			"Average query time: %v\n"+
			"Maximum query time: %v",
		qs.TotalQueries,
		qs.TotalProcessingTime,
		qs.MinQueryTime,
		qs.MedianQueryTime,
		qs.AvgQueryTime,
		qs.MaxQueryTime,
	)
}

func (qs *QueryStats) CalculateStats(queryTaskResults []QueryTaskResult) {
	totalQueries := len(queryTaskResults)

	if totalQueries == 0 {
		return
	}
	queryTimes := make([]time.Duration, 0, totalQueries)
	queryWorkerTimes := make(map[int][]time.Duration)
	queryHostnameTimes := make(map[int]map[string][]time.Duration)

	for _, result := range queryTaskResults {
		queryTimes = append(queryTimes, result.Duration)

		if _, exists := queryWorkerTimes[result.Worker]; !exists {
			queryWorkerTimes[result.Worker] = []time.Duration{}
		}
		queryWorkerTimes[result.Worker] = append(queryWorkerTimes[result.Worker], result.Duration)

		if _, exists := queryHostnameTimes[result.Worker]; !exists {
			queryHostnameTimes[result.Worker] = make(map[string][]time.Duration)
		}
		if _, exists := queryHostnameTimes[result.Worker][result.Hostname]; !exists {
			queryHostnameTimes[result.Worker][result.Hostname] = []time.Duration{}
		}
		queryHostnameTimes[result.Worker][result.Hostname] = append(queryHostnameTimes[result.Worker][result.Hostname], result.Duration)
	}

	qs.calculateAllStats(queryTimes, queryWorkerTimes, queryHostnameTimes)
}

func (qs *QueryStats) calculateAllStats(queryTimes []time.Duration, queryWorkerTimes map[int][]time.Duration, queryHostnameTimes map[int]map[string][]time.Duration) {
	qs.calculateStats(queryTimes)
	qs.QueryWorkerStats = make(map[int]*QueryWorkerStats)

	for worker, queryWorkerTime := range queryWorkerTimes {
		workerStats := QueryWorkerStats{}
		workerStats.calculateStats(queryWorkerTime)
		qs.QueryWorkerStats[worker] = &workerStats
		qs.QueryWorkerStats[worker].QueryHostnameStats = make(map[string]*QueryStats)

		for hostname, queryHostnameTime := range queryHostnameTimes[worker] {
			hostnameStats := QueryStats{}
			hostnameStats.calculateStats(queryHostnameTime)
			qs.QueryWorkerStats[worker].QueryHostnameStats[hostname] = &hostnameStats
		}

	}
}

func (qs *QueryStats) calculateStats(queryTimes []time.Duration) {
	totalQueries := len(queryTimes)
	if totalQueries == 0 {
		return
	}

	var totalProcessingTime time.Duration
	minQueryTime := queryTimes[0]
	maxQueryTime := queryTimes[0]

	for _, qt := range queryTimes {
		totalProcessingTime += qt
		if qt < minQueryTime {
			minQueryTime = qt
		}
		if qt > maxQueryTime {
			maxQueryTime = qt
		}
	}

	avgQueryTime := totalProcessingTime / time.Duration(totalQueries)

	sort.Slice(queryTimes, func(i, j int) bool { return queryTimes[i] < queryTimes[j] })
	medianQueryTime := queryTimes[totalQueries/2]
	if totalQueries%2 == 0 {
		medianQueryTime = (queryTimes[totalQueries/2-1] + queryTimes[totalQueries/2]) / 2
	}

	qs.TotalQueries = totalQueries
	qs.TotalProcessingTime = totalProcessingTime
	qs.MinQueryTime = minQueryTime
	qs.MedianQueryTime = medianQueryTime
	qs.AvgQueryTime = avgQueryTime
	qs.MaxQueryTime = maxQueryTime
}
