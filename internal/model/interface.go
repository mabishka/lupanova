package model

import (
	"context"
)

type StorageLoader interface {
	Load(ctx context.Context) (map[string]string, error)
	GetShortList(ctx context.Context, fullList []FullItem, user string) (map[string]string, error)
	GetShort(ctx context.Context, full string, user string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
	GetUserList(ctx context.Context, user string) ([]StoreItem, error)
	DeleteList(context.Context, []string, string) error
}

type ConnLoader interface {
	Create(context.Context) error
	Ping(context.Context) error
}

type Storage interface {
	GetUserList(ctx context.Context, user string) ([]StoreItem, error)
	GetShortList(ctx context.Context, full []FullItem, user string) ([]ShortItem, error)
	GetShort(ctx context.Context, full string, user string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)
	DeleteList(ctx context.Context, short []string, user string) error

	Load(ctx context.Context, loader StorageLoader) error
}
