package service

import (
	"context"
	"fmt"

	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/pkg/rand"
)

type memLoader struct{}

func (p *memLoader) Load(ctx context.Context) (map[string]string, error) {
	return map[string]string{}, nil
}

func (p *memLoader) GetShortList(ctx context.Context, fullList []model.FullItem) (map[string]string, error) {
	resp := make(map[string]string)
	for _, v := range fullList {
		short, err := rand.CreateShort(config.ShortLen)
		if err != nil {
			return nil, err
		}
		resp[v.Full] = short
	}
	return resp, nil
}

func (p *memLoader) GetShort(ctx context.Context, full string) (string, error) {
	short, err := rand.CreateShort(config.ShortLen)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (p *memLoader) GetFull(ctx context.Context, short string) (string, error) {
	return "", fmt.Errorf("full not found for short %s", short)
}
