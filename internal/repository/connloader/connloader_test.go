package connloader

import (
	"context"
	"testing"

	"github.com/mabishka/lupanova/internal/model"
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
			name:     "negative",
			connName: "conn",
			wantErr:  true,
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
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New("conn=a")
			_, err := p.Load(context.TODO())
			assert.Error(t, err)
			gotErr := p.Store(context.TODO(), test.full, test.short)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestConnLoader_StoreList(t *testing.T) {

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		connName string
		// Named input parameters for target function.
		store   []model.StoreItem
		wantErr bool
	}{
		{
			store:   []model.StoreItem{{Full: "full", Short: "short"}},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New("conn=a")
			_, err := p.Load(context.TODO())
			assert.Error(t, err)
			gotErr := p.StoreList(context.TODO(), test.store)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
