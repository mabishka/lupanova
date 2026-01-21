package handler

import (
	"net/url"

	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/service"
)

type StorageServer struct {
	model.Storage
	u *url.URL
}

func New(address string) *StorageServer {
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	return &StorageServer{Storage: service.New(), u: u}
}

func (p *StorageServer) SetLoader(loader model.Storage) {
	p.Storage = loader
}

func (p *StorageServer) format(path string) string {
	p.u.Path = path
	return p.u.String()
}
