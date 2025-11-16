package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

type mock struct {
}

func (p *mock) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	return nil, errors.New("unsupported")
}

func (p *mock) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, nil
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		conn    Connector
		wantErr bool
	}{
		{
			conn:    &mock{},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Create(context.Background(), test.conn)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		conn    Connector
		want    map[string]string
		wantErr bool
	}{
		{
			conn:    &mock{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadList(context.Background(), test.conn)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetShort(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		conn    Connector
		full    string
		want    string
		wantErr bool
	}{
		{
			conn:    &mock{},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetShort(context.Background(), test.conn, test.full)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetFull(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		conn    Connector
		short   string
		want    string
		wantErr bool
	}{
		{
			conn:    &mock{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetFull(context.Background(), test.conn, test.short)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_store(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		conn    Connector
		full    string
		short   string
		wantErr bool
	}{
		{
			conn:    &mock{},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store(context.Background(), test.conn, test.full, test.short)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_getShort(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		conn    Connector
		full    string
		want    string
		wantErr bool
	}{
		{
			conn:    &mock{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := getShort(context.Background(), test.conn, test.full)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
