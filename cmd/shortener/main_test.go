package main

import "testing"

func Test_main(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "run",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			main()
		})
	}
}
