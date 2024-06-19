package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParams_RawQuery(t *testing.T) {
	type args struct {
		params QueryParams
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid QueryParams",
			args: args{
				params: QueryParams{
					Hostname:  "host_000008",
					StartTime: "2017-01-01 08:59:22",
					EndTime:   "2017-01-01 09:59:22",
				},
			},
			want: `
		SELECT
			time_bucket('1 minute', time) AS minute,
			MAX(cpu_usage) AS max_cpu_usage,
			MIN(cpu_usage) AS min_cpu_usage
		FROM
			metrics
		WHERE
			hostname = 'host_000008' AND
			time >= '2017-01-01 08:59:22' AND
			time <= '2017-01-01 09:59:22'
		GROUP BY
			minute
		ORDER BY
			minute;
		`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.params.RawQuery())
		})
	}
}
