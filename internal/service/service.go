package service

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/mabishka/lupanova/pkg/rand"
)

type Storage interface {
	GetShortUrl(full string) (string, error)
	GetFullUrl(short string) (string, error)
}

type Server struct {
	*sync.RWMutex
	list map[string]string // map [short string] full string
	u    url.URL
}

const shortLen = 8

func New(addr string) *Server {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	if host == "" {
		host = "localhost"
	}
	u := url.URL{
		Host:   net.JoinHostPort(host, port),
		Scheme: "http",
	}
	return &Server{RWMutex: &sync.RWMutex{}, list: make(map[string]string), u: u}
}

func (s *Server) createUrl(path string) string {
	s.u.Path = path
	return s.u.String()
}

func (p *Server) GetShortUrl(full string) (string, error) {
	p.Lock()
	defer p.Unlock()
	for k, v := range p.list {
		if v == full {
			return p.createUrl(k), nil
		}
	}

	short, err := rand.CreateShort(shortLen)
	if err != nil {
		return "", err
	}

	p.list[short] = full
	return p.createUrl(short), nil
}

func (p *Server) GetFullUrl(short string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	short = strings.Trim(short, "/")
	if full, ok := p.list[short]; ok {
		return full, nil
	}

	return "", fmt.Errorf("path %s not found", short)

}
