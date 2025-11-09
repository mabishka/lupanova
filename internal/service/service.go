package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/pkg/rand"
	"go.uber.org/zap"
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
		logger.Log().Error("Server.Load", zap.Error(err))
		return err
	}
	p.shortList = list

	for k, v := range p.shortList {
		p.fullList[v] = k
	}
	return nil
}

func (p *Server) store(ctx context.Context, full, short string) error {
	if p.loader != nil {
		if err := p.loader.Store(ctx, full, short); err != nil {
			logger.Log().Error("Server.store", zap.Error(err))
			return err
		}
	}
	return nil
}

func (p *Server) storelist(ctx context.Context, list []model.StoreItem) error {
	if p.loader != nil {
		if err := p.loader.StoreList(ctx, list); err != nil {
			logger.Log().Error("Server.storelist", zap.Error(err))
			return err
		}
	}
	return nil
}

func checkFull(full string) error {
	if _, err := url.ParseRequestURI(full); err != nil {
		logger.Log().Error("checkFull", zap.Error(err))
		return err
	}
	return nil
}

func (p *Server) GetShortList(ctx context.Context, fullList []model.FullItem) ([]model.ShortItem, error) {
	shortList := make([]model.ShortItem, 0, len(fullList))
	storeList := make([]model.StoreItem, 0, len(fullList))
	for _, v := range fullList {

		if err := checkFull(v.Full); err != nil {
			return nil, err
		}

		if short, ok := p.fullList[v.Full]; ok {
			shortList = append(shortList, model.ShortItem{Corr: v.Corr, Short: short})
			continue
		}

		short, err := rand.CreateShort(shortLen)
		if err != nil {
			return nil, err
		}

		p.shortList[short] = v.Full
		p.fullList[v.Full] = short

		shortList = append(shortList, model.ShortItem{Corr: v.Corr, Short: short})
		storeList = append(storeList, model.StoreItem{Full: v.Full, Short: short})
	}

	if err := p.storelist(ctx, storeList); err != nil {
		return nil, err
	}
	return shortList, nil
}

func (p *Server) GetShort(ctx context.Context, full string) (string, error) {
	p.Lock()
	defer p.Unlock()

	if err := checkFull(full); err != nil {
		return "", err
	}

	if short, ok := p.fullList[full]; ok {
		return short, nil
	}

	short, err := rand.CreateShort(shortLen)
	if err != nil {
		return "", err
	}

	p.shortList[short] = full
	p.fullList[full] = short

	if err := p.store(ctx, full, short); err != nil {
		return "", err
	}

	return short, nil
}

func (p *Server) GetFullList(ctx context.Context, shortList []model.ShortItem) ([]model.FullItem, error) {
	fullList := make([]model.FullItem, 0, len(shortList))
	for _, v := range shortList {
		full, err := p.GetFull(ctx, v.Short)
		if err != nil {
			return nil, err
		}
		fullList = append(fullList, model.FullItem{Corr: v.Corr, Full: full})
	}
	return fullList, nil
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
