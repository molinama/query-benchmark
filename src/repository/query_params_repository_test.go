package repository

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
)

func TestNewQueryParamsRepository(t *testing.T) {
	tests := []struct {
		name    string
		conn    *sql.DB
		wantErr error
	}{
		{
			name:    "Valid DB",
			conn:    &sql.DB{},
			wantErr: nil,
		},
		{
			name:    "Nil DB",
			conn:    nil,
			wantErr: errors.New("cannot create repository"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQueryParamsRepository(tt.conn)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewQueryParamsRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
