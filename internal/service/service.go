package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/pkg/rand"
)

type Server struct {
	*sync.RWMutex
	shortList map[string]string // map [short string] full string
	fullList  map[string]string // map [full string] short string
	loader    model.StorageLoader
}

const shortLen = 8

func New() *Server {
	return &Server{
		RWMutex:   &sync.RWMutex{},
		shortList: make(map[string]string),
		fullList:  make(map[string]string),
	}
}

func (p *Server) Load(ctx context.Context, loader model.StorageLoader) error {
	p.loader = loader
	list, err := loader.Load(ctx)
	if err != nil {
		return err
	}
	p.shortList = list

	for k, v := range p.shortList {
		p.fullList[v] = k
	}
	return nil
}

func (p *Server) store(ctx context.Context, full, short string) error {
	p.shortList[short] = full
	p.fullList[full] = short
	if p.loader != nil {
		if err := p.loader.Store(ctx, full, short); err != nil {
			return err
		}
	}
	return nil
}

func (p *Server) GetShort(ctx context.Context, full string) (string, error) {
	p.Lock()
	defer p.Unlock()

	if short, ok := p.fullList[full]; ok {
		return short, nil
	}

	short, err := rand.CreateShort(shortLen)
	if err != nil {
		return "", err
	}

	if err := p.store(ctx, full, short); err != nil {
		return "", err
	}
	p.shortList[short] = full
	return short, nil
}

func (p *Server) GetFull(ctx context.Context, short string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	short = strings.Trim(short, "/")
	if full, ok := p.shortList[short]; ok {
		return full, nil
	}

	return "", fmt.Errorf("path %s not found", short)

}
