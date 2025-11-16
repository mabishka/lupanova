package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mabishka/lupanova/pkg/utils"
)

func TestCreateShort(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		n    int
		want int
	}{
		{
			name: "length",
			n:    10,
			want: 10,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := utils.CreateShort(test.n)

			assert.NoError(t, err)
			assert.Equal(t, test.want, len(got))
		})
	}
}
