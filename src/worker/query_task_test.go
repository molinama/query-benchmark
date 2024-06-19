package worker

import (
	"errors"
	"testing"
	"time"

	"github.com/molinama/timescale/src/model"
	"github.com/molinama/timescale/src/repository"
	"github.com/stretchr/testify/assert"
)

// MockQueryParamsRepository is a mock implementation of QueryParamsRepository
type MockQueryParamsRepository struct {
	shouldFail bool
}

// NewMockQueryParamsRepository creates a new instance of MockQueryParamsRepository
func NewMockQueryParamsRepository(shouldFail bool) *MockQueryParamsRepository {
	return &MockQueryParamsRepository{
		shouldFail: shouldFail,
	}
}

// RawQuery simulates the RawQuery method
func (m *MockQueryParamsRepository) RawQuery(params *repository.QueryParams) (time.Duration, error) {
	if m.shouldFail {
		return -1, errors.New("mock query error")
	}

	// Simulate a delay to test timing
	time.Sleep(100 * time.Millisecond)
	return 100 * time.Millisecond, nil
}

func TestQueryTask_Execute(t *testing.T) {
	tests := []struct {
		name        string
		shouldFail  bool
		expectError bool
	}{
		{
			name:        "Successful Query",
			shouldFail:  false,
			expectError: false,
		},
		{
			name:        "Query with Error",
			shouldFail:  true,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultsCh := make(chan model.QueryTaskResult, 1)

			mockRepo := NewMockQueryParamsRepository(tt.shouldFail)
			params := &repository.QueryParams{
				Hostname:  "host_000008",
				StartTime: "2017-01-01 08:59:22",
				EndTime:   "2017-01-01 09:59:22",
			}
			task := NewQueryTask(mockRepo, params, resultsCh)

			go task.Execute(1)

			result := <-resultsCh
			assert.Equal(t, 1, result.Worker)
			assert.Equal(t, params.Hostname, result.Hostname)
			assert.Equal(t, params.RawQuery(), result.RawQuery)
			assert.NotEqual(t, 0, result.Duration)
			if tt.expectError {
				assert.Error(t, result.Err)
			} else {
				assert.NoError(t, result.Err)
			}
		})
	}
}
