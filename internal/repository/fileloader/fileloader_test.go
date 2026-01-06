package fileloader

import (
	"context"
	"testing"

	"github.com/mabishka/lupanova/internal/model"
	"github.com/stretchr/testify/assert"
)

const filestorelist = "../../../../filestore_list.json"
const filecreate = "../../../../file_create.json"
const defaultFileName = "../../../storage.json"

func TestFileLoader_exist(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		fileName string
		want     bool
	}{
		{
			fileName: "main.go",
			want:     false,
		},
		{
			fileName: "go.mod",
			want:     false,
		},
		{
			fileName: "fileloader_test.go",
			want:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New(test.fileName)
			got, err := p.exist()
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestFileLoader_create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		fileName string
		wantErr  bool
	}{
		{
			fileName: filecreate,
			wantErr:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New(test.fileName)
			gotErr := p.create()
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestFileLoader_Load(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		fileName string
		want     map[string]string
		wantErr  bool
	}{
		{
			fileName: filestorelist,
			want:     map[string]string{},
			wantErr:  false,
		},
		{
			fileName: "fileloader_test.go",
			want:     map[string]string{},
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := New(test.fileName)
			_, gotErr := p.Load(context.TODO())
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestFileLoader_GetShort(t *testing.T) {
	p := New(defaultFileName)
	_, err := p.Load(context.TODO())
	assert.NoError(t, err)

	haveFull := "full"
	haveShort := "short"
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		fileName string
		// Named input parameters for target function.
		full    string
		short   string
		user    string
		wantErr bool
	}{
		{
			full:    haveFull,
			short:   haveShort,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			short, gotErr := p.GetShort(context.TODO(), test.full, test.user)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.NotEmpty(t, short)
				}
			}
		})
	}
}

func TestFileLoader_GetShortList(t *testing.T) {
	p := New(defaultFileName)
	_, err := p.Load(context.TODO())
	assert.NoError(t, err)

	haveCorr := "aaa"
	haveFull := "full"
	haveShort := "short"
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		// Named input parameters for target function.
		full    []model.FullItem
		user    string
		short   map[string]string
		wantErr bool
	}{
		{
			full:    []model.FullItem{{Full: haveFull, Corr: haveCorr}},
			short:   map[string]string{haveFull: haveShort},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			short, gotErr := p.GetShortList(context.TODO(), test.full, test.user)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.NotEmpty(t, short)
				}
			}
		})
	}
}
