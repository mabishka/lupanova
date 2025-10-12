package config

import (
	"flag"
	"net"
	"strings"
)

// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
// Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например, http://localhost:8000/qsd54gFg).

const (
	defaultServerAddress = ""
	defaultServerPort    = "8080"
	defaultBaseAddress   = "localhost"
	defaultBasePort      = "8080"
)

type Config struct {
	serverAddress string
	baseAddress   string
}

func New() *Config {
	res := &Config{}
	flag.StringVar(&res.serverAddress, "a", ":8080", "адрес запуска HTTP-сервера")
	flag.StringVar(&res.baseAddress, "b", "localhost:8080", "базовый адрес результирующего сокращённого URL")

	flag.Parse()

	res.serverAddress = validateAddress(res.serverAddress, defaultServerAddress, defaultServerPort)
	res.baseAddress = validateAddress(res.baseAddress, defaultBaseAddress, defaultBasePort)

	return res
}

func validateAddress(address, defaultAddress, defaultPort string) string {
	addrList := strings.Split(address, ":")
	if len(addrList) < 1 || len(addrList) > 2 {
		return net.JoinHostPort(defaultAddress, defaultPort)
	}
	if addrList[0] == "" {
		addrList[0] = defaultAddress
	}
	if len(addrList) < 2 || addrList[1] == "" {
		return net.JoinHostPort(addrList[0], defaultPort)
	}
	return net.JoinHostPort(addrList[0], addrList[1])
}

func (c *Config) GetBaseAddress() string {
	return c.baseAddress
}

func (c *Config) GetServerAddress() string {
	return c.serverAddress
}
