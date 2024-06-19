package repository

func (params *QueryParams) PrepareStatement() string {
	return `
		SELECT time_bucket('1 minute', ts) AS minute,
		MAX(usage) AS max_usage,
		MIN(usage) AS min_usage
		FROM cpu_usage
		WHERE host = ?
			AND ts >= ?	
			AND ts < ?	
		GROUP BY minute
		ORDER BY minute;
	`
}
