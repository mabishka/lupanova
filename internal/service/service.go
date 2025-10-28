package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mabishka/lupanova/pkg/rand"
)

type StorageLoader interface {
	Load() (map[string]string, error) // return map [short string] full string
	Store(string, string) error       // store (full, short)
}

type Storage interface {
	GetShort(full string) (string, error)
	GetFull(short string) (string, error)
	Load(loader StorageLoader) error
}

type Server struct {
	*sync.RWMutex
	shortList map[string]string // map [short string] full string
	fullList  map[string]string // map [full string] short string
	loader    StorageLoader
}

const shortLen = 8

func New() *Server {
	return &Server{
		RWMutex:   &sync.RWMutex{},
		shortList: make(map[string]string),
		fullList:  make(map[string]string),
	}
}
func (p *Server) Load(loader StorageLoader) error {
	p.loader = loader
	list, err := loader.Load()
	if err != nil {
		return err
	}
	p.shortList = list

	for k, v := range p.shortList {
		p.fullList[v] = k
	}
	return nil
}

func (p *Server) store(full, short string) {
	p.shortList[short] = full
	p.fullList[full] = short
	if p.loader != nil {
		p.loader.Store(full, short)
	}
}

func (p *Server) GetShort(full string) (string, error) {
	p.Lock()
	defer p.Unlock()

	if short, ok := p.fullList[full]; ok {
		return short, nil
	}

	short, err := rand.CreateShort(shortLen)
	if err != nil {
		return "", err
	}

	p.store(full, short)
	p.shortList[short] = full
	return short, nil
}

func (p *Server) GetFull(short string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	short = strings.Trim(short, "/")
	if full, ok := p.shortList[short]; ok {
		return full, nil
	}

	return "", fmt.Errorf("path %s not found", short)

}
