package rand_test

import (
	"testing"

	"github.com/mabishka/lupanova/pkg/rand"
	"github.com/stretchr/testify/assert"
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
			got, err := rand.CreateShort(test.n)

			assert.NoError(t, err)
			assert.Equal(t, test.want, len(got))
		})
	}
}
