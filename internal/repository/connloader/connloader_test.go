package connloader

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mabishka/lupanova/internal/model"
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

func (p *mock) PingContext(context.Context) error {
	return nil
}

func (p *mock) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, errors.New("unsupported")
}

func TestConnLoader_Ping(t *testing.T) {
	tests := []struct {
		name     string
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
		name     string
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
		name     string
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

	p := ConnLoader{conn: &mock{}}
	haveFull := "full"
	user := uuid.New().String()
	haveShort, _ := p.GetShort(context.TODO(), haveFull, user)
	tests := []struct {
		name     string
		connName string
		full     string
		short    string
		user     string
		wantErr  bool
	}{
		{
			full:    haveFull,
			short:   haveShort,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

	p := ConnLoader{conn: &mock{}}
	haveCorr := "aaa"
	haveFull := "full"
	user := uuid.New().String()
	haveShort, _ := p.GetShort(context.TODO(), haveFull, user)
	tests := []struct {
		name     string
		connName string
		full     []model.FullItem
		short    []model.ShortItem
		user     string
		wantErr  bool
	}{
		{
			full:    []model.FullItem{{Full: haveFull, Corr: haveCorr}},
			short:   []model.ShortItem{{Short: haveShort, Corr: haveCorr}},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

func TestConnLoader_GetUserList(t *testing.T) {
	p := ConnLoader{conn: &mock{}}
	user := uuid.New().String()
	tests := []struct {
		name    string
		user    string
		wantErr bool
	}{
		{
			user:    user,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			list, gotErr := p.GetUserList(context.Background(), test.user)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, 0, len(list))
			}
		})
	}
}

func TestConnLoader_deleteList(t *testing.T) {
	p := ConnLoader{conn: &mock{}}
	user := uuid.New().String()
	haveFull := "full"
	haveShort, _ := p.GetShort(context.TODO(), haveFull, user)

	dataShort := []string{haveShort}
	tests := []struct {
		name    string
		short   []string
		user    string
		wantErr bool
	}{
		{
			short:   dataShort,
			user:    user,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := p.deleteList(context.Background(), test.short, test.user)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestConnLoader_DeleteList(t *testing.T) {
	p := ConnLoader{conn: &mock{}}
	user := uuid.New().String()
	haveFull := "full"
	haveShort, _ := p.GetShort(context.TODO(), haveFull, user)
	tests := []struct {
		name    string
		addr    string
		short   []string
		user    string
		wantErr bool
	}{
		{
			short:   []string{haveShort},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := p.DeleteList(context.Background(), test.short, test.user)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestConnLoader_GetFull(t *testing.T) {
	p := ConnLoader{conn: &mock{}}
	haveFull := "full"
	user := uuid.New().String()
	haveShort, _ := p.GetShort(context.TODO(), haveFull, user)
	tests := []struct {
		name    string
		addr    string
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

func TestConnLoader_GetStat(t *testing.T) {
	p := ConnLoader{conn: &mock{}}
	tests := []struct {
		name    string
		addr    string
		full    string
		short   string
		user    string
		wantErr bool
	}{
		{
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a, b, gotErr := p.GetStat(context.TODO())
			if test.wantErr {
				assert.Error(t, gotErr)
				return
			}
			assert.NoError(t, gotErr)
			assert.NotEmpty(t, a)
			assert.NotEmpty(t, b)

		})
	}
}
