package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setAddress(t *testing.T) {
	type have struct {
		envAddress     string
		flagName       string
		defaultAddress string
		description    string
	}
	tests := []struct {
		name string // description of this test case
		pre  func()
		have have
		want string
	}{
		{
			name: "positiveDefaultAddress",
			pre:  func() {},
			have: have{
				envAddress:     envBaseAddress,
				flagName:       flagBaseAddress,
				defaultAddress: defaultBaseAddress,
				description:    descBaseAddress,
			},
			want: defaultBaseAddress,
		},
		{
			name: "positiveEnvAddress",
			pre:  func() { os.Setenv(envBaseAddress, defaultBaseAddress) },
			have: have{
				envAddress:     envBaseAddress,
				flagName:       flagBaseAddress,
				defaultAddress: defaultBaseAddress,
				description:    descBaseAddress,
			},
			want: defaultBaseAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.pre()
			got := setAddress(test.have.envAddress, test.have.flagName, test.have.defaultAddress, test.have.description)

			assert.Equal(t, test.want, *got)
		})
	}
}

func Test_validateServerAddress(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		address        string
		defaultAddress string
		want           string
	}{
		{
			name:           "positiveAddress",
			address:        defaultServerAddress,
			defaultAddress: defaultServerAddress,
			want:           defaultServerAddress,
		},
		{
			name:           "positiveDefault",
			address:        "",
			defaultAddress: defaultServerAddress,
			want:           defaultServerAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := validateServerAddress(test.address, test.defaultAddress)
			assert.Equal(t, test.want, got)
		})
	}
}

func Test_validateBaseAddress(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		address        string
		defaultAddress string
		want           string
	}{
		{
			name:           "positiveAddress",
			address:        defaultBaseAddress,
			defaultAddress: defaultBaseAddress,
			want:           defaultBaseAddress,
		},
		{
			name:           "positiveDefault",
			address:        "",
			defaultAddress: defaultBaseAddress,
			want:           defaultBaseAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := validateBaseAddress(test.address, test.defaultAddress)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want *Config
	}{
		{
			name: "positiveDefault",
			want: &Config{
				serverAddress: defaultServerAddress,
				baseAddress:   defaultBaseAddress,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New()

			assert.Equal(t, test.want.serverAddress, got.serverAddress)
			assert.Equal(t, test.want.baseAddress, got.baseAddress)
		})
	}
}
