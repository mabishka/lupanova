package connloader

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

func TestConnLoader_GetShort(t *testing.T) {

	p := New("conn=a")
	haveFull := "full"
	user := uuid.New().String()
	haveShort, _ := p.GetShort(context.TODO(), haveFull, user)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		connName string
		// Named input parameters for target function.
		full    string
		short   string
		user    string
		wantErr bool
	}{
		{
			full:    haveFull,
			short:   haveShort,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			_, err := p.Load(context.TODO())
			assert.Error(t, err)
			short, gotErr := p.GetShort(context.TODO(), test.full, test.user)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, short, test.short)
				}
			}
		})
	}
}

func TestConnLoader_StoreList(t *testing.T) {

	p := New("conn=a")
	haveCorr := "aaa"
	haveFull := "full"
	user := uuid.New().String()
	haveShort, _ := p.GetShort(context.TODO(), haveFull, user)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		connName string
		// Named input parameters for target function.
		full    []model.FullItem
		short   []model.ShortItem
		user    string
		wantErr bool
	}{
		{
			full:    []model.FullItem{{Full: haveFull, Corr: haveCorr}},
			short:   []model.ShortItem{{Short: haveShort, Corr: haveCorr}},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New("conn=a")
			_, err := p.Load(context.TODO())
			assert.Error(t, err)
			short, gotErr := p.GetShortList(context.TODO(), test.full, test.user)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, short, test.short)
				}
			}
		})
	}
}
