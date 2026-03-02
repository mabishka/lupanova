package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*

func Test_setAddress(t *testing.T) {

	tests := []struct {
		name string // description of this test case
		pre  func()
		have have
		want string
	}{
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
			got := setParamString(test.have.envAddress, test.have.flagName, test.have.defaultAddress, test.have.description)

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

func Test_setParamString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		address        string
		defaultAddress string
		want           string
	}{
		{
			name:           "positiveDefault",
			address:        "",
			defaultAddress: defaultServerAddress,
			want:           defaultServerAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := setParamString(envServerAddress, flagServerAddress, defaultServerAddress, descServerAddress)
			assert.NotEmpty(t, got)
		})
	}
}


func Test_validateBaseAddress(t *testing.T) {
	cfg := New()
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
*/

func TestConfig_GetBaseAddress(t *testing.T) {

	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: defaultBaseAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetBaseAddress())
		})
	}
}

func TestConfig_GetServerAddress(t *testing.T) {

	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: defaultServerAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetServerAddress())
		})
	}
}

func TestConfig_GetLogLevel(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: defaultLogLevel,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetLogLevel())
		})
	}
}

func TestConfig_GetFileName(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetFileName())
		})
	}
}

func TestConfig_GetConnAddress(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetConnAddress())
		})
	}
}

func TestConfig_GetAuditFile(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetAuditFile())
		})
	}
}

func TestConfig_GetAuditAddress(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want string
	}{
		{
			name: "positive",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.GetAuditAddress())
		})
	}
}

func TestConfig_IsEnableHTTPS(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string // description of this test case
		want bool
	}{
		{
			name: "positive",
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, cfg.IsEnableHTTPS())
		})
	}
}

func TestConfig_getString(t *testing.T) {
	cfg := New()
	tests := []struct {
		name  string
		value configType
		want  string
	}{
		{
			name:  "positive",
			value: configLogLevel,
			want:  defaultLogLevel,
		},
		{
			name:  "negative_name",
			value: "empty",
			want:  "",
		},
		{
			name:  "negative_bool",
			value: configEnableHttps,
			want:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, cfg.getString(test.value), test.want)
		})
	}
}

func TestConfig_getBool(t *testing.T) {
	cfg := New()
	tests := []struct {
		name  string
		value configType
		want  bool
	}{
		{
			name:  "positive",
			value: configEnableHttps,
			want:  false,
		},
		{
			name:  "negative_name",
			value: "empty",
			want:  false,
		},
		{
			name:  "negative_string",
			value: configLogLevel,
			want:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, cfg.getBool(test.value), test.want)
		})
	}
}
