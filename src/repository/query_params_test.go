package repository

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewQueryParams(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name    string
		args    args
		want    *QueryParams
		wantErr error
	}{
		{
			name: "Valid Data",
			args: args{
				data: []string{"host_000008", "2017-01-01 08:59:22", "2017-01-01 09:59:22"},
			},
			want: &QueryParams{
				Hostname:  "host_000008",
				StartTime: "2017-01-01 08:59:22",
				EndTime:   "2017-01-01 09:59:22",
			},
		},
		{
			name: "Invalid Empty Hostname",
			args: args{
				data: []string{"", "2017-01-01 08:59:22", "2017-01-01 09:59:22"},
			},
			wantErr: errors.New("invalid format: hostname cannot be empty"),
		},
		{
			name: "Invalid Empty Hostname With Spaces",
			args: args{
				data: []string{"         ", "2017-01-01 08:59:22", "2017-01-01 09:59:22"},
			},
			wantErr: errors.New("invalid format: hostname cannot be empty"),
		},
		{
			name: "Valid Hostname With Spaces",
			args: args{
				data: []string{"    host_000008     ", "2017-01-01 08:59:22", "2017-01-01 09:59:22"},
			},
			want: &QueryParams{
				Hostname:  "host_000008",
				StartTime: "2017-01-01 08:59:22",
				EndTime:   "2017-01-01 09:59:22",
			},
		},
		{
			name: "Valid Start Time With Spaces",
			args: args{
				data: []string{"host_000008", "  2017-01-01 08:59:22  ", "2017-01-01 09:59:22"},
			},
			want: &QueryParams{
				Hostname:  "host_000008",
				StartTime: "2017-01-01 08:59:22",
				EndTime:   "2017-01-01 09:59:22",
			},
		},
		{
			name: "Valid End Time With Spaces",
			args: args{
				data: []string{"host_000008", "2017-01-01 08:59:22", "  2017-01-01 09:59:22  "},
			},
			want: &QueryParams{
				Hostname:  "host_000008",
				StartTime: "2017-01-01 08:59:22",
				EndTime:   "2017-01-01 09:59:22",
			},
		},
		{
			name: "InvalidDateFormat",
			args: args{
				data: []string{"host_000008", "2017-01-01", "2017-01-01 09:59:22"},
			},
			wantErr: errors.New("invalid format: startTime is not in the correct format (expected 2006-01-02 15:04:05)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQueryParams(tt.args.data)
			if err != nil {
				if !reflect.DeepEqual(err, tt.wantErr) {
					t.Errorf("NewQueryParams() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
