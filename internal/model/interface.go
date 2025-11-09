package model

import "context"

type StorageLoader interface {
	Load(ctx context.Context) (map[string]string, error)        // return map [short string] full string
	Store(ctx context.Context, full string, short string) error // store (full, short)
	StoreList(ctx context.Context, list []StoreItem) error      // store [](full, short)
}

type ConnLoader interface {
	Create(context.Context) error
	Ping(context.Context) error
}

type Storage interface {
	GetShort(ctx context.Context, full string) (string, error)
	GetFull(ctx context.Context, short string) (string, error)

	GetShortList(ctx context.Context, full []FullItem) ([]ShortItem, error)
	GetFullList(ctx context.Context, short []ShortItem) ([]FullItem, error)

	Load(ctx context.Context, loader StorageLoader) error
}

