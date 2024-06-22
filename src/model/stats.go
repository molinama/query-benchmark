package model

import (
	"fmt"
	"sort"
	"time"
)

type Stats struct {
	queryStats
	QueryErrorStats
}

type queryStats struct {
	TotalSuccess        int
	TotalProcessingTime time.Duration
	MinQueryTime        time.Duration
	MedianQueryTime     time.Duration
	AvgQueryTime        time.Duration
	MaxQueryTime        time.Duration
	QueryWorkerStats    map[int]*queryWorkerStats
}

type queryWorkerStats struct {
	queryStats
	QueryHostnameStats map[string]*queryStats
}

type QueryErrorStats struct {
	TotalErrs     int
	QueryTaskErrs []QueryTaskErr
}

func (qs Stats) String() string {
	return fmt.Sprintf(
		"\nSTATS\n"+
			"\nTotal Queries: %d"+
			"\nNumber of queries successfully processed: %d\n"+
			"Total processing time: %v\n"+
			"Minimum query time: %v\n"+
			"Median query time: %v\n"+
			"Average query time: %v\n"+
			"Maximum query time: %v\n"+
			"Total Errors: %d\n",
		qs.TotalSuccess+qs.TotalErrs,
		qs.TotalSuccess,
		qs.TotalProcessingTime,
		qs.MinQueryTime,
		qs.MedianQueryTime,
		qs.AvgQueryTime,
		qs.MaxQueryTime,
		qs.TotalErrs,
	)
}

func (qs *Stats) CalculateStats(queryTaskResults []QueryTaskResult, queryTaskErrs []QueryTaskErr) {
	qs.TotalSuccess = len(queryTaskResults)
	qs.TotalErrs = len(queryTaskErrs)
	qs.QueryTaskErrs = queryTaskErrs

	if qs.TotalSuccess == 0 {
		return
	}
	queryTimes := make([]time.Duration, 0, qs.TotalSuccess)
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

func (qs *queryStats) calculateAllStats(queryTimes []time.Duration, queryWorkerTimes map[int][]time.Duration, queryHostnameTimes map[int]map[string][]time.Duration) {
	qs.calculateStats(queryTimes)
	qs.QueryWorkerStats = make(map[int]*queryWorkerStats)

	for worker, queryWorkerTime := range queryWorkerTimes {
		workerStats := queryWorkerStats{}
		workerStats.calculateStats(queryWorkerTime)
		qs.QueryWorkerStats[worker] = &workerStats
		qs.QueryWorkerStats[worker].QueryHostnameStats = make(map[string]*queryStats)

		for hostname, queryHostnameTime := range queryHostnameTimes[worker] {
			hostnameStats := queryStats{}
			hostnameStats.calculateStats(queryHostnameTime)
			qs.QueryWorkerStats[worker].QueryHostnameStats[hostname] = &hostnameStats
		}

	}
}

func (qs *queryStats) calculateStats(queryTimes []time.Duration) {
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

	qs.TotalProcessingTime = totalProcessingTime
	qs.MinQueryTime = minQueryTime
	qs.MedianQueryTime = medianQueryTime
	qs.AvgQueryTime = avgQueryTime
	qs.MaxQueryTime = maxQueryTime
}
