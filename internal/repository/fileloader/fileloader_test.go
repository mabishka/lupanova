package fileloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			got := p.exist()
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
			fileName: "file1",
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
			fileName: "file1",
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
			got, gotErr := p.Load()
			if test.wantErr {
				assert.Error(t, gotErr)
			} else if assert.NoError(t, gotErr) {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestFileLoader_Store(t *testing.T) {
	p := New("file2")
	_, err := p.Load(); 
	assert.NoError(t, err)
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		fileName string
		// Named input parameters for target function.
		full    string
		short   string
		wantErr bool
	}{
		{
			short: "short",
			full: "full",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := p.Store(test.full, test.short)
			if test.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
