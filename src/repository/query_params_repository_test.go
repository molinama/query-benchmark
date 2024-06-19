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
		db      *sql.DB
		wantErr error
	}{
		{
			name:    "Valid DB",
			db:      &sql.DB{},
			wantErr: nil,
		},
		{
			name:    "Nil DB",
			db:      nil,
			wantErr: errors.New("cannot create repository"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQueryParamsRepository(tt.db)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewQueryParamsRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
