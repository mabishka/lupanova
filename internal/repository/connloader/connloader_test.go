package connloader

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnLoader_Ping(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		connName string
		wantErr  bool
	}{
		{
			name:     "positive",
			connName: "user=user password=user host=localhost port=5433 dbname=practicum sslmode=disable",
			wantErr:  false,
		},
		{
			name:     "negative",
			connName: "conn",
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New(test.connName)
			err := p.Ping(context.TODO())

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConnLoader_create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		connName string
		wantErr  bool
	}{
		{
			name:     "positive",
			connName: "user=user password=user host=localhost port=5433 dbname=practicum sslmode=disable",
			wantErr:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New(test.connName)
			gotErr := p.Create(context.TODO())
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestConnLoader_Load(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		connName string
		want     map[string]string
		wantErr  bool
	}{
		{
			name:     "positive",
			connName: "user=user password=user host=localhost port=5433 dbname=practicum sslmode=disable",
			want:     map[string]string{"short": "full"},
			wantErr:  false,
		},
		{
			name:     "negative",
			connName: "err",
			want:     map[string]string{},
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New(test.connName)
			_, gotErr := p.Load(context.TODO())
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestConnLoader_Store(t *testing.T) {
	p := New("user=user password=user host=localhost port=5433 dbname=practicum sslmode=disable")
	_, err := p.Load(context.TODO())
	assert.NoError(t, err)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		connName string
		// Named input parameters for target function.
		full    string
		short   string
		wantErr bool
	}{
		{
			short:   "short",
			full:    "full",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := p.Store(context.TODO(), test.full, test.short)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
