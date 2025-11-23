package model

import (
	"context"
)

type StorageLoader interface {
	Load(ctx context.Context) (map[string]string, error)
	GetShortList(ctx context.Context, fullList []FullItem) (map[string]string, error)
	GetShort(ctx context.Context, full string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
}

type ConnLoader interface {
	Create(context.Context) error
	Ping(context.Context) error
}

type Storage interface {
	GetShortList(ctx context.Context, full []FullItem) ([]ShortItem, error)
	GetShort(ctx context.Context, full string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)

	Load(ctx context.Context, loader StorageLoader) error
}
