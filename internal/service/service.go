package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mabishka/lupanova/pkg/rand"
)

type Storage interface {
	GetShort(full string) (string, error)
	GetFull(short string) (string, error)
}

type Server struct {
	*sync.RWMutex
	list map[string]string // map [short string] full string
}

const shortLen = 8

func New() *Server {
	return &Server{RWMutex: &sync.RWMutex{}, list: make(map[string]string)}
}

func (p *Server) GetShort(full string) (string, error) {
	p.Lock()
	defer p.Unlock()
	for k, v := range p.list {
		if v == full {
			return k, nil
		}
	}

	short, err := rand.CreateShort(shortLen)
	if err != nil {
		return "", err
	}

	p.list[short] = full
	return short, nil
}

func (p *Server) GetFull(short string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	short = strings.Trim(short, "/")
	if full, ok := p.list[short]; ok {
		return full, nil
	}

	return "", fmt.Errorf("path %s not found", short)

}
