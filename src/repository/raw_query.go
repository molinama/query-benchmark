package repository

import "fmt"

func (params *QueryParams) RawQuery() string {
	return fmt.Sprintf(`
	SELECT 
		time_bucket('1 minute', time) AS minute, 
		MAX(cpu_usage) AS max_cpu_usage, 
		MIN(cpu_usage) AS min_cpu_usage 
	FROM 
		metrics 
	WHERE 
		hostname = '%s' AND 
		time >= '%s' AND 
		time <= '%s' 
	GROUP BY 
		minute 
	ORDER BY 
		minute;
	`, params.Hostname, params.StartTime, params.EndTime)
}
