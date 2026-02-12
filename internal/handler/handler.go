// Package handler обработчики запросов
package handler

import (
	"context"
	"net/url"
	"time"

	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/service"
)

// StorageServer сервер обработки запросов.
type StorageServer struct {
	model.Storage
	u     *url.URL
	audit model.Audit
}

// New создание сервера обработки запросов.
func New(address string) *StorageServer {
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	return &StorageServer{Storage: service.New(), u: u}
}

// SetLoader установка мета хранения данных.
func (p *StorageServer) SetLoader(loader model.Storage) {
	p.Storage = loader
}

// SetAudit установка места отправки аудита.
func (p *StorageServer) SetAudit(audit model.Audit) {
	p.audit = audit
}

func (p *StorageServer) format(path string) string {
	p.u.Path = path
	return p.u.String()
}

func (p *StorageServer) sendAudit(ctx context.Context, action, user, address string) {
	if p.audit == nil {
		return
	}

	p.audit.Send(ctx, &model.AuditData{
		Created: time.Now().Unix(),
		Action:  action,
		User:    user,
		Address: address,
	})
}
