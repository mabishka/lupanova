package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_run(t *testing.T) {

	tests := []struct {
		name    string // description of this test case
		wantErr bool
	}{
		{
			name:    "positive",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			go func() {
				err = run()
			}()
			time.Sleep(time.Second * 5)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
